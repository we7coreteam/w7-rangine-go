package err_handler

type ResponseError struct {
	Msg string
}

func (responseErr ResponseError) Error() string {
	return responseErr.Msg
}
