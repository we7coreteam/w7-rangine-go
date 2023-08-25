package console

import (
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	"github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
	"gorm.io/gen"
	"gorm.io/gen/field"
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

	autoCreateTimeField := gen.FieldGORMTag("create_at", func(tag field.GormTag) field.GormTag {
		tag.Set("autoCreateTime")
		return tag
	})
	autoUpdateTimeField := gen.FieldGORMTag("update_at", func(tag field.GormTag) field.GormTag {
		tag.Set("autoUpdateTime")
		return tag
	})
	//// 模型自定义选项组
	fieldOpts := []gen.ModelOpt{autoCreateTimeField, autoUpdateTimeField}

	tableName, _ := cmd.Flags().GetString("table-name")
	if tableName == "" {
		g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)
	} else {
		g.ApplyBasic(
			g.GenerateModel(tableName, fieldOpts...),
		)
	}
	//g.GenerateModel(tableName)
	g.Execute()
}
