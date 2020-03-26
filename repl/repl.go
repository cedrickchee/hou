package repl

// Package repl implements the Read-Eval-Print-Loop (REPL) or interactive console
// by lexing, parsing and evaluating the input in the interpreter.

import (
	"bufio"
	"fmt"
	"io"

	"github.com/cedrickchee/hou/lexer"
	"github.com/cedrickchee/hou/token"
)

// PROMPT is the REPL prompt displayed for each input.
const PROMPT = ">> "

// Start starts the REPL in a continuous loop.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// A REPL that tokenizes Monkey source code and prints the tokens.
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
