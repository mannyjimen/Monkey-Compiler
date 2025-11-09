package parser

import (
	"testing"

	"github.com/mannyjimen/Monkey-Compiler/ast"
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

	for i, tt := range tests {
		stmt := program.Statements[i]

		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let', got=%q", s.TokenLiteral())
		return false
	}

	//first time seeing type assertion EVER
	//ok returns whether or not s is an *ast.LetStatement
	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s is not a *ast.LetStatement, got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not '%s', got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not %s, got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}
