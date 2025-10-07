package utils

import (
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "gte":
				errors[field] = field + " must be greater than or equal to " + e.Param()
			case "lte":
				errors[field] = field + " must be less than or equal to " + e.Param()
			case "email":
				errors[field] = field + " must be a valid email"
			default:
				errors[field] = field + " is invalid"
			}
		}
	}
	return errors
}
