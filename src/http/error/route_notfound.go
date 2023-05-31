package error

type RouteNotFound struct {
	Err error
}

func (routeNotFound RouteNotFound) Unwrap() error {
	return routeNotFound.Err
}

func (routeNotFound RouteNotFound) Error() string {
	return routeNotFound.Err.Error()
}
