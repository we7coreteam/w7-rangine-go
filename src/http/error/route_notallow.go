package error

type RouteNotAllow struct {
	Err error
}

func (routeNotAllow RouteNotAllow) Unwrap() error {
	return routeNotAllow.Err
}

func (routeNotAllow RouteNotAllow) Error() string {
	return routeNotAllow.Err.Error()
}
