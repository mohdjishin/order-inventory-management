package util

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(data interface{}) (map[string]string, error) {
	err := validate.Struct(data)
	if err != nil {
		validationErrors := make(map[string]string)
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range err.(validator.ValidationErrors) {
				validationErrors[fieldErr.Field()] = fieldErr.Tag()
			}
		}
		return validationErrors, err
	}
	return nil, nil
}
