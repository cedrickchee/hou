package token

// Package token defines the tokens our lexer is going to output.

// There is a limited number of different token types in the Monkey language.
// That means we can define the possible TokenTypes as constants.
const (
	//
	// Special type
	//
	ILLEGAL = "ILLEGAL" // a token/character we don't know about
	EOF     = "EOF"     // stands for "end of file", which tells parser that it can stop

	//
	// Identifiers + literals
	//
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // an integer, e.g: 1343456
	STRING = "STRING" // a string, e.g: "foobar"

	//
	// Operators
	//
	ASSIGN   = "=" // the assignment operator
	PLUS     = "+" // the addition operator
	MINUS    = "-" // the substraction operator
	BANG     = "!" // the factorial operator
	ASTERISK = "*" // the multiplication operator
	SLASH    = "/" // the division operator

	LT = "<" // the less than comparision operator
	GT = ">" // the greater than comparision operator

	EQ     = "==" // the equality operator
	NOT_EQ = "!=" // the inequality operator

	//
	// Delimiters
	//
	COMMA     = "," // a comma
	SEMICOLON = ";" // a semi-colon
	COLON     = ":" // a colon

	LPAREN   = "(" // a left paranthesis
	RPAREN   = ")" // a right parenthesis
	LBRACE   = "{" // a left brace
	RBRACE   = "}" // a right brace
	LBRACKET = "[" // a left bracket
	RBRACKET = "]" // a right bracket

	//
	// Keywords
	//
	FUNCTION = "FUNCTION" // the `fn` keyword (function)
	LET      = "LET"      // the `let` keyword (let)
	TRUE     = "TRUE"     // the `true` keyword (true)
	FALSE    = "FALSE"    // the `false` keyword (false)
	IF       = "IF"       // the `if` keyword (if)
	ELSE     = "ELSE"     // the `else` keyword (else)
	RETURN   = "RETURN"   // the `return` keyword (return)
)

// Language keywords table
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// TokenType distinguishes between different types of tokens.
type TokenType string

// Token holds a single token type and its literal value.
type Token struct {
	Type    TokenType
	Literal string
}

// LookupIdent looks up the identifier in ident and returns the appropriate
// token type depending on whether the identifier is user-defined or a keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok // language keyword
	}
	return IDENT // user-defined identifier
}
