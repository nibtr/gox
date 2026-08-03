package runtime

import (
	"fmt"

	"github.com/nibtr/gox/lexer"
)

type Environment struct {
	values map[string]any
	// immediate parent scope
	enclosing *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func NewEnvironmentWithEnclosing(enclosing *Environment) *Environment {
	return &Environment{
		values:    make(map[string]any),
		enclosing: enclosing,
	}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) assign(name lexer.Token, value any) error {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}

	return &RuntimeError{
		Token:   &name,
		Message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}

func (e *Environment) assignAt(distance int, name lexer.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e *Environment) get(name lexer.Token) (any, error) {
	v, ok := e.values[name.Lexeme]
	if ok {
		return v, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, &RuntimeError{
		Token:   &name,
		Message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
}

func (e *Environment) getAt(distance int, name string) (any, error) {
	// no need check for undefined variable because we know
	// it will be because the resolver already found it before
	return e.ancestor(distance).values[name], nil
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}
