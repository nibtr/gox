package runtime

import (
	"fmt"
	"time"

	"github.com/nibtr/gox/ast"
)

// Clock returns the current time in seconds since the Unix epoch.
type Clock struct{}

func (f *Clock) Arity() int {
	return 0
}

func (f *Clock) Call(i *interpreter, args []any) (any, error) {
	return float64(time.Now().Unix()), nil
}

func (f *Clock) String() string {
	return "<native fn>"
}

// Function represents a general function
type Function struct {
	declaration *ast.FunctionStmt
	closure     *Environment
}

func (f *Function) Call(i *interpreter, args []any) (any, error) {
	env := NewEnvironmentWithEnclosing(f.closure)
	for i := range f.declaration.Params {
		// safe to assume assume the parameter and argument lists have the same length
		// visitCallExpr() checks the arity before calling call()
		env.define(f.declaration.Params[i].Lexeme, args[i])
	}

	err := i.executeBlock(f.declaration.Body, env)
	if ret, ok := err.(*Return); ok {
		return ret.Value, nil
	}

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (f *Function) Arity() int {
	return len(f.declaration.Params)
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}
