package backend

import "fmt"

// Errors handles errors list for hadnling arrays
type Errors struct {
	errors []error
}

// Error returns the error string
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

// Add adds an error to the list
func (e *Errors) Add(format string, args ...interface{}) {
	e.errors = append(e.errors, fmt.Errorf(format, args...))
}
