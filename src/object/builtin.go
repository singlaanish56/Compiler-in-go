package object

import "fmt"

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {

			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to len not supported, got %s", args[0].Type())
			}

		},
		},
	},
	{
		"puts",
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return nil
		},
		},
	},
	{
		"first",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to the first should be an array, got %s", args[0].Type())
			}

			arg := args[0].(*Array)
			if len(arg.Elements) > 0 {
				return arg.Elements[0]
			}

			return nil
		},
		},
	},
	{
		"last",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to the first should be an array, got %s", args[0].Type())
			}

			arg := args[0].(*Array)
			length := len(arg.Elements)
			if length > 0 {
				return arg.Elements[length-1]
			}

			return nil
		},
		},
	},
	{
		"rest",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to the first should be an array, got %s", args[0].Type())
			}

			arg := args[0].(*Array)
			length := len(arg.Elements)
			if length > 0 {
				newElements := make([]Object, length-1)
				copy(newElements, arg.Elements[1:length])
				return &Array{Elements: newElements}
			}

			return nil
		},
		},
	},
	{
		"push",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != ARRAY_OBJ {
				return newError("argument to the first should be an array, got %s", args[0].Type())
			}

			arg := args[0].(*Array)
			length := len(arg.Elements)

			newElements := make([]Object, length+1)
			copy(newElements, arg.Elements)
			newElements[length] = args[1]

			return nil
		},
		},
	},
}

var builtins = map[string]*Builtin{
	"len":   GetBuiltinByName("len"),
	"puts":  GetBuiltinByName("puts"),
	"first": GetBuiltinByName("first"),
	"rest":  GetBuiltinByName("rest"),
	"push":  GetBuiltinByName("push"),
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}

	return nil
}

func newError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
