package parser

import (
	"fmt"
	"strings"
)

// Parser tracing helps putting tracing statements in the methods of Parser to
// see what was happening when parsing certain expressions.
//
// The file includes two function definitions that are really helpful when
// trying to understand what the parser does: `trace` and `untrace`. Use them
// like this:

/*
parser/parser.go

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    defer untrace(trace("parseExpressionStatement"))
	   // [...]
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    defer untrace(trace("parseExpression"))
	   // [...]
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    defer untrace(trace("parseIntegerLiteral"))
    // [...]
}

func (p *Parser) parsePrefixExpression() ast.Expression {
    defer untrace(trace("parsePrefixExpression"))
    // [...]
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    defer untrace(trace("parseInfixExpression"))
    // [...]
}
*/

// With these tracing statements included we can now use our parser and see what
// it does. Here is the output when parsing the expression statement
// `-1 * 2 + 3` in the test suite:

/*
$ go test -v -run TestOperatorPrecedenceParsing ./parser
=== RUN		TestOperatorPrecedenceParsing
BEGIN parseExpressionStatement
	BEGIN parseExpression
		BEGIN parsePrefixExpression
			BEGIN parseExpression
				BEGIN parseIntegerLiteral
				END parseIntegerLiteral
			END parseExpression
		END parsePrefixExpression
		BEGIN parseInfixExpression
			BEGIN parseExpression
				BEGIN parseIntegerLiteral
				END parseIntegerLiteral
			END parseExpression
		END parseInfixExpression
		BEGIN parseInfixExpression
			BEGIN parseExpression
				BEGIN parseIntegerLiteral
				END parseIntegerLiteral
			END parseExpression
		END parseInfixExpression
	END parseExpression
END parseExpressionStatement
--- PASS: TestOperatorPrecedenceParsing (0.00s)
PASS
*/

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}
