package err_handler

import (
	"github.com/pkg/errors"
)

type ErrHandler struct {
	err error
}

type ErrCustomerHandler func(err error)
type FinallyHandler func(err error)

func Throw(message string, previous error) error {
	var err error
	if previous == nil {
		err = errors.New(message)
	} else {
		err = errors.Wrap(previous, message)
	}

	return err
}

func Found(err error) bool {
	if err != nil {
		return true
	}
	return false
}

func Try(err error) *ErrHandler {
	return &ErrHandler{err: err}
}

func (errHandler ErrHandler) Is(target ...error) bool {
	if !Found(errHandler.err) {
		return false
	}

	for _, arg := range target {
		if errors.Is(errHandler.err, arg) {
			return true
		}
	}
	return false
}

func (errHandler ErrHandler) Catch(err any, handler ErrCustomerHandler) ErrHandler {
	if Found(errHandler.err) && errors.As(errHandler.err, err) {
		handler(errHandler.err)
	}

	return errHandler
}

func (errHandler ErrHandler) Finally(handler FinallyHandler) {
	if Found(errHandler.err) {
		handler(errHandler.err)
	}
}
