package vm

import (
	"fmt"
	"github.com/singlaanish56/Compiler-in-go/code"
	"github.com/singlaanish56/Compiler-in-go/object"
	"github.com/singlaanish56/Compiler-in-go/compiler"
)


const StackSize = 2048
const GlobalSize= 65536
const MaxFrames = 1024
var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type Frame struct{
	fn *object.CompiledFunction
	ip int
}

func NewFrame(fn * object.CompiledFunction) *Frame{
	return &Frame{fn, -1}
}

func (f *Frame) Instructions() code.Instructions{
	return f.fn.Instructions
}

type VM struct{
	constants []object.Object
	
	frames []*Frame
	framesIndex int

	stack []object.Object
	stackPointer int

	globalStore []object.Object
}

func New(bytecode *compiler.Bytecode) *VM{
	mainFn := &object.CompiledFunction{bytecode.Instructions}
	mainFrame := NewFrame(mainFn)
	frames := make([]*Frame, MaxFrames)
	frames[0]=mainFrame

	return &VM{
		constants: bytecode.Constants,
		frames: frames,
		framesIndex: 1,
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

	var i int
	var ins code.Instructions
	var op code.Opcode


	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1{
		vm.currentFrame().ip++

		i= vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[i])

		switch op{
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[i+1:])
			vm.currentFrame().ip+=2
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
			pos := int(code.ReadUint16(ins[i+1:]))
			vm.currentFrame().ip+=2

			condition := vm.pop()
			if !isTruthy(condition){
				vm.currentFrame().ip=pos-1
			}
		case code.OpJump:
			pos := int(code.ReadUint16(ins[i+1:]))
			vm.currentFrame().ip=pos-1
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[i+1:])
			vm.currentFrame().ip+=2
			vm.globalStore[globalIndex]=vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[i+1:])
			vm.currentFrame().ip+=2
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
			numelements := int(code.ReadUint16(ins[i+1:]))
			vm.currentFrame().ip+=2

			array := vm.buildArray(vm.stackPointer-numelements, vm.stackPointer)
			vm.stackPointer -= numelements
			err := vm.push(array)
			if err != nil{
				return err
			}
		case code.OpHash:
			numOfElements := int(code.ReadUint16(ins[i+1:]))
			vm.currentFrame().ip+=2

			hash, err := vm.buildHash(vm.stackPointer-numOfElements, vm.stackPointer)
			if err !=nil{
				return err
			}

			vm.stackPointer-=numOfElements

			err = vm.push(hash)
			if err != nil{
				return err
			}
		case code.OpIndex:
			index:= vm.pop()
			objectToBeIndexed := vm.pop()

			err:= vm.executeIndexExpression(objectToBeIndexed, index)
			if err != nil{
				return err
			}
		case code.OpCall:
			fn, ok := vm.stack[vm.stackPointer-1].(*object.CompiledFunction)
			if !ok{
				return fmt.Errorf("object is not a function, got=%T", vm.stack[vm.stackPointer-1])
			}
			frame := NewFrame(fn)
			vm.pushFrame(frame)
		case code.OpReturnValue:
			returnValue := vm.pop()

			vm.popFrame()
			vm.pop()

			err := vm.push(returnValue)
			if err != nil{
				return err
			}
		case code.OpReturn:
			vm.popFrame()
			vm.pop()

			err := vm.push(Null)
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

func (vm *VM) currentFrame() *Frame{
	return vm.frames[vm.framesIndex - 1]
}

func (vm *VM) pushFrame(frame *Frame){
	vm.frames[vm.framesIndex]= frame
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame{
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
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

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error){
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i:=startIndex;i<endIndex;i+=2{
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}

		hashKey, ok := key.(object.Hashable)
		if !ok{
			return nil, fmt.Errorf("unhashable type %s", key.Type())
		}

		hashedPairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) executeIndexExpression(objectToBeIndexed, index object.Object) error{
	switch{
		case objectToBeIndexed.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
			return vm.executeArrayIndex(objectToBeIndexed, index)
		case objectToBeIndexed.Type() == object.HASHPAIR_OBJ:
			return vm.executeHashIndex(objectToBeIndexed, index)
		default:
			return fmt.Errorf("index operator not supported: %s %s", objectToBeIndexed.Type(), index.Type())
	}
}

func (vm *VM) executeArrayIndex(array, index object.Object) error{
	arrayObject, ok := array.(*object.Array)
	if !ok{
		return fmt.Errorf("object is not an array, got=%T", array)
	}

	indexObject, ok := index.(*object.Integer)
	if !ok{
		return fmt.Errorf("object is not an integer, got=%T", index)
	}

	i:= indexObject.Value
	maxIndex := int64(len(arrayObject.Elements))
	if i<0 || i>=maxIndex{
		return vm.push(Null)
	}

	return vm.push(arrayObject.Elements[i])
}

func (vm *VM) executeHashIndex(hash, index object.Object) error{
	hashObject, ok := hash.(*object.Hash)
	if !ok{
		return fmt.Errorf("object is not a hash, got=%T", hash)
	}

	hashKey, ok := index.(object.Hashable)
	if !ok{
		return fmt.Errorf("unhashable type %s", index.Type())
	}

	pair, ok := hashObject.Pairs[hashKey.HashKey()]
	if !ok{
		return vm.push(Null)
	}

	return vm.push(pair.Value)
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
