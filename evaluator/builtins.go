package evaluator

import (
	"github.com/cedrickchee/hou/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			// Error checking that makes sure that we can't call this function
			// with the wrong number of arguments.
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				// Error checking that makes sure that we can't call this
				// function with an argument of an unsupported type.
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
}
