package object

// Package object implements the object system (or value system) of Monkey used
// to both represent values as the evaluator encounters and constructs them as
// well as how the user interacts with values.

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cedrickchee/hou/ast"
)

const (
	// INTEGER_OBJ is the Integer object type.
	INTEGER_OBJ = "INTEGER"

	// BOOLEAN_OBJ is the Boolean object type.
	BOOLEAN_OBJ = "BOOLEAN"

	// STRING_OBJ is the String object type.
	STRING_OBJ = "STRING"

	// NULL_OBJ is the Null object type.
	NULL_OBJ = "NULL"

	// RETURN_VALUE_OBJ is the Return value object type.
	RETURN_VALUE_OBJ = "RETURN_VALUE"

	// ERROR_OBJ is the Error object type.
	ERROR_OBJ = "ERROR"

	// FUNCTION_OBJ is the Function object type.
	FUNCTION_OBJ = "FUNCTION"
)

// ObjectType represents the type of an object.
type ObjectType string

// Object represents a value and implementations are expected to implement
// `Type()` and `Inspect()` functions. The reason Object being an interface
// instead of struct is that every value needs a different internal
// representation and it’s easier to define two different struct types than
// trying to fit booleans and integers into the same struct field.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer is the integer type used to represent integer literals and holds an
// internal int64 value.
// Whenever we encounter an integer literal in the source code we first turn it
// into an ast.IntegerLiteral and then, when evaluating that AST node, we turn
// it into an object.Integer, saving the value inside our struct and passing
// around a reference to this struct.
type Integer struct {
	Value int64
}

// Type returns the type of the object.
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Inspect returns a stringified version of the object for debugging.
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Boolean is the boolean type and used to represent boolean literals and holds
// an internal bool value.
type Boolean struct {
	Value bool
}

// Type returns the type of the object.
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Inspect returns a stringified version of the object for debugging.
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Null is the null type and used to represent the absence of a value.
type Null struct{}

// Type returns the type of the object.
func (n *Null) Type() ObjectType { return NULL_OBJ }

// Inspect returns a stringified version of the object for debugging.
func (n *Null) Inspect() string { return "null" }

// ReturnValue is the return value type and used to hold the value of another
// object. This is used for `return` statements and this object is tracked
// through the evaluator and when encountered stops evaluation of the program,
// or body of a function.
type ReturnValue struct {
	Value Object
}

// Type returns the type of the object.
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// Inspect returns a stringified version of the object for debugging.
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error is the error type and used to hold a message denoting the details of
// error encountered. This object is tracked through the evaluator and when
// encountered stops evaulation of the program or body of a function.
// In a production-ready interpreter we'd want to attach a stack trace to such
// error objects, add the line and column numbers of its origin.
type Error struct {
	Message string
}

// Type returns the type of the object.
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// Inspect returns a stringified version of the object for debugging.
func (e *Error) Inspect() string { return "ERROR:" + e.Message }

// Function is the function type that holds the function's formal parameters,
// body and an environment to support closures.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type returns the type of the object.
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

// Inspect returns a stringified version of the object for debugging.
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// String is the string type used to represent string literals and holds an
// internal string value.
type String struct {
	Value string
}

// Type returns the type of the object.
func (s *String) Type() ObjectType { return STRING_OBJ }

// Inspect returns a stringified version of the object for debugging.
func (s *String) Inspect() string { return s.Value }
