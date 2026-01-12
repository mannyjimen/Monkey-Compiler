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
	t.Helper()
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
	return foo;`

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
	t.Helper()
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

// func TestNilValues(t *testing.T) {
// 	input := `
// 	let a = 5;
// 	return b;`

// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.ParseProgram()
// 	checkParseErrors(t, p)

// 	ls, ok := program.Statements[0].(*ast.LetStatement)

// 	if !ok {
// 		t.Errorf("1st stmt not an *ast.LetStatement, got %T type", program.Statements[0])
// 	}

// 	if ls.Value != nil {
// 		t.Errorf("LetStatement value not being set to nil (temp test)")
// 	}

// 	rs, ok := program.Statements[1].(*ast.ReturnStatement)

// 	if !ok {
// 		t.Errorf("2nd stmt not an *ast.ReturnStatement, got %T type", program.Statements[1])
// 	}

// 	if rs.ReturnValue != nil {
// 		t.Errorf("ReturnStatement's ReturnValue is not being set to nil (temp test)")
// 	}
// }

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParseErrors(t, p)

	checkOneStatementInProgram(t, program)
	stmt := checkAndGetExpressionStatement(t, program.Statements[0])

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

	checkOneStatementInProgram(t, program)
	stmt := checkAndGetExpressionStatement(t, program.Statements[0])

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression literal is not an *ast.IntegerLiteral, got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("integer literal's value is not 5, got %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("integerLiteral.TokenLiteral is incorrect, expected '5', got %q", literal.TokenLiteral())
	}

}

func TestBooleanExpression(t *testing.T) {
	input := `true`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	checkOneStatementInProgram(t, program)
	stmt := checkAndGetExpressionStatement(t, program.Statements[0])
	expr := stmt.Expression

	literal, ok := expr.(*ast.Boolean)
	if !ok {
		t.Fatalf("Expression literal is not an *ast.Boolean, got %T", expr)
	}

	if literal.Value != true {
		t.Errorf("Boolean literal value is not incorrect, expected true, got %t", literal.Value)
	}

	if literal.TokenLiteral() != "true" {
		t.Errorf("*ast.Boolean.TokenLiteral is incorrect, expected 'true', got %q", literal.TokenLiteral())
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	checkOneStatementInProgram(t, program)
	stmt := checkAndGetExpressionStatement(t, program.Statements[0])

	expr := stmt.Expression

	ifExpr, ok := expr.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expr is not an *ast.IfExpression, got %T", expr)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatalf("Consequence has incorrect amount of statements, expected 1, got %d", len(ifExpr.Consequence.Statements))
	}

	exprStmt, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence's only statement is not an *ast.ExpressionStatement, got %T",
			ifExpr.Consequence.Statements[0])
	}

	if !testIdentifier(t, exprStmt.Expression, "x") {
		return
	}

	if ifExpr.Alternative != nil {
		t.Fatalf("Alternative is not nil, got %+v", ifExpr.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	checkOneStatementInProgram(t, program)
	stmt := checkAndGetExpressionStatement(t, program.Statements[0])

	expr := stmt.Expression

	ifExpr, ok := expr.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expr is not an *ast.IfExpression, got %T", expr)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatalf("Consequence has incorrect amount of statements, expected 1, got %d", len(ifExpr.Consequence.Statements))
	}

	consequence, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence's only statement is not an *ast.ExpressionStatement, got %T",
			ifExpr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	//alternative testing

	if ifExpr.Alternative == nil {
		t.Fatalf("Alternative was incorrectly parsed as nil")
	}

	if len(ifExpr.Alternative.Statements) != 1 {
		t.Fatalf("Alternative has incorrect amount of statements, expected 1, got %d", len(ifExpr.Alternative.Statements))
	}

	alternative, ok := ifExpr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative's only statement is not an *ast.ExpressionStatement, got %T",
			ifExpr.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralExpression(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)

	checkOneStatementInProgram(t, program)
	stmt := checkAndGetExpressionStatement(t, program.Statements[0])
	expr := stmt.Expression

	funcLitExpr, ok := expr.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr is not an *ast.FunctionLiteral, got %T", expr)
	}

	//testing parsed parameters

	if len(funcLitExpr.Parameters) != 2 {
		t.Fatalf("FunctionLiteral parameter count is incorrect, expected 2, got %d", len(funcLitExpr.Parameters))
	}

	testLiteralExpression(t, funcLitExpr.Parameters[0], "x")
	testLiteralExpression(t, funcLitExpr.Parameters[1], "y")

	//testing parsed block

	if len(funcLitExpr.Body.Statements) != 1 {
		t.Fatalf("FunctionLiteral block statement count is incorrect, expected 1, got %d",
			len(funcLitExpr.Body.Statements))
	}

	bodyStmt := checkAndGetExpressionStatement(t, funcLitExpr.Body.Statements[0])

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")

}

func TestParsingFunctionParameters(t *testing.T) {
	tests := []struct {
		input      string
		parameters []string
	}{{"fn() { x };", []string{}},
		{"fn(x) { x };", []string{"x"}},
		{"fn(x, y, z) { x };", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		checkOneStatementInProgram(t, program)
		stmt := checkAndGetExpressionStatement(t, program.Statements[0])
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.parameters) {
			t.Fatalf("Incorrect parameter count, expected %d, got %d", len(tt.parameters), len(function.Parameters))
		}

		for i, param := range function.Parameters {
			testLiteralExpression(t, param, tt.parameters[i])
		}
	}
}

