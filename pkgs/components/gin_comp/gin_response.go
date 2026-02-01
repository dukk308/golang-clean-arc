package gin_comp

import (
	"net/http"

	"github.com/dukk308/beetool.dev-go-starter/internal/common"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/base"
	"github.com/dukk308/beetool.dev-go-starter/pkgs/constants"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ResponseError(c *gin.Context, err error) {
	c.Set(constants.ContextKeyError, err)
	if base.IsDomainError(err) {
		domainErr := base.ToDomainError(err)
		statusCode := getStatusCodeForError(domainErr.Code)
		c.JSON(statusCode, domainErr)
		return
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make([]map[string]interface{}, 0, len(validationErrors))
		for _, fieldError := range validationErrors {
			errors = append(errors, map[string]interface{}{
				"field":   fieldError.Field(),
				"tag":     fieldError.Tag(),
				"message": getValidationErrorMessage(fieldError),
			})
		}
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"errors": errors,
			"code":   string(base.ErrorCodeValidation),
		})
		return
	}

	normalizedErr := base.ToDomainError(err)
	statusCode := getStatusCodeForError(normalizedErr.Code)
	c.JSON(statusCode, normalizedErr)
}

func getStatusCodeForError(code string) int {
	switch code {
	case string(base.ErrorCodeValidation), string(base.ErrorCodeInvalidInput):
		return http.StatusBadRequest
	case string(base.ErrorCodeUnauthorized):
		return http.StatusUnauthorized
	case string(base.ErrorCodeForbidden):
		return http.StatusForbidden
	case string(base.ErrorCodeNotFound):
		return http.StatusNotFound
	case string(base.ErrorCodeConflict):
		return http.StatusConflict
	case string(base.ErrorCodeBusinessRule):
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func getValidationErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fieldError.Field() + " is required"
	case "min":
		return fieldError.Field() + " must be at least " + fieldError.Param() + " characters"
	case "max":
		return fieldError.Field() + " must be at most " + fieldError.Param() + " characters"
	case "email":
		return fieldError.Field() + " must be a valid email address"
	default:
		return fieldError.Field() + " is invalid"
	}
}

func ResponseSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, common.NewResponseSuccess(data))
}

func ResponseSuccessCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, common.NewResponseSuccess(data))
}
