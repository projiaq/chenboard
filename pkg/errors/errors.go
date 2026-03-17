package errors

import "fmt"

// ErrorCode 错误码
type ErrorCode int

const (
	// 成功
	CodeSuccess ErrorCode = 0

	// 参数错误 1000-1999
	CodeInvalidParams   ErrorCode = 1000
	CodeResourceNotFound ErrorCode = 1001
	CodeResourceExists   ErrorCode = 1002

	// 认证错误 2000-2999
	CodeUnauthorized    ErrorCode = 2000
	CodeForbidden       ErrorCode = 2001
	CodeTokenExpired    ErrorCode = 2002
	CodeInvalidToken    ErrorCode = 2003

	// 内部错误 3000-3999
	CodeInternalError  ErrorCode = 3000
	CodeDatabaseError  ErrorCode = 3001
	CodeExternalError  ErrorCode = 3002
)

// AppError 应用错误
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 实现 errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误
var (
	ErrInvalidParams    = New(CodeInvalidParams, "Invalid parameters")
	ErrResourceNotFound = New(CodeResourceNotFound, "Resource not found")
	ErrResourceExists   = New(CodeResourceExists, "Resource already exists")
	ErrUnauthorized     = New(CodeUnauthorized, "Unauthorized")
	ErrForbidden        = New(CodeForbidden, "Forbidden")
	ErrTokenExpired     = New(CodeTokenExpired, "Token expired")
	ErrInvalidToken     = New(CodeInvalidToken, "Invalid token")
	ErrInternalError    = New(CodeInternalError, "Internal server error")
	ErrDatabaseError    = New(CodeDatabaseError, "Database error")
)