func TestParsingPrefixOperators(t *testing.T) {
	prefixTests := []struct {
		testCase     string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-4;", "-", 4},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.testCase)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		//checking length of program
		checkOneStatementInProgram(t, program)
		stmt := checkAndGetExpressionStatement(t, program.Statements[0])

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

func TestParsingInfixOperators(t *testing.T) {
	infixTests := []struct {
		testCase   string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5+5;", 5, "+", 5},
		{"5-5;", 5, "-", 5},
		{"5*5;", 5, "*", 5},
		{"5/5;", 5, "/", 5},
		{"5<5;", 5, "<", 5},
		{"5>5;", 5, ">", 5},
		{"5==5;", 5, "==", 5},
		{"5!=5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.testCase)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		checkOneStatementInProgram(t, program)
		stmt := checkAndGetExpressionStatement(t, program.Statements[0])

		expr, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not an *ast.InfixExpression, got %T", stmt.Expression)
		}

		if !testIntegerLiteral(t, expr.Left, tt.leftValue) {
			return
		}

		if expr.Operator != tt.operator {
			t.Errorf("expr operator is incorrect,  expected %s, got %s", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		testCase string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"5+4*6;",
			"(5 + (4 * 6))",
		},
		{
			"a + b - c / 4 + d * 5",
			"(((a + b) - (c / 4)) + (d * 5))",
		},
		{
			"5 < 4 != 4 > 6",
			"((5 < 4) != (4 > 6))",
		},
		{
			"5 < 4 * 7 == 4 - 3 > 6",
			"((5 < (4 * 7)) == ((4 - 3) > 6))",
		},
		{
			"5 - 4 * 7 / 4 + 6 - c * d + z",
			"((((5 - ((4 * 7) / 4)) + 6) - (c * d)) + z)",
		},
		{
			"5 - 4 == true;",
			"((5 - 4) == true)",
		},
		{
			"true",
			"true",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"-(2 + 2)",
			"(-(2 + 2))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.testCase)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		output := program.String()

		if output != tt.expected {
			t.Errorf("Parsed statement incorrectly, expected %q, got %q",
				tt.expected,
				output)
		}
	}
}

func TestParsingPrefixWithLiterals(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5", "!", 5},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		exprStmt := program.Statements[0].(*ast.ExpressionStatement)
		expr := exprStmt.Expression

		if !testPrefixExpression(t, expr, tt.operator, tt.value) {
			continue
		}
	}
}
func TestParsingInfixWithLiterals(t *testing.T) {
	tests := []struct {
		input    string
		left     any
		operator string
		right    any
	}{
		{"5+4;", 5, "+", 4},
		{"hello * goodbye;", "hello", "*", "goodbye"},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseErrors(t, p)

		exprStmt := program.Statements[0].(*ast.ExpressionStatement)
		expr := exprStmt.Expression

		if !testInfixExpression(t, expr, tt.left, tt.operator, tt.right) {
			continue
		}
	}
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, value int64) bool {
	t.Helper()
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

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	t.Helper()
	ident, ok := expr.(*ast.Identifier)

	if !ok {
		t.Errorf("Expression is not an *ast.Identifier, got %T", expr)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value is incorrect, expected %q, got %q", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral is incorrect, expected %s, got %s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	t.Helper()
	boolExpr, ok := expr.(*ast.Boolean)

	if !ok {
		t.Errorf("expr is not an *ast.Boolean, got %T", expr)
		return false
	}

	if boolExpr.Value != value {
		t.Errorf("bool expression value is incorrect, expected %t, got %t", value, boolExpr.Value)
		return false
	}

	if boolExpr.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bool expression value is incorrect, expected %t, got %s",
			value,
			boolExpr.TokenLiteral())
		return false
	}

	return true
}

// tests whether expression `expr` matches literal `expected`
func testLiteralExpression(t *testing.T, expr ast.Expression, expected any) bool {
	t.Helper()
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}

	t.Errorf("Type of expr is not handled, got %T", expr)
	return false
}

func testPrefixExpression(t *testing.T, expr ast.Expression,
	operator string, right any) bool {

	t.Helper()
	prefix, ok := expr.(*ast.PrefixExpression)
	if !ok {
		t.Errorf("expr is not an *ast.PrefixExpression, got %T", expr)
		return false
	}

	if prefix.Operator != operator {
		t.Errorf("expr operator is incorrect, expected %q, got %q", operator, prefix.Operator)
		return false
	}

	if !testLiteralExpression(t, prefix.Right, right) {
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, expr ast.Expression,
	left any, operator string, right any) bool {

	t.Helper()
	infix, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expr is not an *ast.InfixExpression, got %T", expr)
		return false
	}

	if !testLiteralExpression(t, infix.Left, left) {
		return false
	}

	if operator != infix.Operator {
		t.Errorf("expr operator is incorrect, expected %q, got %q", operator, infix.Operator)
		return false
	}

	if !testLiteralExpression(t, infix.Right, right) {
		return false
	}

	return true
}

func checkOneStatementInProgram(t *testing.T, program *ast.Program) {
	t.Helper()
	if len(program.Statements) != 1 {
		t.Fatalf("Incorrect amount of statements in program, expected 1, got %d", len(program.Statements))
	}
}

func checkAndGetExpressionStatement(t *testing.T, stmt ast.Statement) *ast.ExpressionStatement {
	t.Helper()
	exprStmt, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an *ast.ExpressionStatement, got %T", stmt)
	}

	return exprStmt
}
