package ast

// Packge ast implement the Abstract Syntax Tree (AST) that represents the
// parsed source code before being passed on to the interpreter for evaluation.

import (
	"bytes"

	"github.com/cedrickchee/hou/token"
)

// Node defines an interface for all nodes in the AST.
type Node interface {
	// Returns the literal value of the token it's associated with.
	// This method will be used only for debugging and testing.
	TokenLiteral() string
	// Returns a stringified version of the AST for debugging.
	String() string
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

// String returns a stringified version of the AST for debugging.
func (p *Program) String() string {
	// Creates a buffer and writes the return value of each statements String()
	// method to it.
	var out bytes.Buffer

	for _, s := range p.Statements {
		// Delegates most of program work to the Statements of *ast.Program.
		out.WriteString(s.String())
	}

	return out.String()
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

// String returns a stringified version of the `let` node.
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

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

// String returns a stringified version of the identifier node.
func (i *Identifier) String() string {
	return i.Value
}

// ReturnStatement the `return` statement that represents the AST node that
// holds a return value to the outter stack in the call stack.
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns a stringified version of the `return` node.
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement represents an expression node.
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns a stringified version of the expression node
func (es *ExpressionStatement) String() string {
	// The nil-checks will be taken out, later on, when we can fully build
	// expressions.
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral represents a literal integer node.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns a stringified version of the expression node.
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// PrefixExpression represents a prefix expression node.
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns a stringified version of the expression node.
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	// We deliberately add parentheses around the operator and its operand,
	// the expression in Right. That allows us to see which operands belong to
	// which operator.
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents an infix expression node.
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }

// String returns a stringified version of the expression node.
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}
