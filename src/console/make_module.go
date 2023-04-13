package console

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"html/template"
	"os"
	"strings"
)

type commandArgs struct {
	name string
}

type TemplateData struct {
	Name string
}

var argsValue commandArgs

type MakeModuleCommand struct {
	Abstract
}

func (makeModuleCommand MakeModuleCommand) GetName() string {
	return "make:module"
}

func (makeModuleCommand MakeModuleCommand) GetDescription() string {
	return "Make business module skeleton"
}

func (makeModuleCommand MakeModuleCommand) Configure(cmd *cobra.Command) {
	cmd.Flags().StringVar(&argsValue.name, "name", "", "Created module name")
	cmd.MarkFlagRequired("name")
}

func (makeModuleCommand MakeModuleCommand) Handle(cmd *cobra.Command, args []string) {
	baseDir, _ := os.Getwd()

	exist, err := os.Stat(fmt.Sprintf("%s/app/%s", baseDir, argsValue.name))
	if err == nil && exist.IsDir() {
		panic(errors.New("module is exists"))
	}
	//创建目录
	for _, dir := range makeModuleCommand.templateDir() {
		os.MkdirAll(fmt.Sprintf("%s/app/%s/%s", baseDir, argsValue.name, dir), 0755)
	}

	//创建文件
	for fileName, code := range makeModuleCommand.templateCode() {
		path := fmt.Sprintf("%s/app/%s/%s", baseDir, argsValue.name, fileName)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			panic(errors.New("error make template file " + err.Error()))
		}

		templateParser := template.New(fileName)
		templateParser, err = templateParser.Parse(strings.Trim(code, "\n"))
		if err != nil {
			panic(errors.New("error parsing template"))
		}
		templateParser.Execute(file, TemplateData{
			Name: argsValue.name,
		})
	}

	color.Println("Please copy the register provider code to the 'main.go' file.")
	color.Println("********************************************************************")
	color.Red.Printf(" app.GetProviderManager().RegisterProvider(new(%s.Provider)).Register() \n", argsValue.name)
	color.Println("********************************************************************")
}

func (makeModuleCommand MakeModuleCommand) templateCode() map[string]string {
	return map[string]string{
		"provider.go": `
package {{.Name}}

import (
	"github.com/gin-gonic/gin"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
	http_server "github.com/we7coreteam/w7-rangine-go/src/http/server"
	"github.com/we7coreteam/w7api/app/{{.Name}}/command"
	"github.com/we7coreteam/w7api/app/{{.Name}}/http/controller"
	"github.com/we7coreteam/w7api/app/{{.Name}}/http/middleware"
			
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	provider.GetConsole().RegisterCommand(new(command.Test))

	http_server.RegisterRouters(func(engine *gin.Engine) {
		engine.GET("/{{.Name}}/index", middleware.Home{}.Process, controller.Home{}.Index)
	})
}`,
		"command/test.go": `
package command

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/src/console"
)

type Test struct {
	console.Abstract
}

func (test Test) GetName() string {
	return "{{.Name}}:test"
}

func (test Test) GetDescription() string {
	return "test command"
}

func (test Test) Handle(cmd *cobra.Command, args []string) {
	color.Infoln("{{.Name}} test")
}`,
		"http/controller/home.go": `
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/we7coreteam/w7-rangine-go/src/http/controller"
)

type Home struct {
	controller.Abstract
}

func (home Home) Index(ctx *gin.Context) {
	home.JsonResponseWithoutError(ctx, "hello world!")
}`,
		"http/middleware/home.go": `
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/we7coreteam/w7-rangine-go/src/facade"
	"github.com/we7coreteam/w7-rangine-go/src/http/middleware"
	"time"
)

type Home struct {
	middleware.Abstract
}

func (home Home) Process(ctx *gin.Context) {
	log, _ := facade.GetLoggerFactory().Channel("default")
	log.Info("route middleware test, req time: " + time.Now().String())

	ctx.Next()
}`,
	}
}

func (makeModuleCommand MakeModuleCommand) templateDir() []string {
	return []string{
		"command",
		"http/controller",
		"http/middleware",
		"logic",
		"model",
	}
}
