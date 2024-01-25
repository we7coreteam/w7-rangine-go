package console

type Console interface {
	RegisterCommand(cmd Command)
	Run()
}
