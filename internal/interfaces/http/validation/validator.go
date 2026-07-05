package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct and returns a user-friendly error message
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var errorMessages []string
	for _, fieldError := range validationErrors {
		errorMessages = append(errorMessages, formatFieldError(fieldError))
	}

	return fmt.Errorf("%s", strings.Join(errorMessages, "; "))
}

func formatFieldError(fieldError validator.FieldError) string {
	field := strings.ToLower(fieldError.Field())

	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fieldError.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fieldError.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
