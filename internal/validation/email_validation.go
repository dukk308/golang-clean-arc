package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func ValidateEmail(v *validator.Validate) error {
	if err := v.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		return emailRegex.MatchString(fl.Field().String())
	}); err != nil {
		return err
	}

	return nil
}
