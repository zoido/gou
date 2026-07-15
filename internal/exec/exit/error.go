package exit

import (
	"errors"
	"fmt"
)

// New creates a new [Error] with the given exit code and message.
func New(code int, msg string) Error {
	return Error{
		err:  errors.New(msg),
		code: code,
	}
}

// Newf creates a new [Error] with the given exit code and formatted message.
func Newf(code int, format string, args ...any) Error {
	return Error{
		err:  fmt.Errorf(format, args...),
		code: code,
	}
}

// Error represents an error with an associated exit code.
type Error struct {
	err  error
	code int
}

func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) Unwrap() error {
	return e.err
}

// ExitCode returns the exit code associated with the error.
func (e Error) ExitCode() int {
	return e.code
}
