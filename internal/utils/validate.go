package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	validateErrs, _ := err.(validator.ValidationErrors)
	for _, e := range validateErrs {
		errField := strings.ToLower(e.Field())

		if errField == "firstname" {
			errField = "first_name"
		}

		if errField == "lastname" {
			errField = "last_name"
		}

		switch e.Tag() {
		case "required":
			errors[errField] = "This field is required."
		case "email":
			errors[errField] = "Invalid email format."
		case "min":
			errors[errField] = fmt.Sprintf("Must be at least %s characters long", e.Param())
		case "max":
			errors[errField] = fmt.Sprintf("Must not exceed %s characters long", e.Param())
		default:
			errors[errField] = fmt.Sprintf("Failed on constraint: %s", e.Tag())
		}
	}

	return errors
}
