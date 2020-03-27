package parser

// Package parser implements the parser that takes as input tokens from the
// lexer and produces as output an AST (Abstract Syntax Tree).

import (
	"fmt"
	"strconv"

	"github.com/cedrickchee/hou/ast"
	"github.com/cedrickchee/hou/lexer"
	"github.com/cedrickchee/hou/token"
)

// Define the precedences of the language.
// These constants is able to answer: "does the * operator have a higher
// precedence than the == operator? Does a prefix operator have a higher
// preference than a call expression?"
const (
	_           int = iota
	LOWEST          // lowest possible precedence
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	PREFIX          // -X or !X
	CALL            // myFunction(X)
)

// Precedence table.
// It associates token types with their precedence.
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

// Pratt parser's idea is the association of parsing functions with token types.
// Whenever this token type is encountered, the parsing functions are called to
// parse the appropriate expression and return an AST node that represents it.
// Each token type can have up to two parsing functions associated with it,
// depending on whether the token is found in a prefix or an infix position.
type (
	prefixParseFn func() ast.Expression
	// This function argument is "left side" of the infix operator that’s being
	// parsed.
	infixParseFn func(ast.Expression) ast.Expression
)

// Parser implements the parser.
type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token

	// maps to get the correct prefixParseFn or infixParseFn for the current
	// token type.
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New constructs a new Parser with a Lexer as input.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Initialize the prefixParseFns map.
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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
		return p.parseExpressionStatement()
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

// The top-level method that kicks off expression parsing.
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Check whether there's a parsing function associated with p.curToken.Type in
// the prefix position.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		// noPrefixParseFnError give us better error messages when
		// program.Statements does not contain one statement but simply one nil.
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	// The heart of our Pratt parser.
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// Try to find infixParseFns for the next token.
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)

		// Loop until it encounters a token that has a higher precedence.
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	// This method doesn’t advance the tokens, it doesn’t call nextToken.
	// That’s important.
	// All of our parsing functions, prefixParseFn or infixParseFn, are going to
	// follow this protocol:
	// start with curToken being the type of token you’re associated with and
	// return with curToken being the last token that’s part of your expression
	// type. Never advance the tokens too far.
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	// Advances our tokens in order to correctly parse a prefix expression
	// like `-5` more than one token has to be "consumed".
	p.nextToken()

	// parseExpression() value changes depending on the caller's knowledge and
	// its context.
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken, // the operator of the infix expression
		Operator: p.curToken.Literal,
		Left:     left,
	}

	// Precedence of the operator token.
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
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

// Helper method that add entries to the prefixParseFns map.
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// Helper method that add entries to the infixParseFns map.
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Returns the precedence associated with the token type of peekToken.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// Returns the precedence associated with the token type of curToken.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}
