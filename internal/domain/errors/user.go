package errors

import "fmt"

type DomainError struct {
	Code    string
	Message string
	Field   string
}

func (e *DomainError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s (field: %s)", e.Code, e.Message, e.Field)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// User-specific domain errors
var (
	ErrUserNotFound = &DomainError{
		Code:    "USER_NOT_FOUND",
		Message: "User not found",
	}

	ErrUserAlreadyExists = &DomainError{
		Code:    "USER_ALREADY_EXISTS",
		Message: "User with this email already exists",
		Field:   "email",
	}

	ErrInvalidUserEmail = &DomainError{
		Code:    "INVALID_EMAIL",
		Message: "Invalid email format",
		Field:   "email",
	}

	ErrInvalidUserPassword = &DomainError{
		Code:    "INVALID_PASSWORD",
		Message: "Password does not meet requirements",
		Field:   "password",
	}

	ErrUserInactive = &DomainError{
		Code:    "USER_INACTIVE",
		Message: "User account is inactive",
	}

	ErrUserSuspended = &DomainError{
		Code:    "USER_SUSPENDED",
		Message: "User account is suspended",
	}

	ErrFailedToCheckUserExistance = &DomainError{
		Code:    "FAILED_TO_CHECK_USER_EXISTENCE",
		Message: "failed to check user existence",
	}

	ErrFailedToCreateUser = &DomainError{
		Code:    "FAILED_TO_CREATE_USER",
		Message: "failed to create user",
	}

	ErrFailedToListUsers = &DomainError{
		Code:    "FAILED_TO_LIST_USERS",
		Message: "failed to list users",
	}
)

// Helper functions to create specific errors
func NewUserValidationError(field, message string) *DomainError {
	return &DomainError{
		Code:    "VALIDATION_ERROR",
		Message: message,
		Field:   field,
	}
}
