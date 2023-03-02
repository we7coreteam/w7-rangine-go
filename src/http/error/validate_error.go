package error

import (
	"github.com/go-playground/validator/v10"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/error"
)

type ValidateErr struct {
	errorhandler.ResponseError
	ValidateErrs validator.ValidationErrors
}
