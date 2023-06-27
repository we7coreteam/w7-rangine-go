package error

import (
	"github.com/go-playground/validator/v10"
)

type ValidateFail struct {
	Err          error
	ValidateErrs validator.ValidationErrors
}

func (validateFail ValidateFail) Unwrap() error {
	return validateFail.Err
}

func (validateFail ValidateFail) Error() string {
	return validateFail.Err.Error()
}
