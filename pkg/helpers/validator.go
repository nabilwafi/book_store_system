package helpers

import (
	"errors"
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct using the global validator instance
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, getErrorMessage(err))
		}
		return errors.New(strings.Join(validationErrors, ", "))
	}
	return nil
}

// getErrorMessage returns a user-friendly error message based on the validation tag
func getErrorMessage(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
		}
		return fmt.Sprintf("%s must be at least %s", field, fe.Param())
	case "max":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("%s cannot exceed %s characters", field, fe.Param())
		}
		return fmt.Sprintf("%s cannot exceed %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "dive":
		return fmt.Sprintf("%s contains invalid items", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}