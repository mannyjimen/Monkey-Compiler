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

func TestInstructionsString(t *testing.T) {
	tests := []struct {
		instructions []Instructions
		expected     string
	}{
		{
			instructions: []Instructions{
				Make(OpConstant, 1),
				Make(OpConstant, 2),
				Make(OpConstant, 65534),
			},
			expected: `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65534
`,
		},
	}

	for _, tt := range tests {
		flattened := Instructions{}
		for _, instr := range tt.instructions {
			flattened = append(flattened, instr...)
		}

		actual := flattened.String()
		if actual != tt.expected {
			t.Errorf("Instructions String() method is incorrect, expected %q, got %q",
				tt.expected, actual)
		}
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{{OpConstant, []int{65534}, 2}}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("could not find definition: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("wrong number of bytes read, expected %d, got %d",
				tt.bytesRead, n)
		}

		for i, o := range operandsRead {
			if o != tt.operands[i] {
				t.Errorf("wrong operand read, expected %d, got %d",
					tt.operands[i], o)
			}
		}
	}
}
