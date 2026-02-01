package ddd

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type DomainError struct {
	Message    string      `json:"message"`
	Code       string      `json:"code"`
	StatusCode int         `json:"statusCode"`
	Field      string      `json:"field,omitempty"`
	Detail     interface{} `json:"details,omitempty"`
	Err        error       `json:"-"`
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) GetCode() string {
	return e.Code
}

func (e *DomainError) GetStatusCode() int {
	return e.StatusCode
}

func (e *DomainError) GetField() string {
	return e.Field
}

func (e *DomainError) GetDetails() interface{} {
	return e.Detail
}

func (e *DomainError) GetMessage() string {
	return e.Message
}

func (e *DomainError) Wrap(err error) *DomainError {
	return &DomainError{
		Message:    e.Message,
		Code:       e.Code,
		StatusCode: e.StatusCode,
		Field:      e.Field,
		Detail:     e.Detail,
		Err:        err,
	}
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func (e *DomainError) RootCause() error {
	if e.Err == nil {
		return nil
	}
	for {
		wrapped := errors.Unwrap(e.Err)
		if wrapped == nil {
			break
		}
		e.Err = wrapped
	}
	return e.Err
}

func NewDomainError(message string, code string, statusCode int) *DomainError {
	return &DomainError{
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
		Field:      "",
		Detail:     nil,
		Err:        nil,
	}
}

func NewDomainErrorWithCause(message string, code string, err error) *DomainError {
	return &DomainError{
		Message: message,
		Code:    code,
		Err:     err,
	}
}

func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

func AsDomainError(err error) (*DomainError, bool) {
	domainErr, ok := err.(*DomainError)
	return domainErr, ok
}

func ToDomainError(err error) *DomainError {
	if err == nil {
		return nil
	}

	// If the error is a DomainError, return it
	if IsDomainError(err) {
		return err.(*DomainError)
	}

	// If the error is a database error, normalize it
	normalized := normalizeDBError(err)
	if normalized != nil {
		return normalized
	}

	// If the error is not a DomainError or a database error, return a generic internal error
	return NewDomainError(err.Error(), string(ErrorCodeInternal), 500)
}

func normalizeDBError(err error) *DomainError {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewNotFoundError("record not found")
	}

	errMsg := strings.ToLower(err.Error())

	if isUniqueConstraintViolation(errMsg) {
		field := extractFieldFromConstraintError(errMsg)
		if field != "" {
			return NewConflictError(fmt.Sprintf("%s already exists", field))
		}
		return NewConflictError("duplicate entry violates unique constraint")
	}

	if isForeignKeyViolation(errMsg) {
		return NewConflictError("operation violates foreign key constraint")
	}

	if isNotNullViolation(errMsg) {
		return NewValidationError("required field cannot be null")
	}

	if isCheckConstraintViolation(errMsg) {
		return NewValidationError("data violates check constraint")
	}

	return nil
}

func isUniqueConstraintViolation(errMsg string) bool {
	return strings.Contains(errMsg, "duplicate key") ||
		strings.Contains(errMsg, "unique constraint") ||
		strings.Contains(errMsg, "sqlstate 23505") ||
		strings.Contains(errMsg, "duplicate entry") ||
		strings.Contains(errMsg, "1062") ||
		strings.Contains(errMsg, "unique violation")
}

func isForeignKeyViolation(errMsg string) bool {
	return strings.Contains(errMsg, "foreign key constraint") ||
		strings.Contains(errMsg, "sqlstate 23503") ||
		strings.Contains(errMsg, "1452") ||
		strings.Contains(errMsg, "foreign key violation")
}

func isNotNullViolation(errMsg string) bool {
	return strings.Contains(errMsg, "not null constraint") ||
		strings.Contains(errMsg, "sqlstate 23502") ||
		strings.Contains(errMsg, "1048") ||
		strings.Contains(errMsg, "cannot be null")
}

func isCheckConstraintViolation(errMsg string) bool {
	return strings.Contains(errMsg, "check constraint") ||
		strings.Contains(errMsg, "sqlstate 23514")
}

func extractFieldFromConstraintError(errMsg string) string {
	if strings.Contains(errMsg, "uni_users_email") || strings.Contains(errMsg, "email") {
		return "email"
	}
	if strings.Contains(errMsg, "uni_users_username") || strings.Contains(errMsg, "username") {
		return "username"
	}
	if strings.Contains(errMsg, "uni_notes_slug") || strings.Contains(errMsg, "slug") {
		return "slug"
	}

	constraintPatterns := []string{
		"constraint \"uni_",
		"constraint 'uni_",
		"key '",
		"key \"",
	}

	for _, pattern := range constraintPatterns {
		idx := strings.Index(errMsg, pattern)
		if idx != -1 {
			start := idx + len(pattern)
			end := start
			for end < len(errMsg) && errMsg[end] != '"' && errMsg[end] != '\'' && errMsg[end] != ' ' {
				end++
			}
			if end > start {
				field := errMsg[start:end]
				if strings.HasPrefix(field, "uni_") {
					field = strings.TrimPrefix(field, "uni_")
					field = strings.TrimPrefix(field, "users_")
					return field
				}
			}
		}
	}

	return ""
}

type ErrorCode string

const (
	ErrorCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrorCodeConflict     ErrorCode = "CONFLICT"
	ErrorCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrorCodeInvalidInput ErrorCode = "INVALID_INPUT"
	ErrorCodeBusinessRule ErrorCode = "BAD_REQUEST"
)

func NewValidationError(message string) *DomainError {
	return NewDomainError(
		message,
		string(ErrorCodeValidation),
		400,
	)
}

func NewNotFoundError(message string) *DomainError {
	return NewDomainError(
		message,
		string(ErrorCodeNotFound),
		404,
	)
}

func NewUnauthorizedError(message string) *DomainError {
	return NewDomainError(
		message,
		string(ErrorCodeUnauthorized),
		401,
	)
}

func NewForbiddenError(message string) *DomainError {
	return NewDomainError(
		message,
		string(ErrorCodeForbidden),
		403,
	)
}

func NewConflictError(message string) *DomainError {
	return NewDomainError(
		message,
		string(ErrorCodeConflict),
		409,
	)
}

func NewInternalError(message string) *DomainError {
	return NewDomainError(
		message,
		string(ErrorCodeInternal),
		500,
	)
}

func NewInvalidInputError(message string) *DomainError {
	return NewDomainError(
		message,
		string(ErrorCodeInvalidInput),
		400,
	)
}

func NewBusinessRuleError(message string) *DomainError {
	return NewDomainError(
		message,
		string(ErrorCodeBusinessRule),
		422,
	)
}
