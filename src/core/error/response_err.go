package error

type ResponseError struct {
	error

	Msg string
}

func (responseErr ResponseError) Error() string {
	return responseErr.Msg
}
