package console

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"html/template"
	"os"
	"strings"
)

type MakeCmdTemplateData struct {
	Name    string
	Module  string
	Command string
}

type MakeCmdCommand struct {
	Abstract
}

func (self MakeCmdCommand) GetName() string {
	return "make:cmd"
}

func (self MakeCmdCommand) GetDescription() string {
	return "Gen cmd file"
}

func (self MakeCmdCommand) Configure(command *cobra.Command) {
	command.Flags().String("name", "", "Set command name")
	command.MarkFlagRequired("name")
	command.Flags().String("module-name", "", "Set module name")
	command.MarkFlagRequired("module-name")
}

func (self MakeCmdCommand) Handle(cmd *cobra.Command, args []string) {
	baseDir, _ := os.Getwd()
	moduleName, _ := cmd.Flags().GetString("module-name")
	exist, err := os.Stat(fmt.Sprintf("%s/app/%s", baseDir, moduleName))
	if err == nil && !exist.IsDir() {
		cmd.PrintErrln("Error: Module is not exists")
		cmd.Usage()
		return
	}

	fileName, _ := cmd.Flags().GetString("name")
	path := fmt.Sprintf("%s/app/%s/command/%s.go", baseDir, moduleName, fileName)
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)

	templateParser := template.New(fileName)
	templateParser, _ = templateParser.Parse(strings.Trim(self.templateCode(), "\n"))
	templateParser.Execute(file, MakeCmdTemplateData{
		Name:    strings.ToUpper(fileName[:1]) + fileName[1:],
		Command: fmt.Sprintf("%s:%s", moduleName, fileName),
		Module:  moduleName,
	})

	color.Printf("Generate command file: %s \n", path)
	color.Printf("Please copy the register provider code to the '%s/provider.go' file. \n", moduleName)
	color.Println("********************************************************************")
	color.Red.Printf(" console.RegisterCommand(new(command.%s)) \n", strings.ToUpper(fileName[:1])+fileName[1:])
	color.Println("********************************************************************")
}

func (self MakeCmdCommand) templateCode() string {
	return `
package command

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/v2/src/console"
)

type {{.Name}} struct {
	console.Abstract
}

func (self {{.Name}}) GetName() string {
	return "{{.Command}}"
}

func (self {{.Name}}) GetDescription() string {
	return "{{.Name}} command"
}

func (self {{.Name}}) Configure(command *cobra.Command)  {
	command.Flags().String("name", "test", "test name params")
}

func (self {{.Name}}) Handle(cmd *cobra.Command, args []string) {
	color.Infoln("Run {{.Name}} command")
}
`
}
