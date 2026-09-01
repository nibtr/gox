package runtime

import (
	"fmt"

	"github.com/nibtr/gox/lexer"
)

type Class struct {
	Name    string
	Methods map[string]*Function
}

func (f *Class) String() string {
	return f.Name
}

func (f *Class) Call(i *Interpreter, args []any) (any, error) {
	instance := &ClassInstance{Klass: f, Fields: make(map[string]any)}
	return instance, nil
}

func (f *Class) Arity() int {
	return 0
}

func (f *Class) findMethod(name string) *Function {
	if v, ok := f.Methods[name]; ok {
		return v
	}
	return nil
}

type ClassInstance struct {
	Klass  *Class
	Fields map[string]any
}

func (ci *ClassInstance) String() string {
	return fmt.Sprintf("%v instance", ci.Klass.Name)
}

func (ci *ClassInstance) Get(name lexer.Token) (any, error) {
	if v, ok := ci.Fields[name.Lexeme]; ok {
		return v, nil
	}

	method := ci.Klass.findMethod(name.Lexeme)
	if method != nil {
		return method, nil
	}

	return nil, &RuntimeError{
		Token:   &name,
		Message: fmt.Sprintf("Undefined property '%v'.", name.Lexeme),
	}
}

func (ci *ClassInstance) Set(name lexer.Token, value any) {
	ci.Fields[name.Lexeme] = value
}
