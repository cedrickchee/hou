package evaluator

// Package evaluator implements the evaluator -- a tree-walker implementation
// that recursively walks the parsed AST (Abstract Syntax Tree) and evaluates
// the nodes according to their semantic meaning.

import (
	"fmt"

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
func Eval(node ast.Node, env *object.Environment) object.Object {
	// Traverse the AST by starting at the top of the tree, receiving an
	// *ast.Program, and then traverse every node in it.
	// Use object.Environment and keep track of the environment by passing it
	// around.

	switch node := node.(type) {

	// Statements
	case *ast.Program:
		// Traverse the tree and evaluate every statement of the *ast.Program.
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		// If the statement is an *ast.ExpressionStatement we evaluate its
		// expression. An expression statement (not a return statement and not
		// a let statement).
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		// Evaluate the expression associated with the return statement.
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		// Keep track of values using Environment.
		env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		// The first step is to evaluate its operand and then use the result of
		// this evaluation with the operator.
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		// We just reuse the Parameters and Body fields of the AST node.
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		// Using Eval to get the function we want to call.
		// Whether that's an *ast.Identifier or an *ast.FunctionLiteral: Eval
		// returns an *object.Function.
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		// Evaluate the arguments of a call expression.
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		// Call the function. Apply the function to the arguments.
		return applyFunction(function, args)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	// evalProgram was renamed from evalStatements and make less generic because
	// we can’t reuse evalStatements function for evaluating block statements.
	// We are using evalBlockStatement for evaluating block statements.

	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			// Check if the last evaluation result is such an object.ReturnValue
			// and if so, we stop the evaluation and return the unwrapped value.
			return result.Value
		case *object.Error:
			// Error handling — stop the evaluation.
			return result
		}
	}

	return result
}

func evalBlockStatement(
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	// Evaluate an *ast.BlockStatement.

	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		// Here we explicitly don't unwrap the return value and only check the
		// Type() of each evaluation result. If it's object.RETURN_VALUE_OBJ we
		// simply return the *object.ReturnValue, without unwrapping its .Value,
		// so it stops execution in a possible outer block statement and bubbles
		// up to evalProgram, where it finally get's unwrapped.
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
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
		// If the operator is not supported we don't return NULL since we now
		// have error handling implemented.
		return newError("unknown operator: %s%s", operator, right.Type())
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
		return newError("unknown operator: -%s", right.Type())
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
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	// Deciding what to evaluate.

	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}

	return val
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

func newError(format string, a ...interface{}) *object.Error {
	// Helper function to help create new Error type.
	// Error type wraps the formatted error messages.
	//
	// This function finds its use in every place where we didn't know what to
	// do before and returned NULL instead.
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	// Evaluating the arguments is nothing more than evaluating a list of
	// expressions and keeping track of the produced values. But we also
	// have to stop the evaluation process as soon as it encounters an
	// error.

	var result []object.Object

	// This part is where we decided to evaluate the arguments from
	// left-to-right.
	for _, e := range exps {
		// Evaluate ast.Expression in the context of the current environment.
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	// Convert the fn parameter to a *object.Function reference.
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	// Creates a new *object.Environment that's enclosed by the function's
	// environment.
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		// In this new, enclosed environment, binds the arguments of the
		// function call to the function's parameter names.
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	// The result of the evaluation (obj) is unwrapped if it's an
	// *object.ReturnValue. That’s necessary, because otherwise a return s
	// tatement would bubble up through several functions and stop the
	// evaluation in all of them. But we only want to stop the evaluation of the
	// last called function's body.
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}
