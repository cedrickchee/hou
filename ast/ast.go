package ast

// Packge ast implement the Abstract Syntax Tree (AST) that represents the
// parsed source code before being passed on to the interpreter for evaluation.

import (
	"bytes"
	"strings"

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

// String returns a stringified version of the AST `let` node for debugging.
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
// holds a return value to the outer stack in the call stack.
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns a stringified version of the AST `return` node for debugging.
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement represents an expression statement and holds an
// expression.
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns a stringified version of the AST for debugging.
func (es *ExpressionStatement) String() string {
	// The nil-checks will be taken out, later on, when we can fully build
	// expressions.
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral represents a literal integer and holds an integer value.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns a stringified version of the AST for debugging.
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// PrefixExpression represents a prefix expression and holds the operator as
// as well as the right-hand side expression.
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns a stringified version of the AST for debugging.
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

// InfixExpression represents an infix expression and holds the left-hand
// expression, operator and right-hand expression.
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }

// String returns a stringified version of the AST for debugging.
func (oe *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}

// Boolean represents a boolean value and holds the underlying boolean value.
type Boolean struct {
	Token token.Token
	Value bool // save either true or false.
}

func (b *Boolean) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// String returns a stringified version of the AST for debugging.
func (b *Boolean) String() string { return b.Token.Literal }

// IfExpression represents an `if` expression and holds the condition,
// consequence and alternative expressions
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a stringified version of the AST for debugging.
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// BlockStatement represents a block statement and holds a series of statements.
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// String returns a stringified version of the AST for debugging.
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// FunctionLiteral represents a literal function and has two main parts,
// the list of parameters and the block statement that is the function's body.
type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

// The type of AST node for FunctionLiteral is expression.
func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns a stringified version of the AST for debugging.
func (fl *FunctionLiteral) String() string {
	// The abstract structure of a function literal is:
	// 		fn <parameters> <block statement>
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpression represents a call expression and holds the function to be
// called as well as the arguments to be passed to that function.
type CallExpression struct {
	Token token.Token // The '(' token
	// Identifier or FunctionLiteral.
	// FunctionLiteral is the 'fn' token.
	// Identifier is the function name, example 'add'.
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// String returns a stringified version of the AST for debugging.
func (ce *CallExpression) String() string {
	// Call expression structure:
	// 		<expression>(<comma separated expressions>)

	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// StringLiteral represents a literal string and holds a string value.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node.
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

// String returns a stringified version of the AST for debugging.
func (sl *StringLiteral) String() string { return sl.Token.Literal }
