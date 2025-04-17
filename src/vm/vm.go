package vm

import (
	"fmt"
	"github.com/singlaanish56/Compiler-in-go/code"
	"github.com/singlaanish56/Compiler-in-go/object"
	"github.com/singlaanish56/Compiler-in-go/compiler"
)


const StackSize = 2048

type VM struct{
	constants []object.Object
	instructions code.Instructions

	stack []object.Object
	stackPointer int
}

func New(bytecode *compiler.Bytecode) *VM{
	return &VM{
		instructions: bytecode.Instructions,
		constants: bytecode.Constants,
		stack: make([]object.Object, StackSize),
		stackPointer: 0,
	}
}

func (vm *VM) StackTop() object.Object{
	if vm.stackPointer == 0{
		return nil
	}

	return vm.stack[vm.stackPointer - 1]
}

func (vm *VM) Run() error{
	for i:=0; i< len(vm.instructions); i++{
		op := code.Opcode(vm.instructions[i])
		switch op{
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[i+1:])
			i+=2
			err := vm.push(vm.constants[constIndex])
			if err != nil{
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()

			leftVal  := left.(*object.Integer).Value
			rightVal := right.(*object.Integer).Value
			result := leftVal + rightVal
			vm.push(&object.Integer{Value: result})
		}
	}

	return nil
}

func (vm *VM) push(o object.Object) error{
	if vm.stackPointer >= StackSize{
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.stackPointer] = o
	vm.stackPointer++
	return nil
}

func (vm *VM) pop() object.Object{
	o := vm.stack[vm.stackPointer - 1]
	vm.stackPointer--
	return o
}