package lexer

import (
	"testing"

	"github.com/mannyjimen/Monkey-Compiler/token"
)

func TestNextToken(t *testing.T) {
	input := `let hello = 20;
	let goodbye = 10;
	
	let add = fn(x, y) {
		x + y;
	};
	
	let result = add(hello, goodbye);

	!-/*7<>
	return if else true false
	5 != 10;
	6 == 6;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "hello"},
		{token.ASSIGN, "="},
		{token.INT, "20"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "goodbye"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "hello"},
		{token.COMMA, ","},
		{token.IDENT, "goodbye"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "7"},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.RETURN, "return"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.INT, "5"},
		{token.NOT_EQ, "!="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "6"},
		{token.EQ, "=="},
		{token.INT, "6"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
