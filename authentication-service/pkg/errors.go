package pkg

import (
    "fmt"
)

const (
	ALREADY_EXISTS_ERROR  = "already_exists"
	INTERNAL_ERROR        = "internal"
	INVALID_ERROR         = "invalid"
	NOT_FOUND_ERROR       = "not_found"
	NOT_IMPLEMENTED_ERROR = "not_implemented"
	AUTHENTICATION_ERROR  = "authentication"
)

type Error struct {
	Code string
	Message string
}

// Errorf is a helper function to return an Error with a given code and formatted message.
func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}