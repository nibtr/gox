package runtime

import (
	"fmt"

	"github.com/nibtr/gox/lexer"
)

type Class struct {
	Name string
}

func (f *Class) String() string {
	return f.Name
}

func (f *Class) Call(i *Interpreter, args []any) (any, error) {
	instance := &ClassInstance{klass: f, fields: make(map[string]any)}
	return instance, nil
}

func (f *Class) Arity() int {
	return 0
}

type ClassInstance struct {
	klass  *Class
	fields map[string]any
}

func (ci *ClassInstance) String() string {
	return fmt.Sprintf("%v instance", ci.klass.Name)
}

func (ci *ClassInstance) Get(name lexer.Token) (any, error) {
	if v, ok := ci.fields[name.Lexeme]; ok {
		return v, nil
	}

	return nil, &RuntimeError{
		Token:   &name,
		Message: fmt.Sprintf("Undefined property '%v'.", name.Lexeme),
	}
}
