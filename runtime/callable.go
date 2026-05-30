package runtime

type callable interface {
	// arity() returns the number of arguments a function or operation expects.
	arity() int
	// call calls the function or operation with a list of arguments
	call(i *interpreter, arguments []any) any
}
