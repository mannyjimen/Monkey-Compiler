package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"

	//operators
	ASSIGN = ""
	PLUS   = "+"

	//delimeters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
