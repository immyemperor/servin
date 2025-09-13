package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
	// System errors
	ErrTypeSystem     ErrorType = "SYSTEM"
	ErrTypeIO         ErrorType = "IO"
	ErrTypeNetwork    ErrorType = "NETWORK"
	ErrTypePermission ErrorType = "PERMISSION"

	// Container errors
	ErrTypeContainer ErrorType = "CONTAINER"
	ErrTypeImage     ErrorType = "IMAGE"
	ErrTypeVolume    ErrorType = "VOLUME"
	ErrTypeNamespace ErrorType = "NAMESPACE"
	ErrTypeCgroup    ErrorType = "CGROUP"

	// User errors
	ErrTypeValidation ErrorType = "VALIDATION"
	ErrTypeNotFound   ErrorType = "NOT_FOUND"
	ErrTypeConflict   ErrorType = "CONFLICT"
	ErrTypeConfig     ErrorType = "CONFIG"
)

// ServinError represents a structured error with context
type ServinError struct {
	Type      ErrorType
	Operation string
	Message   string
	Cause     error
	Context   map[string]interface{}
	Stack     string
}

// Error implements the error interface
func (e *ServinError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %s (caused by: %v)", e.Type, e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Operation, e.Message)
}

// Unwrap returns the underlying error
func (e *ServinError) Unwrap() error {
	return e.Cause
}

// GetContext returns the error context
func (e *ServinError) GetContext() map[string]interface{} {
	return e.Context
}

// GetStack returns the stack trace
func (e *ServinError) GetStack() string {
	return e.Stack
}

// NewError creates a new ServinError
func NewError(errType ErrorType, operation, message string) *ServinError {
	return &ServinError{
		Type:      errType,
		Operation: operation,
		Message:   message,
		Context:   make(map[string]interface{}),
		Stack:     getStackTrace(),
	}
}

// WrapError wraps an existing error with additional context
func WrapError(err error, errType ErrorType, operation, message string) *ServinError {
	return &ServinError{
		Type:      errType,
		Operation: operation,
		Message:   message,
		Cause:     err,
		Context:   make(map[string]interface{}),
		Stack:     getStackTrace(),
	}
}

// WithContext adds context to the error
func (e *ServinError) WithContext(key string, value interface{}) *ServinError {
	e.Context[key] = value
	return e
}

// WithContextMap adds multiple context values
func (e *ServinError) WithContextMap(ctx map[string]interface{}) *ServinError {
	for k, v := range ctx {
		e.Context[k] = v
	}
	return e
}

// getStackTrace captures the current stack trace
func getStackTrace() string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])

	var builder strings.Builder
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()
		builder.WriteString(fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	return builder.String()
}

// IsType checks if an error is of a specific type
func IsType(err error, errType ErrorType) bool {
	if servinErr, ok := err.(*ServinError); ok {
		return servinErr.Type == errType
	}
	return false
}

// GetType returns the error type if it's a ServinError
func GetType(err error) (ErrorType, bool) {
	if servinErr, ok := err.(*ServinError); ok {
		return servinErr.Type, true
	}
	return "", false
}

// Common error constructors

// NewSystemError creates a system-level error
func NewSystemError(operation, message string) *ServinError {
	return NewError(ErrTypeSystem, operation, message)
}

// NewIOError creates an I/O error
func NewIOError(operation, message string) *ServinError {
	return NewError(ErrTypeIO, operation, message)
}

// NewNetworkError creates a network error
func NewNetworkError(operation, message string) *ServinError {
	return NewError(ErrTypeNetwork, operation, message)
}

// NewPermissionError creates a permission error
func NewPermissionError(operation, message string) *ServinError {
	return NewError(ErrTypePermission, operation, message)
}

// NewContainerError creates a container error
func NewContainerError(operation, message string) *ServinError {
	return NewError(ErrTypeContainer, operation, message)
}

// NewImageError creates an image error
func NewImageError(operation, message string) *ServinError {
	return NewError(ErrTypeImage, operation, message)
}

// NewVolumeError creates a volume error
func NewVolumeError(operation, message string) *ServinError {
	return NewError(ErrTypeVolume, operation, message)
}

// NewValidationError creates a validation error
func NewValidationError(operation, message string) *ServinError {
	return NewError(ErrTypeValidation, operation, message)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(operation, message string) *ServinError {
	return NewError(ErrTypeNotFound, operation, message)
}

// NewConflictError creates a conflict error
func NewConflictError(operation, message string) *ServinError {
	return NewError(ErrTypeConflict, operation, message)
}

// NewConfigError creates a configuration error
func NewConfigError(operation, message string) *ServinError {
	return NewError(ErrTypeConfig, operation, message)
}
