package error

import (
	"github.com/pkg/errors"
)

type ErrHandler struct {
	err error
}

type ErrCustomerHandler func(err error)
type FinallyHandler func(err error)

func (errHandler ErrHandler) Is(target ...error) bool {
	for _, arg := range target {
		if errors.Is(errHandler.err, arg) {
			return true
		}
	}
	return false
}

func Throw(message string, previous error) error {
	var err error
	if previous == nil {
		err = errors.New(message)
	} else {
		err = errors.Wrap(previous, message)
	}

	return err
}

func (errHandler ErrHandler) Found() bool {
	if errHandler.err != nil {
		return true
	}
	return false
}

func Try(err error) *ErrHandler {
	return &ErrHandler{err: err}
}

func (errHandler ErrHandler) Catch(err error, handler ErrCustomerHandler) ErrHandler {
	if errors.Is(errHandler.err, err) {
		handler(errHandler.err)
	}

	return errHandler
}

func (errHandler ErrHandler) Finally(handler FinallyHandler) {
	handler(errHandler.err)
}
