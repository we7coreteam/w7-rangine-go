package database

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go-support/src/database"
	"github.com/we7coreteam/w7-rangine-go-support/src/logger"
)

type Provider struct {
}

func (provider Provider) Register(config *viper.Viper, loggerFactory logger.Factory, container container.Container) {
	var dbConfigMap map[string]Config
	err := config.UnmarshalKey("database", &dbConfigMap)
	if err != nil {
		panic(err)
	}

	dbFactory := NewDatabaseFactory()
	logger, err := loggerFactory.Channel("default")
	if err == nil {
		dbFactory.SetLogger(logger)
	}
	dbFactory.Register(dbConfigMap)
	if config.GetString("app.env") == "debug" {
		dbFactory.SetDebug()
	}

	err = container.NamedSingleton("db-factory", func() database.Factory {
		return dbFactory
	})
	if err != nil {
		panic(err)
	}
}
