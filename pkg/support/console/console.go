package console

type ConsoleInterface interface {
	RegisterCommand(cmd CommandInterface)
	Run()
}
