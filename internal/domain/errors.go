package domain

import "errors"

var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInactiveUser       = errors.New("user is not active")

	// OTP errors
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
	ErrOTPAlreadyVerified = errors.New("OTP already verified")

	// Authentication errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")

	// Validation errors
	ErrInvalidEmail = errors.New("invalid email format")
	ErrWeakPassword = errors.New("password does not meet requirements")

	// Cache errors
	ErrCacheMiss        = errors.New("cache miss")
	ErrCacheUnavailable = errors.New("cache service unavailable")
)

// ValidationError represents a validation error with details
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
