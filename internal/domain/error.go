package domain

import "fmt"

// ErrCode is an error code that maps to an HTTP status
type ErrCode int

const (
	ErrCodeBadRequest   ErrCode = 400
	ErrCodeNotFound     ErrCode = 404
	ErrCodeUpstreamFail ErrCode = 502
)

// AppError is the unified application-layer error type.
// Handlers can use errors.As to inspect it without parsing strings.
type AppError struct {
	Code    ErrCode
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewBadRequest(msg string) *AppError {
	return &AppError{Code: ErrCodeBadRequest, Message: msg}
}

func NewUpstreamFail(msg string) *AppError {
	return &AppError{Code: ErrCodeUpstreamFail, Message: msg}
}
