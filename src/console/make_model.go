package console

import (
	"github.com/spf13/cobra"
	yamlgen "github.com/we7coreteam/gorm-gen-yaml"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/facade"
	"github.com/we7coreteam/w7-rangine-go/v3/src/core/err_handler"
	"gorm.io/gen"
)

type MakeModelCommand struct {
	Abstract
}

func (self MakeModelCommand) GetName() string {
	return "make:model"
}

func (self MakeModelCommand) GetDescription() string {
	return "Gen gorm model"
}

func (self MakeModelCommand) Configure(command *cobra.Command) {
	command.Flags().String("table-name", "", "Table name")
	command.Flags().String("db-channel", "default", "")
	command.Flags().String("yaml-file", "", "Specify the yaml configuration file")
}

func (self MakeModelCommand) Handle(cmd *cobra.Command, args []string) {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./common/dao",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		ModelPkgPath: "./common/entity",
	})
	dbChannel, _ := cmd.Flags().GetString("db-channel")
	db, err := facade.GetDbFactory().Channel(dbChannel)
	if err_handler.Found(err) {
		cmd.PrintErrln("Db channel not found")
		cmd.Usage()
		return
	}
	g.UseDB(db)

	//// 模型自定义选项组
	fieldOpts := []gen.ModelOpt{}

	tableName, _ := cmd.Flags().GetString("table-name")
	yamlFile, _ := cmd.Flags().GetString("yaml-file")

	if tableName == "" && yamlFile == "" {
		g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)
	}

	if yamlFile != "" {
		yamlgen.NewYamlGenerator(yamlFile).UseGormGenerator(g).Generate()
	} else {
		g.ApplyBasic(
			g.GenerateModel(tableName),
		)
	}
	//g.GenerateModel(tableName)
	g.Execute()
}
