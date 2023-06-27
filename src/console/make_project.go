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
	"path/filepath"
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
	cmd.Flags().String("target-dir", "./", "Created project path, the default is the current directory")
	cmd.Flags().String("name", "", "Project name, such as github.com/we7coreteam/w7-rangine-go-skeleton or w7-rangine-go-skeleton")
	cmd.MarkFlagRequired("name")
}

func (self MakeProjectCommand) Handle(cmd *cobra.Command, args []string) {
	moduleName, _ := cmd.Flags().GetString("name")
	projectPath, _ := cmd.Flags().GetString("target-dir")
	projectPath, _ = filepath.Abs(projectPath)

	projectPath = strings.TrimRight(projectPath, "/")
	exist, err := os.Stat(fmt.Sprintf("%s", projectPath))

	if err == nil && exist.IsDir() {
		dir, _ := os.ReadDir(projectPath)
		if len(dir) > 0 {
			cmd.PrintErrln("Error: target dir not empty")
			cmd.Usage()
			return
		}
	}
	os.MkdirAll(fmt.Sprintf("%s", projectPath), 0755)

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
			os.MkdirAll(fmt.Sprintf("%s/%s", projectPath, path), 0755)
			continue
		}

		file, _ := zipFile.Open()
		defer file.Close()
		content, _ := io.ReadAll(file)
		if path != "" {
			err = os.WriteFile(fmt.Sprintf("%s/%s", projectPath, path), []byte(self.replaceFileContent(path, content, moduleName)), 0755)
		}
	}

	color.Println("Please run command.")
	color.Println("********************************************************************")
	color.Red.Printf(" cd %s \n", projectPath)
	color.Red.Printf(" go get -u && go mod tidy \n")
	color.Red.Printf(" go build -o ./bin/rangine . \n")
	color.Println("********************************************************************")

}

func (self MakeProjectCommand) replaceFileContent(file string, content []byte, name string) string {
	if file == ".run/rangine.run.xml" {
		return strings.ReplaceAll(string(content), "w7-rangine-go-skeleton", filepath.Base(name))
	}

	if file == "app/home/provider.go" || file == "go.mod" || file == "main.go" {
		return strings.ReplaceAll(string(content), "github.com/we7coreteam/w7-rangine-go-skeleton", name)
	}

	return string(content)
}
