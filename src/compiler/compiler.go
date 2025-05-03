package compiler

import (
	"fmt"
	"sort"

	"github.com/singlaanish56/Compiler-in-go/ast"
	"github.com/singlaanish56/Compiler-in-go/code"
	"github.com/singlaanish56/Compiler-in-go/object"
)

type Compiler struct{
	instructions code.Instructions
	constants []object.Object
	lastInstruction EmittedInstruction
	previousInstruction EmittedInstruction
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

func New() *Compiler{
	return &Compiler{
		instructions : code.Instructions{},
		constants: []object.Object{},
		lastInstruction: EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
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
		c.emit(code.OpSetGlobal, symbol.Position)
	case *ast.BlockStatement:
		for _,s := range node.Statements{
			err := c.Compile(s)
			if err != nil{
				return err
			}
		}
	case *ast.ArrayLiteral:
		for _, element := range node.Elements{
			err := c.Compile(element)
			if err != nil{
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil{
			return err
		}
		c.emit(code.OpPop)
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

		if c.lastInstructionIsPop(){
			c.removeLastPop()
		}

		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.instructions)
		c.changeOperand(jumpNotTruthyPosition, afterConsequencePos)

		if node.Alternative == nil{
			c.emit(code.OpNull)
		}else{

			err := c.Compile(node.Alternative)
			if err != nil{
				return err
			}

			if c.lastInstructionIsPop(){
				c.removeLastPop()
			}
		}

		alternativePos := len(c.instructions)
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
	case *ast.Variable:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok{
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		c.emit(code.OpGetGlobal, symbol.Position)
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
	}

	return nil
}
	
func (c *Compiler) Bytecode() *Bytecode{
	return &Bytecode{
		Instructions: c.instructions,
		Constants: c.constants,
	}
}

func (c *Compiler) emit(operation code.Opcode, operands... int) int{
	ins := code.Make(operation, operands...)
	lastInstructionPos := c.addInstruction(ins)

	c.setLastInstruction(operation, lastInstructionPos)

	return lastInstructionPos
}

func (c *Compiler) addInstruction(ins []byte) int{
	currentInstructionLen := len(c.instructions)
	c.instructions= append(c.instructions, ins...)
	return currentInstructionLen	
}

func (c *Compiler) addConstant(object object.Object) int{
	c.constants=append(c.constants, object)
	return len(c.constants) -1
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int){
	previous := c.lastInstruction
	last:= EmittedInstruction{Opcode: op, Position: pos}

	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) lastInstructionIsPop() bool{
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop(){
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) changeOperand(operationPosition, operand int){
	op := code.Opcode(c.instructions[operationPosition])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(operationPosition, newInstruction)
}

func (c *Compiler) replaceInstruction(operationPosition int, newInstruction []byte){
	for i:=0;i<len(newInstruction);i++{
		c.instructions[operationPosition+i] = newInstruction[i]
	}
}