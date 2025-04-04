package code

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

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