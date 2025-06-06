package server

type IServer interface {
	GetServerName() string
	Start()
	Stop()
	GetOptions() map[string]string
}
