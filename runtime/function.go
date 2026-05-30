package runtime

import "time"

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
