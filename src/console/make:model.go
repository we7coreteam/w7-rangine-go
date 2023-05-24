package console

import (
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	"github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
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
	command.MarkFlagRequired("table-name")
}

func (self MakeModelCommand) Handle(cmd *cobra.Command, args []string) {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./common/entity",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
		ModelPkgPath: "entity",
	})
	dbChannel, _ := cmd.Flags().GetString("db-channel")
	db, err := facade.GetDbFactory().Channel(dbChannel)
	if err_handler.Found(err) {
		cmd.PrintErrln("Db channel not found")
		cmd.Usage()
		return
	}
	g.UseDB(db)

	tableName, _ := cmd.Flags().GetString("table-name")
	g.GenerateModel(tableName)
	g.Execute()
}
