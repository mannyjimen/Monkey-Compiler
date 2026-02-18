package ast

import (
	"fmt"
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

func TestStringForStringLit(t *testing.T) {
	astStringLit := StringLiteral{Token: token.Token{Type: token.STRING, Literal: "hello"}, Value: "hello"}

	expected := fmt.Sprintf("%chello%c", 34, 34)
	if astStringLit.String() != expected {
		t.Errorf("StringLiteral .String() incorrect, got %s, expected %s", astStringLit.String(), expected)
	}
}
