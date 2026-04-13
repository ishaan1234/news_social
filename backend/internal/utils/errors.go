package utils

import "errors"

// Common reusable errors
var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInternal       = errors.New("internal server error")
)

// Mapping errors to HTTP status codes
func GetStatusCode(err error) int {
	switch err {
	case ErrInvalidInput:
		return 400
	case ErrUnauthorized:
		return 401
	case ErrForbidden:
		return 403
	case ErrNotFound:
		return 404
	default:
		return 500
	}
}