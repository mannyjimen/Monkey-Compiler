package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
	}

	for _, tt := range tests {
		//first time using variadic function
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Fatalf("instruction has incorrect length, expected %d, got %d",
				len(instruction), len(tt.expected))
		}

		for i, b := range instruction {
			if tt.expected[i] != instruction[i] {
				t.Errorf("instruction at pos %d is incorrect, expected %d, got %d",
					i, tt.expected[i], b)
			}
		}
	}
}
