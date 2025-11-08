package parser

import (
	"testing"

	"github.com/mannyjimen/Monkey-Compiler/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 7;
	let barfoo = 44;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, contains %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"barfoo"},
	}

	//making compiler happy :) for now
	_ = tests

	// for i, tt := range tests {
	// 	statement := program.Statements[i]

	// }
}
