package lexer

import "github.com/cedrickchee/hou/token"

// Package lexer implements the lexical analysis that is used to transform the
// source code input into a stream of tokens for parsing by the parser.
// The lexer only supports ASCII characters instead of the full Unicode range
// for now to keep things simple.

// Lexer represents the lexer and contains the source input and internal state.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New returns a new Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken returns the next token read from the input stream.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			// Note: save l.ch in a local variable before calling l.readChar()
			// again. This way we don’t lose the current character and can
			// safely advance the lexer so it leaves the NextToken() with
			// l.position and l.readPosition in the correct state.
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// Early exit here. We don't need the call to readChar() below.
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// Helper method to make the usage of these lexer fields easier to understand.
// It gives us the next character and advance our position in the input string.
func (l *Lexer) readChar() {
	// First, check whether we've reached the end of input.
	if l.readPosition >= len(l.input) {
		// 0 is the ASCII code for the "NUL" character and signifies either
		// "we haven't read anything yet" or "end of file".
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	// After that, l.readPosition always point to the next position where we're
	// going to read from next and l.position always points to the position
	// where we last read.
	l.position = l.readPosition
	l.readPosition++

	// Note: Unicode support
	// ---------------------
	// In order to fully support Unicode and UTF-8 we would need to change l.ch
	// from a byte to rune and change the way we read the next characters,
	// since they could be multiple bytes wide now.
}

// peekChar is similar to readChar except that it doesn’t increment l.position
// and l.readPosition.
// We only want to “peek” ahead in the input and not move around in it, so we
// know what a call to readChar would return.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// Reads in an identifier and advances our lexer’s positions until it encounters
// a non-letter-character.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		// Call readChar until it encounters either a closing double quote or
		// the end of the input.
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// In Monkey whitespace only acts as a separator of tokens and doesn’t have
// meaning, so we need to skip over it entirely.
// Otherwise, we get an ILLEGAL token for the whitespace character. Example,
// between “let five”.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Helper function just checks whether the given argument is a letter.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit returns whether the passed in byte is a Latin digit between 0 and 9.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
