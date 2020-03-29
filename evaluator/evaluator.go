package evaluator

// Package evaluator implements the evaluator -- a tree-walker implementation
// that recursively walks the parsed AST (Abstract Syntax Tree) and evaluates
// the nodes according to their semantic meaning.

import (
	"github.com/cedrickchee/hou/ast"
	"github.com/cedrickchee/hou/object"
)

// Eval evaluates the node and returns an object.
func Eval(node ast.Node) object.Object {
	// Traverse the AST by starting at the top of the tree, receiving an
	// *ast.Program, and then traverse every node in it.

	switch node := node.(type) {

	// Statements
	case *ast.Program:
		// Traverse the tree and evaluate every statement of the *ast.Program.
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		// If the statement is an *ast.ExpressionStatement we evaluate its
		// expression. An expression statement (not a return statement and not
		// a let statement).
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}
