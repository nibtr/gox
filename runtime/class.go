package runtime

type Class struct {
	Name string
}

func (f *Class) String() string {
	return f.Name
}
