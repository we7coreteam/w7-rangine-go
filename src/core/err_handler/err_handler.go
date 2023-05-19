package err_handler

import (
	"github.com/pkg/errors"
)

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
