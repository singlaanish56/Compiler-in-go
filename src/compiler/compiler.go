package compiler

import (
	"github.com/singlaanish56/Compiler-in-go/code"
	"github.com/singlaanish56/Compiler-in-go/object"
	"github.com/singlaanish56/Compiler-in-go/ast"
)

type Compiler struct{
	instructions code.Instructions
	constants []object.Object
}

type Bytecode struct{
	Instructions code.Instructions
	Constants []object.Object
}

func New() *Compiler{
	return &Compiler{
		instructions : code.Instructions{},
		constants: []object.Object{},
	}
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
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil{
			return err
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
	case *ast.IntegerLiteral:
		integerObject := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integerObject))

	}

	return nil
}
	
func (c *Compiler) Bytecode() *Bytecode{
	return &Bytecode{
		Instructions: c.instructions,
		Constants: c.constants,
	}
}

//returns the index of the last instructions that is  set in the list
func (c *Compiler) emit(operation code.Opcode, operands... int) int{
	ins := code.Make(operation, operands...)
	lastInstructionPos := c.addInstruction(ins)
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