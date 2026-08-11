package runtime

import "fmt"

type Class struct {
	Name string
}

func (f *Class) String() string {
	return f.Name
}

func (f *Class) Call(i *Interpreter, args []any) (any, error) {
	instance := &ClassInstance{klass: f}
	return instance, nil
}

func (f *Class) Arity() int {
	return 0
}

type ClassInstance struct {
	klass *Class
}

func (ci *ClassInstance) String() string {
	return fmt.Sprintf("%v instance", ci.klass.Name)
}
