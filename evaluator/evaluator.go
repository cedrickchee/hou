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

	case *ast.BlockStatement:
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

	case *ast.PrefixExpression:
		// The first step is to evaluate its operand and then use the result of
		// this evaluation with the operator.
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node)
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

func evalPrefixExpression(operator string, right object.Object) object.Object {
	// Checks if the operator is supported.
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		// If the operator is not supported we return NULL. Is that the best
		// choice? Maybe, maybe not. For now, it's definitely the easiest
		// choice, since we don't have any error handling implemented yet.
		return NULL
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	// The behavior of the bang ! operator.
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	// Check if the operand is an integer.
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	// Allocate a new object to wrap a negated version of this value.
	return &object.Integer{Value: -value}
}

func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	// Monkey's object system doesn't allow pointer comparison for integer
	// objects. It has to unwrap the value before a comparison can be made.
	// Thus the comparison between booleans is faster.

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		// The check for integer operands has to be higher up in the switch
		// statement.
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		// Using pointer comparison to check for equality between booleans.
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		// Using pointer comparison to check for equality between booleans.
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	// Deciding what to evaluate.

	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
