package console

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
	"io"
	"net/http"
	"os"
	"strings"
)

type MakeProjectCommand struct {
	Abstract
}

func (self MakeProjectCommand) GetName() string {
	return "make:project"
}

func (self MakeProjectCommand) GetDescription() string {
	return "Create project skeleton"
}

func (self MakeProjectCommand) Configure(cmd *cobra.Command) {
	cmd.Flags().String("target-dir", "", "Created project path")
	//cmd.MarkFlagRequired("target-dir")
}

func (self MakeProjectCommand) Handle(cmd *cobra.Command, args []string) {
	projectName, _ := cmd.Flags().GetString("target-dir")
	if projectName == "" && len(args) > 0 {
		projectName = args[0]
	}
	if projectName == "" {
		cmd.PrintErrln("Error: required flag(s) \"target-dir\" not set")
		cmd.Usage()
		return
	}
	projectName = strings.TrimRight(projectName, "/")
	exist, err := os.Stat(fmt.Sprintf("%s", projectName))
	println(projectName)
	if err == nil && exist.IsDir() {
		dir, _ := os.ReadDir(projectName)
		if len(dir) > 0 {
			cmd.PrintErrln("Error: target dir not empty")
			cmd.Usage()
			return
		}
	}
	os.MkdirAll(fmt.Sprintf("%s", projectName), 0755)

	response, _ := http.Get("https://codeload.github.com/we7coreteam/w7-rangine-go-skeleton/zip/refs/heads/main")
	defer response.Body.Close()
	data, _ := io.ReadAll(response.Body)
	if len(data) == 0 {
		cmd.PrintErrln("Error: Failed to download zip file, Please download from https://codeload.github.com/we7coreteam/w7-rangine-go-skeleton/zip/refs/heads/main")
		cmd.Usage()
		return
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err_handler.Found(err) {
		cmd.PrintErrln("Error: " + err.Error())
		cmd.Usage()
		return
	}

	for _, zipFile := range reader.File {
		path := strings.Replace(zipFile.Name, "w7-rangine-go-skeleton-main/", "", -1)
		if path != "" && zipFile.FileInfo().IsDir() {
			os.MkdirAll(fmt.Sprintf("%s/%s", projectName, path), 0755)
			continue
		}

		file, _ := zipFile.Open()
		defer file.Close()
		content, _ := io.ReadAll(file)
		if path != "" {
			err = os.WriteFile(fmt.Sprintf("%s/%s", projectName, path), content, 0755)
		}
	}

	color.Println("Please run command.")
	color.Println("********************************************************************")
	color.Red.Printf(" cd %s \n", projectName)
	color.Red.Printf(" go get -u && go mod tidy \n")
	color.Red.Printf(" go build -o ./bin/rangine . \n")
	color.Println("********************************************************************")

}
