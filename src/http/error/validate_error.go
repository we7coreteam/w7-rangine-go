package error

import (
	"github.com/go-playground/validator/v10"
)

type ValidateErr struct {
	Err          error
	ValidateErrs validator.ValidationErrors
}

func (validateErr ValidateErr) Unwrap() error {
	return validateErr.Err
}

func (validateErr ValidateErr) Error() string {
	return validateErr.Err.Error()
}
