package domain

import "regexp"

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type EmailVO struct {
	Value string `json:"value"`
}

func NewEmailVO(email string) *EmailVO {
	return &EmailVO{
		Value: email,
	}
}

func (e *EmailVO) Validate() error {
	if e.Value == "" {
		return ErrInvalidEmail
	}
	if !emailRegex.MatchString(e.Value) {
		return ErrInvalidEmail
	}

	return nil
}
