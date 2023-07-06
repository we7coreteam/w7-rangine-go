package error

import (
	"github.com/go-playground/validator/v10"
)

type ValidateFail struct {
	Msg          string
	ValidateErrs validator.ValidationErrors
}

func (validateFail ValidateFail) Error() string {
	return validateFail.Msg
}
