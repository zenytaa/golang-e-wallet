package apperror

import "runtime/debug"

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

type AppError struct {
	stack []byte
}

func NewAppError() *AppError {
	return &AppError{
		stack: debug.Stack(),
	}
}

func (e *AppError) Error() string {
	return "got error!"
}

func (e *AppError) GetStackTrace() []byte {
	return e.stack
}
