package repl

// Package repl implements the Read-Eval-Print-Loop (REPL) or interactive console
// by lexing, parsing and evaluating the input in the interpreter.

import (
	"bufio"
	"fmt"
	"io"

	"github.com/cedrickchee/hou/evaluator"
	"github.com/cedrickchee/hou/lexer"
	"github.com/cedrickchee/hou/object"
	"github.com/cedrickchee/hou/parser"
)

// PROMPT is the REPL prompt displayed for each input.
const PROMPT = ">> "

// MONKEYFACE is the REPL's face if we run into any parser errors. You get to
// see a monkey :D
const MONKEYFACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

// Start starts the REPL in a continuous loop.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// A REPL that tokenizes and parses Monkey source code and prints
		// the AST.
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			// Print string representation of the object to stdout.
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

// Print parser errors to stdout.
func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEYFACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
