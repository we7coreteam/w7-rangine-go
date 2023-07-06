package helper

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ValidateAndGetErrFields(obj any) []string {
	err := binding.Validator.ValidateStruct(obj)

	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fields := make([]string, len(validationErrors))
			for index, e := range validationErrors {
				fields[index] = e.Field()
			}
			return fields
		}

		return []string{err.Error()}
	}

	return nil
}
