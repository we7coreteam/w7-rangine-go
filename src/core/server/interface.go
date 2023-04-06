package server

type Interface interface {
	GetServerName() string
	Start()
}
