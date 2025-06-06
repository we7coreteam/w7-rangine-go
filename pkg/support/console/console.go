package console

type IConsole interface {
	RegisterCommand(cmd ICommand)
	Run()
}
