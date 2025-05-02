package vm

import (
	"fmt"
	"github.com/singlaanish56/Compiler-in-go/code"
	"github.com/singlaanish56/Compiler-in-go/object"
	"github.com/singlaanish56/Compiler-in-go/compiler"
)


const StackSize = 2048
const GlobalSize= 65536
var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct{
	constants []object.Object
	instructions code.Instructions

	stack []object.Object
	stackPointer int

	globalStore []object.Object
}

func New(bytecode *compiler.Bytecode) *VM{
	return &VM{
		instructions: bytecode.Instructions,
		constants: bytecode.Constants,
		stack: make([]object.Object, StackSize),
		stackPointer: 0,
		globalStore: make([]object.Object, GlobalSize),
	}
}

func NewWithGlobalStore(bytecode *compiler.Bytecode, store []object.Object) *VM{
	vm := New(bytecode)
	vm.globalStore= store
	return vm
}

func (vm *VM) StackTop() object.Object{
	if vm.stackPointer == 0{
		return nil
	}

	return vm.stack[vm.stackPointer - 1]
}

func(vm * VM) LastPoppedStackElement() object.Object{
	return vm.stack[vm.stackPointer]
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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil{
				 return err
			}
		case code.OpTrue:
			if err:=vm.push(True);err != nil{
				return err
			}
		case code.OpFalse:
			if err:=vm.push(False);err != nil{
				return err
			}
		case code.OpGreaterThan, code.OpLessThan, code.OpEqual, code.OpNotEqual:
			if err := vm.executeComparison(op);err != nil{
				return err
			}
		case code.OpBang:
			if err := vm.executeBangOperation(); err != nil{
				return err
			}
		case code.OpMinus:
			if err := vm.executeMinusOperation(); err != nil{
				return err
			}
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[i+1:]))
			i+=2

			condition := vm.pop()
			if !isTruthy(condition){
				i=pos-1
			}
		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[i+1:]))
			i=pos-1
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[i+1:])
			i+=2
			vm.globalStore[globalIndex]=vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[i+1:])
			i+=2
			err := vm.push(vm.globalStore[globalIndex])
			if err != nil{
				return err
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil{
				return err
			}
		case code.OpArray:
			numelements := int(code.ReadUint16(vm.instructions[i+1:]))
			i+=2

			array := vm.buildArray(vm.stackPointer-numelements, vm.stackPointer)
			vm.stackPointer -= numelements
			err := vm.push(array)
			if err != nil{
				return err
			}
		case code.OpPop:
			vm.pop()
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

func (vm *VM) executeBinaryOperation(op code.Opcode) error{
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	
	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ{
		return vm.executeIntegerBinaryOperation(op, left, right)
	}else if leftType==object.STRING_OBJ && rightType == object.STRING_OBJ{
		return vm.executeStringBinaryOperation(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (vm *VM) executeIntegerBinaryOperation(op code.Opcode, left, right object.Object) error{
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	var result int64
	switch op{
	case code.OpAdd:
		result = leftVal + rightVal
	case code.OpSub:
		result = leftVal - rightVal
	case code.OpMul:
		result = leftVal * rightVal
	case code.OpDiv:
		result = leftVal / rightVal

	default:
		return fmt.Errorf("unsupported binary operation %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) executeStringBinaryOperation(operation code.Opcode, left,right object.Object) error{
	if operation != code.OpAdd{
		return fmt.Errorf("unkown string operation, %d", operation)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return vm.push(&object.String{leftValue+rightValue})
}

func (vm *VM) executeComparison(op code.Opcode) error{
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ || right.Type() == object.INTEGER_OBJ{
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op{
	case code.OpEqual:
		return vm.push(toBooleanObject(right==left))
	case code.OpNotEqual:
		return vm.push(toBooleanObject(right != left))
	default:
		return fmt.Errorf("unsupported comparison operation %d", op)
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error{
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op{
	case code.OpEqual:
		return vm.push(toBooleanObject(leftVal == rightVal))
	case code.OpNotEqual:
		return vm.push(toBooleanObject(leftVal != rightVal))
	case code.OpGreaterThan:
		return vm.push(toBooleanObject(leftVal > rightVal))
	case code.OpLessThan:
		return vm.push(toBooleanObject(leftVal < rightVal))
	default:
		return fmt.Errorf("unsupported comparison operation %d", op)
	}
}

func (vm *VM) executeBangOperation() error{
	right := vm.pop()

	switch right{
		case True:
			return vm.push(False)
		case False:
			return vm.push(True)
		case Null:
			return vm.push(True)
		default:
			return fmt.Errorf("unsupported bang operation %s", right.Type())	
	}
}

func (vm *VM) executeMinusOperation() error{
	right := vm.pop()

	if right.Type() != object.INTEGER_OBJ{
		return fmt.Errorf("unsupported type for minus operation %s", right.Type())
	}

	rightVal := right.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -rightVal})
}

func (vm *VM) buildArray(startIndex, endIndex int) *object.Array{
	elements := make([]object.Object, endIndex-startIndex)
	for i:=startIndex;i<endIndex;i++{
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}
func toBooleanObject(val bool) *object.Boolean{
	if val{
		 return True
	}

	return False
}

func isTruthy(obj object.Object) bool{
	switch obj := obj.(type){

		case *object.Boolean:
			return obj.Value
		case *object.Null:
			return false
		default:
			return true
	}
}
