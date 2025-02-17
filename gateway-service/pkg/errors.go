package pkg

import (
	"errors"
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
	Code    string
	Message string
}

func Errorf(code string, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func ErrorCode(err error) string {
	var e *Error

	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}

	return INTERNAL_ERROR
}

func ErrorMessage(err error) string {
	var e *Error

	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}

	return "Internal error."
}

type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func ErrorResponse(msg string, code int) APIError {
	return APIError{
		StatusCode: code,
		Message:    msg,
	}
	// return gin.H{
	// 	"message": msg,
	// 	"error":   err.Error(),
	// }
}

// Error implements the error interface. Not used by the application otherwise.
func (e *Error) Error() string {
	return fmt.Sprintf("error: code=%s message=%s", e.Code, e.Message)
}
