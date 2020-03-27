package parser

// Package parser implements the parser that takes as input tokens from the
// lexer and produces as output an AST (Abstract Syntax Tree).

import (
	"fmt"

	"github.com/cedrickchee/hou/ast"
	"github.com/cedrickchee/hou/lexer"
	"github.com/cedrickchee/hou/token"
)

// Parser implements the parser.
type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token
}

// New constructs a new Parser with a Lexer as input.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

// Errors check if the parser encountered any errors.
func (p *Parser) Errors() []string {
	return p.errors
}

// Add an error to errors when the type of peekToken doesn’t match the
// expectation.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// Helper method that advances both curToken and peekToken.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram starts the parsing process and is the entry point for all other
// sub-parsers that are responsible for other nodes in the AST.
func (p *Parser) ParseProgram() *ast.Program {
	// Construct the root node of the AST.
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Iterate over every token in the input until it encounters an token.EOF
	// token.
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// Parse a statement.
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Constructs an *ast.LetStatement node with the token it’s currently
	// sitting on (a token.LET token).
	stmt := &ast.LetStatement{Token: p.curToken}

	// Advances the tokens while making assertions about the next token.
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Use token.IDENT token to construct an *ast.Identifier node.
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Expects an equal sign and jumps over the expression following the
	// equal sign.
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// "assertion functions".
// Enforce the correctness of the order of tokens by checking the type of the
// next token.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
