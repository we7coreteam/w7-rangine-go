package error

type NotAllowErr struct {
	Err error
}

func (notAllowErr NotAllowErr) Unwrap() error {
	return notAllowErr.Err
}

func (notAllowErr NotAllowErr) Error() string {
	return notAllowErr.Err.Error()
}
