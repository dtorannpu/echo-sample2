package errors

import "fmt"

type ErrInvalidTodoID struct {
	Value string
	Err   error
}

func (e *ErrInvalidTodoID) Error() string {
	return fmt.Sprintf("invalid todo id %q: %v", e.Value, e.Err)
}

func (e *ErrInvalidTodoID) Unwrap() error {
	return e.Err
}
