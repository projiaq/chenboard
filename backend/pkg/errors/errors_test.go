package errors

import (
	"errors"
	"testing"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		want string
	}{
		{
			name: "without wrapped error",
			err:  New(CodeInvalidParams, "test error"),
			want: "[1000] test error",
		},
		{
			name: "with wrapped error",
			err:  Wrap(CodeDatabaseError, "db error", errors.New("connection failed")),
			want: "[3001] db error: connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("AppError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	err := Wrap(CodeInternalError, "outer error", innerErr)

	if unwrapped := err.Unwrap(); unwrapped != innerErr {
		t.Errorf("AppError.Unwrap() = %v, want %v", unwrapped, innerErr)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name string
		err  *AppError
		code ErrorCode
	}{
		{"ErrInvalidParams", ErrInvalidParams, CodeInvalidParams},
		{"ErrResourceNotFound", ErrResourceNotFound, CodeResourceNotFound},
		{"ErrUnauthorized", ErrUnauthorized, CodeUnauthorized},
		{"ErrTokenExpired", ErrTokenExpired, CodeTokenExpired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("%s.Code = %d, want %d", tt.name, tt.err.Code, tt.code)
			}
		})
	}
}
