package database

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/v2/src/console"
	"github.com/we7coreteam/w7-rangine-go/v2/src/core/logger"
)

type Provider struct {
}

func (provider Provider) Register(config *viper.Viper, loggerFactory *logger.Factory, container container.Container, consoleManager *console.Console) {
	var dbConfigMap map[string]Config
	err := config.UnmarshalKey("database", &dbConfigMap)
	if err != nil {
		panic(err)
	}

	dbFactory := NewDatabaseFactory()
	dbFactory.SetLoggerFactory(loggerFactory)
	dbFactory.Register(dbConfigMap)
	if config.GetString("app.env") == "debug" {
		dbFactory.SetDebug()
	}

	err = container.NamedSingleton("db-factory", func() *Factory {
		return dbFactory
	})
	if err != nil {
		panic(err)
	}

	consoleManager.RegisterCommand(&MakeModelCommand{
		DbFactory: dbFactory,
	})
}
