package parser

import (
	"github.com/mannyjimen/Monkey-Compiler/ast"
	"github.com/mannyjimen/Monkey-Compiler/lexer"
	"github.com/mannyjimen/Monkey-Compiler/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token
}

// concise same names is fine for idiomatic go ({l: l})
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	//setting up parser with curr and peek tokens
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	return nil
}
