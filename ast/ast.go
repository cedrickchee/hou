package ast

// Packge ast implement the Abstract Syntax Tree (AST) that represents the
// parsed source code before being passed on to the interpreter for evaluation.

import "github.com/cedrickchee/hou/token"

// Node defines an interface for all nodes in the AST.
type Node interface {
	// Returns the literal value of the token it's associated with.
	// This method will be used only for debugging and testing.
	TokenLiteral() string
}

// Statement defines the interface for all statement nodes.
type Statement interface {
	// Some of these nodes implement the Statement interface.
	Node
	statementNode()
}

// Expression defines the interface for all expression nodes.
type Expression interface {
	// Some of these nodes implement the Expression interface.
	Node
	expressionNode()
}

// =============================================================================
// Implementation of Node
// =============================================================================

// Program is the root node of every AST. Every valid program is a series of
// statements.
type Program struct {
	// A program consists of a slice of AST nodes that implement the Statement
	// interface.
	Statements []Statement
}

// TokenLiteral prints the literal value of the token associated with this node.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement the `let` statement represents the AST node that binds an
// expression to an identifier
type LetStatement struct {
	Token token.Token // the token.LET token
	// Name hold the identifier of the binding and Value for the expression
	// that produces the value.
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Identifier is a node that holds the literal value of an identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// To hold the identifier of the binding, the x in let x = 5; , we have the
// Identifier struct type, which implements the Expression interface.
func (i *Identifier) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
