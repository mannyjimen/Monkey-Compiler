package ast

import (
	"testing"

	"github.com/mannyjimen/Monkey-Compiler/token"
)

// "let myVar = anotherVar"
func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		//Errorf because this error doesn't make subsequent tests impossible (philosophical)
		t.Errorf("program.String() incorrect, got=%q", program.String())
	}
}
