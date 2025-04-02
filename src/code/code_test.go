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