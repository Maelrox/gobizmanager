package errors

import "net/http"

// ErrorType represents the type of error
type ErrorType string

const (
	// Validation errors
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	// Authentication errors
	ErrorTypeAuthentication ErrorType = "AUTHENTICATION_ERROR"
	// Authorization errors
	ErrorTypeAuthorization ErrorType = "AUTHORIZATION_ERROR"
	// Resource not found errors
	ErrorTypeNotFound ErrorType = "NOT_FOUND_ERROR"
	// Internal server errors
	ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
)

// ErrorCode represents a specific error code
type ErrorCode string

const (
	// Validation error codes
	ErrorCodeInvalidID        ErrorCode = "INVALID_ID"
	ErrorCodeValidationFailed ErrorCode = "VALIDATION_FAILED"

	// Authentication error codes
	ErrorCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrorCodeUserNotFound ErrorCode = "USER_NOT_FOUND"

	// Authorization error codes
	ErrorCodePermissionDenied    ErrorCode = "PERMISSION_DENIED"
	ErrorCodeRootAccessDenied    ErrorCode = "ROOT_ACCESS_DENIED"
	ErrorCodeRootRoleAssignment  ErrorCode = "ROOT_ROLE_ASSIGNMENT"
	ErrorCodeRoleCompanyMismatch ErrorCode = "ROLE_COMPANY_MISMATCH"

	// Resource not found error codes
	ErrorCodePermissionNotFound  ErrorCode = "PERMISSION_NOT_FOUND"
	ErrorCodeRoleNotFound        ErrorCode = "ROLE_NOT_FOUND"
	ErrorCodeCompanyUserNotFound ErrorCode = "COMPANY_USER_NOT_FOUND"

	// Internal error codes
	ErrorCodePermissionCheckFailed ErrorCode = "PERMISSION_CHECK_FAILED"
	ErrorCodeRoleListFailed        ErrorCode = "ROLE_LIST_FAILED"
	ErrorCodePermissionListFailed  ErrorCode = "PERMISSION_LIST_FAILED"
	ErrorCodeInternal              ErrorCode = "INTERNAL_ERROR"

	// Role assignment error codes
	ErrorCodeRoleAssignFailed ErrorCode = "ROLE_ASSIGN_FAILED"

	// Permission module action error codes
	ErrorCodePermissionAlreadyAssociated ErrorCode = "PERMISSION_ALREADY_ASSOCIATED"
	ErrorCodeModuleActionNotFound        ErrorCode = "MODULE_ACTION_NOT_FOUND"
	ErrorCodePermissionCreateFailed      ErrorCode = "PERMISSION_CREATE_FAILED"

	// Permission update error codes
	ErrorCodePermissionUpdateFailed ErrorCode = "PERMISSION_UPDATE_FAILED"

	// Permission error codes
	ErrorCodePermissionAlreadyExists ErrorCode = "PERMISSION_ALREADY_EXISTS"
)

// Error represents a structured error response
type Error struct {
	Type    ErrorType `json:"type"`
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// GetHTTPStatus returns the appropriate HTTP status code for the error type
func (e *Error) GetHTTPStatus() int {
	switch e.Type {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeAuthentication:
		return http.StatusUnauthorized
	case ErrorTypeAuthorization:
		return http.StatusForbidden
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// NewError creates a new error with the given type, code, and message
func NewError(errorType ErrorType, code ErrorCode, message string) *Error {
	return &Error{
		Type:    errorType,
		Code:    code,
		Message: message,
	}
}
