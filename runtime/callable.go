package runtime

type Callable interface {
	// Arity() returns the number of arguments a function or operation expects.
	Arity() int
	// Call calls the function or operation with a list of arguments
	Call(i *interpreter, args []any) (any, error)
}
