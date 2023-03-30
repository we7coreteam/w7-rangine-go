package error

type ResponseErrI interface {
	IsResponseErr() bool
}

type ResponseError struct {
	ResponseErrI

	Msg string
}

func (responseErr ResponseError) Error() string {
	return responseErr.Msg
}

func (responseErr ResponseError) IsResponseErr() bool {
	return true
}
