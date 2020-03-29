package evaluator

// Package evaluator implements the evaluator -- a tree-walker implementation
// that recursively walks the parsed AST (Abstract Syntax Tree) and evaluates
// the nodes according to their semantic meaning.

import (
	"github.com/cedrickchee/hou/ast"
	"github.com/cedrickchee/hou/object"
)

var (
	// TRUE is a cached Boolean object holding the `true` value.
	TRUE = &object.Boolean{Value: true}

	// FALSE is a cached Boolean object holding the `false` value.
	FALSE = &object.Boolean{Value: false}

	// NULL is a cached Null object. There should only be one reference to a
	// null value, just as there's only one 'true' and one 'false'.
	// No kinda-but-not-quite-null, no half-null and no
	// basically-thesame-as-the-other-null.
	NULL = &object.Null{}
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

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
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

// Helper function to reference true or false to only two instances of
// object.Boolean: TRUE and FALSE.
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	// We shouldn't create a new object.Boolean every time we encounter a true
	// or false. There is no difference between two trues. The same goes for
	// false. We shouldn't use new instances every time. There are only two
	// possible values, so let's reference them instead of allocating new
	// object.Booleans (creating new ones). That is a small performance
	// improvement we get without a lot of work.

	if input {
		return TRUE
	}
	return FALSE
}
