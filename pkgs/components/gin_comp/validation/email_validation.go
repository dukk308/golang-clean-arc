package validation

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func RegisterValidations() error {
	if binding.Validator == nil {
		return errors.New("binding validator is not initialized")
	}

	engine := binding.Validator.Engine()
	if engine == nil {
		return errors.New("validator engine is nil")
	}

	v, ok := engine.(*validator.Validate)
	if !ok {
		return errors.New("validator engine is not *validator.Validate")
	}

	if err := ValidateEmail(v); err != nil {
		return err
	}

	return nil
}

func ValidateEmail(v *validator.Validate) error {
	if err := v.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		return emailRegex.MatchString(fl.Field().String())
	}); err != nil {
		return err
	}

	return nil
}