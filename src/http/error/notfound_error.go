package error

type NotFoundErr struct {
	Err error
}

func (notFoundErr NotFoundErr) Unwrap() error {
	return notFoundErr.Err
}

func (notFoundErr NotFoundErr) Error() string {
	return notFoundErr.Err.Error()
}
