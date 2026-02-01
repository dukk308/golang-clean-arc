package validation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func jsonFormQueryTagName(field reflect.StructField) string {
	for _, tag := range []string{"json", "form", "query"} {
		name := strings.SplitN(field.Tag.Get(tag), ",", 2)[0]
		if name != "" && name != "-" {
			return name
		}
	}
	return field.Name
}

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

	v.RegisterTagNameFunc(jsonFormQueryTagName)

	if err := ValidateEmail(v); err != nil {
		return err
	}

	return nil
}
