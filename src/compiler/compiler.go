package compiler

import (
	"fmt"
	"sort"

	"github.com/singlaanish56/Compiler-in-go/ast"
	"github.com/singlaanish56/Compiler-in-go/code"
	"github.com/singlaanish56/Compiler-in-go/object"
)

type Compiler struct{
	constants []object.Object
	compilerScopes []CompilationScope
	scopeIndex int
	symbolTable *SymbolTable
}

type Bytecode struct{
	Instructions code.Instructions
	Constants []object.Object
}

type EmittedInstruction struct{
	Opcode code.Opcode
	Position int
}

type CompilationScope struct{
	instructions code.Instructions
	lastInstruction EmittedInstruction
	previousInstruction EmittedInstruction
}

func New() *Compiler{
	mainScope := CompilationScope{
		instructions: code.Instructions{},
		lastInstruction: EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	return &Compiler{
		constants: []object.Object{},
		compilerScopes: []CompilationScope{mainScope},
		scopeIndex: 0,
		symbolTable: NewSymbolTable(),
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler{
	compiler := New()
	compiler.symbolTable=s;
	compiler.constants=constants
	return compiler
}

func (c *Compiler) Compile(node ast.ASTNode) error{
	switch node := node.(type){
	case *ast.AstRootNode:
		for _, statement := range node.Statements{
			err := c.Compile(statement)
			if err !=nil{
				return err
			}
		}
	case *ast.LetStatement:
		err:= c.Compile(node.Value)
		if err != nil{
			return err
		}
		symbol := c.symbolTable.Define(node.Variable.Value)
		if symbol.Scope == GlobalScope{
			c.emit(code.OpSetGlobal, symbol.Position)
		}else{
			c.emit(code.OpSetLocal, symbol.Position)
		}
	case *ast.BlockStatement:
		for _,s := range node.Statements{
			err := c.Compile(s)
			if err != nil{
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil{
			return err
		}
		c.emit(code.OpPop)
	case *ast.ReturnStatement:
		err := c.Compile(node.Value)
		if err != nil{
			return err
		}

		c.emit(code.OpReturnValue)
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil{
			return err
		}

		//dummy code to jump to
		jumpNotTruthyPosition := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil{
			return err
		}

		if c.lastInstructionIs(code.OpPop){
			c.removeLastPop()
		}

		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPosition, afterConsequencePos)

		if node.Alternative == nil{
			c.emit(code.OpNull)
		}else{

			err := c.Compile(node.Alternative)
			if err != nil{
				return err
			}

			if c.lastInstructionIs(code.OpPop){
				c.removeLastPop()
			}
		}

		alternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, alternativePos)
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err!=nil{
			return err
		}

		switch node.Operator{
			case "!":
				c.emit(code.OpBang)
			case "-":
				c.emit(code.OpMinus)
			default:
				fmt.Errorf("unknown prefix operator %s", node.Operator)
		}
	case *ast.InfixExpression:
		err := c.Compile(node.Left)
		if err != nil{
			return err
		}

		err  =c.Compile(node.Right)
		if err != nil{
			return err
		}

		switch node.Operator{
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "<":
			c.emit(code.OpLessThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err !=nil{
			return err
		}

		err = c.Compile(node.Index)
		if err != nil{
			return err
		}

		c.emit(code.OpIndex)
	case *ast.FunctionExpression:
		c.enterScope()

		for _, param := range node.Parameters{
			c.symbolTable.Define(param.Value)
		}
		
		err := c.Compile(node.Body)
		if err != nil{
			return err
		}

		if c.lastInstructionIs(code.OpPop){
			c.replaceLastPopWithReturn()
		}

		if !c.lastInstructionIs(code.OpReturnValue){
			c.emit(code.OpReturn)
		}
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()
		compiledFn := &object.CompiledFunction{
			Instructions: instructions,
			NumberOfLocals: numLocals,
		}
		c.emit(code.OpConstant, c.addConstant(compiledFn))
	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil{
			return err
		}

		for _, arg := range node.Arguments{
			err := c.Compile(arg)
			if err != nil{
				return err
			}
		}

		c.emit(code.OpCall, len(node.Arguments))
	case *ast.Variable:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok{
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		if symbol.Scope == GlobalScope{
			c.emit(code.OpGetGlobal, symbol.Position)
		}else{
			c.emit(code.OpGetLocal, symbol.Position)
		}
	case *ast.IntegerLiteral:
		integerObject := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integerObject))
	case *ast.BooleanLiteral:
		if node.Value{
			c.emit(code.OpTrue)
		}else{
			c.emit(code.OpFalse)
		}
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))
	case *ast.HashLiteral:

		keys := []ast.Expression{}
		for k := range node.Pairs{
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i,j int) bool{
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys{
			err := c.Compile(k)
			if err != nil{
				return err
			}

			err = c.Compile(node.Pairs[k])
			if err !=nil{
				return err
			}
		}

		c.emit(code.OpHash, len(node.Pairs)*2)
	case *ast.ArrayLiteral:
		for _, element := range node.Elements{
			err := c.Compile(element)
			if err != nil{
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	}

	return nil
}
	
func (c *Compiler) Bytecode() *Bytecode{
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants: c.constants,
	}
}

func (c *Compiler) emit(operation code.Opcode, operands... int) int{
	ins := code.Make(operation, operands...)
	lastInstructionPos := c.addInstruction(ins)

	c.setLastInstruction(operation, lastInstructionPos)

	return lastInstructionPos
}


func (c *Compiler) addConstant(object object.Object) int{
	c.constants=append(c.constants, object)
	return len(c.constants) -1
}

func (c *Compiler) currentInstructions() code.Instructions{
	return c.compilerScopes[c.scopeIndex].instructions
}

func (c *Compiler) addInstruction(ins []byte) int{
	posNewInstruction := len(c.currentInstructions())
	updateInstruction := append(c.currentInstructions(), ins...)

	c.compilerScopes[c.scopeIndex].instructions = updateInstruction

	return posNewInstruction
}	

func (c *Compiler) setLastInstruction(op code.Opcode, pos int){
	previous := c.compilerScopes[c.scopeIndex].lastInstruction	
	last:= EmittedInstruction{Opcode: op, Position: pos}

	c.compilerScopes[c.scopeIndex].previousInstruction = previous
	c.compilerScopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool{
	if len(c.currentInstructions()) ==0{
		return false
	}

	return c.compilerScopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop(){
	last := c.compilerScopes[c.scopeIndex].lastInstruction
	prev := c.compilerScopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]
	
	c.compilerScopes[c.scopeIndex].instructions = new
	c.compilerScopes[c.scopeIndex].lastInstruction = prev
}

func (c *Compiler) replaceLastPopWithReturn(){
	lastPos := c.compilerScopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.compilerScopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) changeOperand(operationPosition, operand int){
	op := code.Opcode(c.currentInstructions()[operationPosition])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(operationPosition, newInstruction)
}

func (c *Compiler) replaceInstruction(operationPosition int, newInstruction []byte){
	curr := c.currentInstructions()
	for i:=0;i<len(newInstruction);i++{
		curr[operationPosition+i] = newInstruction[i]
	}
}

func (c *Compiler) enterScope(){
	newScope := CompilationScope{
		instructions: code.Instructions{},
		lastInstruction: EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.compilerScopes = append(c.compilerScopes, newScope)
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
	c.scopeIndex++
}

func(c *Compiler) leaveScope() code.Instructions{
	curr := c.currentInstructions()

	c.compilerScopes = c.compilerScopes[:len(c.compilerScopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer
	return curr
}