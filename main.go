package main

// Package main implements the main process which invokes the interpreter's
// REPL and waits for user input before lexing, parsing nad evaulating.

import (
	"fmt"
	"os"
	"os/user"

	"github.com/cedrickchee/hou/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "Hello %s! This is the Hou programming language!\n", user.Username)
	fmt.Fprintf(os.Stdout, "Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
