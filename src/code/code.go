package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte
func (i Instructions) String() string{ return ""}


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
}

const (
	OpConstant Opcode = iota
)