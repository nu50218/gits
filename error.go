package main

import "fmt"

type ErrorType int

const (
	ErrorTypeGeneral = iota
	ErrorTypeInvalidArgument
)

type Error struct {
	Type         ErrorType
	Format       string
	WrappedError error
}

func (e *Error) Unwrap() error {
	return e.WrappedError
}

func (e *Error) Error() string {
	return fmt.Sprintf(e.Format, e.WrappedError)
}

func NewError(errorType ErrorType, format string, wrapedError error) error {
	return &Error{
		Type:         errorType,
		Format:       format,
		WrappedError: wrapedError,
	}
}
