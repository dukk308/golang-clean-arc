package ddd

import "fmt"

type DomainError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Err     error  `json:"-"`
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func NewDomainError(message string, code string) *DomainError {
	return &DomainError{
		Message: message,
		Code:    code,
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
	if IsDomainError(err) {
		return err.(*DomainError)
	}
	return NewDomainError(err.Error(), string(ErrorCodeInternal))
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
	ErrorCodeBusinessRule ErrorCode = "BUSINESS_RULE_VIOLATION"
)

func NewValidationError(message string) *DomainError {
	return NewDomainError(message, string(ErrorCodeValidation))
}

func NewNotFoundError(message string) *DomainError {
	return NewDomainError(message, string(ErrorCodeNotFound))
}

func NewUnauthorizedError(message string) *DomainError {
	return NewDomainError(message, string(ErrorCodeUnauthorized))
}

func NewForbiddenError(message string) *DomainError {
	return NewDomainError(message, string(ErrorCodeForbidden))
}

func NewConflictError(message string) *DomainError {
	return NewDomainError(message, string(ErrorCodeConflict))
}

func NewInternalError(message string) *DomainError {
	return NewDomainError(message, string(ErrorCodeInternal))
}

func NewInvalidInputError(message string) *DomainError {
	return NewDomainError(message, string(ErrorCodeInvalidInput))
}

func NewBusinessRuleError(message string) *DomainError {
	return NewDomainError(message, string(ErrorCodeBusinessRule))
}
