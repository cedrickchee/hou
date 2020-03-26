package token

// Package token defines the tokens our lexer is going to output.

// There is a limited number of different token types in the Monkey language.
// That means we can define the possible TokenTypes as constants.
const (
	// special type that signifies a token/character we don't know about
	ILLEGAL = "ILLEGAL"
	// special type stands for "end of file", which tells parser that it can stop
	EOF = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 1343456

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// TokenType distinguishes between different types of tokens.
type TokenType string

// Token holds a single token type and its literal value.
type Token struct {
	Type    TokenType
	Literal string
}
