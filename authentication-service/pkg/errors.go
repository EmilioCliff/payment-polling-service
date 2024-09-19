// package pkg

// import (
//     "fmt"
// )

// const (
// 	ALREADY_EXISTS_ERROR  = "already_exists"
// 	INTERNAL_ERROR        = "internal"
// 	INVALID_ERROR         = "invalid"
// 	NOT_FOUND_ERROR       = "not_found"
// 	NOT_IMPLEMENTED_ERROR = "not_implemented"
// 	AUTHENTICATION_ERROR  = "authentication"
// )

// type Error struct {
// 	Code string
// 	Message string
// }

// // Errorf is a helper function to return an Error with a given code and formatted message.
// func Errorf(code string, format string, args ...interface{}) *Error {
// 	return &Error{
// 		Code:    code,
// 		Message: fmt.Sprintf(format, args...),
// 	}
// }

package pkg

import (
	"errors"
	"fmt"
)

// Application error codes.
const (
	ALREADY_EXISTS_ERROR  = "already_exists"
	INTERNAL_ERROR        = "internal"
	INVALID_ERROR         = "invalid"
	NOT_FOUND_ERROR       = "not_found"
	NOT_IMPLEMENTED_ERROR = "not_implemented"
	AUTHENTICATION_ERROR  = "authentication"
)

// Error represents an application-specific error. Application errors can be
// unwrapped by the caller to extract out the code & message.
type Error struct {
	// Machine-readable error code.
	Code string

	// Human-readable error message.
	Message string
}

// Error implements the error interface. Not used by the application otherwise.
func (e *Error) Error() string {
	return fmt.Sprintf("error: code=%s message=%s", e.Code, e.Message)
}

// ErrorCode unwraps an application error and returns its code.
// Non-application errors always return INTERNAL_ERROR.
func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return INTERNAL_ERROR
}

// ErrorMessage unwraps an application error and returns its message.
// Non-application errors always return "Internal error".
func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error."
}

// Errorf is a helper function to return an Error with a given code and formatted message.
func Errorf(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}