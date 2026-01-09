package parser

import (
	"fmt"
	"testing"

	"github.com/mannyjimen/Monkey-Compiler/ast"
	"github.com/mannyjimen/Monkey-Compiler/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 7;
	let barfoo = 50550;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	} else if len(program.Statements) != 3 {
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
		t.Errorf("s.TokenLiteral not 'let', got %q", s.TokenLiteral())
		return false
	}

	//first time seeing type assertion EVER
	//ok returns whether or not s is an *ast.LetStatement
	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s is not a *ast.LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not '%s', got %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not %s, got %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return foo(33);`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	} else if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, contains %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement, got %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}

}

func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

/*
TEMPORARY TEST
want to test
parsing input, now we have two statements in program.Statements

want to see if Let and Return statements' respective Value fields
are getting parsed to nil
*/

func TestNilValues(t *testing.T) {
	input := `
	let a = 5;
	return b;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	ls, ok := program.Statements[0].(*ast.LetStatement)

	if !ok {
		t.Errorf("1st stmt not an *ast.LetStatement, got %T type", program.Statements[0])
	}

	if ls.Value != nil {
		t.Errorf("LetStatement value not being set to nil (temp test)")
	}

	rs, ok := program.Statements[1].(*ast.ReturnStatement)

	if !ok {
		t.Errorf("2nd stmt not an *ast.ReturnStatement, got %T type", program.Statements[1])
	}

	if rs.ReturnValue != nil {
		t.Errorf("ReturnStatement's ReturnValue is not being set to nil (temp test)")
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Incorrect amount of statements in program, expected 1, got %d", len(program.Statements))
	}

	//check that the statement was parsed as *ast.ExpressionStatement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	//check that expr stmt's expression is an identifier
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expression is not an *ast.Identifier, got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value is incorrect, expected 'foobar', got %q", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral is incorrect, expected 'foobar', got %q", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	//first, check if program has 1 statement
	if len(program.Statements) != 1 {
		t.Fatalf("Incorrect amount of statements in program, expected 1, got %d", len(program.Statements))
	}

	//check if statement is expr statement
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement, got %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression literal is not an *ast.IntegerLiteral, got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("integer literal's value is not 5, got %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("integerLiteral.TokenLiteral is incorrect, expected '5', got %s", literal.TokenLiteral())
	}

}

func TestParsingPrefixOperators(t *testing.T) {
	prefixTests := []struct {
		testCase     string
		operator     string
		integerValue int64
	}{
		{`!5;`, `!`, 5},
		{`-4;`, `-`, 4},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.testCase)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		//checking length of program
		if len(program.Statements) != 1 {
			t.Fatalf("Incorrect number of statements parsed, expected 1, got %d", len(program.Statements))
		}

		// fmt.Println(program.Statements[0])

		//checking that it was parsed as an ExpressionStatement
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an expression statement, got %T", program.Statements[0])
		}

		//checking that is was parsed as PrefixExpression
		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not a ast.PrefixExpression, got %T", stmt.Expression)
		}

		if expr.Operator != tt.operator {
			t.Fatalf("expr operator is not %s, got %s", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.integerValue) {
			return
		}

	}
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, value int64) bool {
	integerLitExpr, ok := expr.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("Expression is not an *ast.IntegerLiteral, got %T", expr)
		return false
	}

	if integerLitExpr.Value != value {
		t.Errorf("Integer value is incorrect, expected %q, got %q", value, integerLitExpr)
		return false
	}

	if integerLitExpr.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("*ast.IntegerLiteral.TokenLiteral() method is incorrect, expected %d, got %s",
			value, integerLitExpr.TokenLiteral())
		return false
	}

	return true
}
