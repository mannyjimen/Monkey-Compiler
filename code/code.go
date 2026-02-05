package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//defining bytecode format

type Instructions []byte

func (instr Instructions) String() string {
	var out bytes.Buffer

	for i := 0; i < len(instr); {
		opByte := instr[i]
		def, err := Lookup(opByte)
		if err != nil {
			return fmt.Sprintf("ERROR: %s\n", err)
		}

		operands, bytesRead := ReadOperands(def, instr[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, fmtInstructions(def, operands))

		i += 1 + bytesRead
	}
	return out.String()
}

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpBang
	OpMinus
	OpJump
	OpJumpNotTruthy
)

// constructing bytecode instruction, op code and operand constant pool addresses
// returns slice of bytes (instruction turned to bytecode)
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		// case 0:
		// 	return instruction
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}

		offset += width
	}

	return instruction
}

type Definition struct {
	Name          string
	OperandWidths []int
}

// opcode : {string representation, byte widths of operands expected after opcode}
var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpJump:          {"OpJump", []int{2}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// disassembles bytecode instruction operands
func ReadOperands(def *Definition, instr Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(instr[offset:]))
		}
		offset += width
	}

	return operands, offset
}

func ReadUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func fmtInstructions(def *Definition, operands []int) string {
	if len(operands) != len(def.OperandWidths) {
		return fmt.Sprintf("ERROR: operand len %d does not match defined len %d",
			len(operands), len(def.OperandWidths))
	}

	switch len(operands) {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: operand len of %d not handled for %s", len(operands), def.Name)
}
