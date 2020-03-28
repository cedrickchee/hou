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

// Precedence table for infix expression.
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
	token.LPAREN:   CALL,
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
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpression)

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

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	// Take care of optional semicolons.
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// The top-level method that kicks off expression parsing.
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// Take care of optional semicolons.
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

func (p *Parser) parseBoolean() ast.Expression {
	// The structure of our parser serves us well.
	// That actually is one of the beauties of Pratt's approach: it's so easy
	// to extend.

	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	// In no other parsing function did we use expectPeek so extensively.
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// This method also follows our parsing function protocol: the tokens get
	// advanced just enough so that parseBlockStatement sits on the `{` with
	// p.curToken being of type token.LBRACE.
	expression.Consequence = p.parseBlockStatement()

	// This support the else part of an if-else-condition. It check if it even
	// exists and if so, parse the block statement that comes directly after
	// the else.
	// Allows an optional 'else' but doesn’t add a parser error if there is none.
	if p.peekTokenIs(token.ELSE) {
		// If we have a token.ELSE, we advance the tokens 2 times.
		// The first time with a call to nextToken, since we already know that
		// the p.peekToken is the 'else'. Then with a call to expectPeek since
		// now the next token has to be the opening brace of a block statement.
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	// Calls parseStatement until it encounters either a `}`, which signifies
	// the end of the block statement, or a token.EOF, which tells us that
	// there’s no more tokens left to parse.
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	// One of the great things about our parser is that once we define function
	// literals as expressions and provide a function to correctly parse them
	// the rest works.

	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	// Method to parse the literal's parameters.

	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// Constructs the slice of parameters by repeatedly building identifiers
	// from the comma separated list. It also makes an early exit if the list is
	// empty and it carefully handles lists of varying sizes.
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers

	// For a method like this it really pays off to have another set of tests
	// that check the edge cases: an empty parameter list, a list with one
	// parameter and a list with multiple parameters.
	// Please see TestFunctionParameterParsing.
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

// Parse the function's argument list.
func (p *Parser) parseCallArguments() []ast.Expression {
	// This method looks strikingly similar to parseFunctionParameters, except
	// that it's more generic and returns a slice of ast.Expression and not
	// *ast.Identifier (because call expression AST structure is:
	// <expression>(<comma separated expressions>))

	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
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
