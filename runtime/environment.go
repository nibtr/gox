package runtime

import (
	"fmt"

	"github.com/nibtr/gox/lexer"
)

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) assign(name lexer.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; !ok {
		return &RuntimeError{
			Token:   &name,
			Message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
		}
	}

	e.values[name.Lexeme] = value
	return nil
}

func (e *Environment) get(name lexer.Token) (any, error) {
	v, ok := e.values[name.Lexeme]
	if ok {
		return v, nil
	}
	return nil, &RuntimeError{
		Token:   &name,
		Message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}
