package code

import(
	"testing"
)
func TestMakeBytecode(t *testing.T){

	tests := []struct{
		op Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
	}

	for _, tt := range tests{
		inst := Make(tt.op, tt.operands...)

		if len(inst) != len(tt.expected){
			t.Errorf("instruction has the wrong length expected=%d, got=%d", len(tt.expected), len(inst))
			return 
		}

		for i, b := range tt.expected{
			if inst[i] != b{
				t.Errorf("wroing bytes op %d , expected=%d, get=%d",i,b, inst[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T){
	instructions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpConstant 1
	0003 OpConstant 2
	0006 OpConstant 65535
	`

	concatted := Instructions{}

	for _, ins := range instructions{
		concatted = append(concatted, ins...)
	}

	if concatted.String() != expected{
		t.Errorf("wrong instruction string expected=%q, got=%q", expected, concatted.String())

	}
}

func TestReadOperands(t *testing.T){
	tests := []struct{
		operation Opcode
		operands []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests{
		instructions := Make(tt.operation, tt.operands...)


		definition , err := Lookup(byte(tt.operation))
		if err != nil{
			t.Errorf("lookup failed %s", err)
		}

		operandsRead, bytesRead := ReadOperands(definition, instructions[1:])
		if bytesRead != tt.bytesRead{
			t.Fatalf("wrong bytes read, expected=%d, got=%d", tt.bytesRead, bytesRead)
		}

		for i, want := range tt.operands{
			if operandsRead[i] != want{
				t.Errorf("wrong operand at %d, expected=%d, got=%d", i, want, operandsRead[i])
			}
		}
	}
}