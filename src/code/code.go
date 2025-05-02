package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte
func (ins Instructions) String() string{
	var out bytes.Buffer

	i := 0
	for i< len(ins){
		def, err := Lookup(ins[i])
		if err != nil{
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands , bytesRead := ReadOperands(def, ins[i+1:])
		
		fmt.Fprintf(&out , "%04d %s\n\t", i, ins.instructionToFmt(def, operands))

		i += 1+ bytesRead
	}

	return out.String()
}

func (ins Instructions) instructionToFmt(def *Definition, operands []int) string{
	operandsExpectedLen := len(def.OperandWidths)
	if len(operands) != operandsExpectedLen{
		return fmt.Sprintf("ERROR: wrong number of operands. expected=%d, got=%d", operandsExpectedLen, len(operands))
	}

	switch operandsExpectedLen{
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

type Opcode byte



type Definition struct{
	Name string
	OperandWidths []int
}

func Make(op Opcode, operands...int) []byte{
	def, ok := definitions[op]
	if !ok{
		return []byte{}
	}

	instructionLen:=1
	for _, w := range def.OperandWidths{
		instructionLen+=w
	}

	instruction := make([]byte, instructionLen)
	instruction[0]=byte(op)

	offset :=1
	for i,o := range operands{
		width := def.OperandWidths[i]
		switch width{
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset +=width

	}

	return instruction
}

func ReadOperands(definition *Definition, ins Instructions) ([]int ,int){
	operands := make([]int, len(definition.OperandWidths))
	offset:=0
	for i, width := range definition.OperandWidths{
		switch width{
		case 2 :
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16{
	return binary.BigEndian.Uint16(ins)
}

func Lookup(op byte) (*Definition, error){
	def , ok := definitions[Opcode(op)]
	if !ok{
		return nil, fmt.Errorf("opcode is undefined %d", op)
	}

	return def, nil
}

var definitions = map[Opcode] *Definition{
	OpConstant : {"OpConstant", []int{2}},
	OpAdd : {"OpAdd", []int{}},
	OpPop : {"OpPop", []int{}},
	OpSub : {"OpSub", []int{}},
	OpMul : {"OpMul", []int{}},
	OpDiv : {"OpDiv", []int{}},
	OpTrue : {"OpTrue", []int{}},
	OpFalse : {"OpFalse", []int{}},
	OpEqual : {"OpEqual", []int{}},
	OpNotEqual : {"OpNotEqual", []int{}},
	OpGreaterThan : {"OpGreaterThan", []int{}},
	OpLessThan : {"OpLessThan", []int{}},
	OpMinus : {"OpMinus", []int{}},
	OpBang : {"OpBang", []int{}},	
	OpJumpNotTruthy : {"OpJumpNotTruthy", []int{2}},
	OpJump : {"OpJump", []int{2}},
	OpNull:{"OpNull", []int{}},
	OpSetGlobal: {"OpSetGlobal", []int{2}},
	OpGetGlobal: {"OpGetGlobal", []int{2}},
	OpArray:{"OpArray",[]int{2}},
}

const (
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpLessThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
	OpNull
	OpSetGlobal
	OpGetGlobal
	OpArray
)
