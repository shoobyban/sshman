package backend

import "fmt"

type Errors struct {
	errors []error
}

func (e *Errors) Error() string {
	var b []byte
	if e == nil {
		return ""
	}
	for _, v := range e.errors {
		b = append(b, v.Error()...)
		b = append(b, '\n')
	}

	return string(b)
}

func (e *Errors) Add(format string, args ...interface{}) {
	e.errors = append(e.errors, fmt.Errorf(format, args...))
}
