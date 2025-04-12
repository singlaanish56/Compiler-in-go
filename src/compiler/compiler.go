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
	return nil
}

func (c *Compiler) Bytecode() *Bytecode{
	return &Bytecode{
		Instructions: c.instructions,
		Constants: c.constants,
	}
}
