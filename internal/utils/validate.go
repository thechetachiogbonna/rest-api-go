package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = newValidator()

func newValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return field.Name
		}
		return name
	})
	return validate
}

func FormatValidationError(err error) map[string]string {
	formatted := make(map[string]string)

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string]string{"error": err.Error()}
	}

	for _, fieldError := range validationErrors {
		switch fieldError.Tag() {
		case "required":
			formatted[fieldError.Field()] = "This field is required."
		case "email":
			formatted[fieldError.Field()] = "Invalid email format."
		case "min":
			formatted[fieldError.Field()] = fmt.Sprintf("Must be at least %s.", validationLimit(fieldError))
		case "max":
			formatted[fieldError.Field()] = fmt.Sprintf("Must be no more than %s.", validationLimit(fieldError))
		default:
			formatted[fieldError.Field()] = fmt.Sprintf("Failed on constraint: %s", fieldError.Tag())
		}
	}

	return formatted
}

func validationLimit(fieldError validator.FieldError) string {
	switch fieldError.Kind() {
	case reflect.String:
		return fmt.Sprintf("%s characters", fieldError.Param())
	case reflect.Array, reflect.Slice, reflect.Map:
		return fmt.Sprintf("%s items", fieldError.Param())
	default:
		return fieldError.Param()
	}
}
