package database

import (
	"github.com/we7coreteam/w7-rangine-go-support/src/database"
	"github.com/we7coreteam/w7-rangine-go-support/src/provider"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	var dbConfigMap map[string]Config
	err := provider.GetConfig().UnmarshalKey("database", &dbConfigMap)
	if err != nil {
		panic(err)
	}

	dbFactory := NewDatabaseFactory()
	logger, err := provider.GetLoggerFactory().Channel("default")
	if err == nil {
		dbFactory.SetLogger(logger)
	}
	dbFactory.Register(dbConfigMap)
	if provider.GetConfig().GetString("app.env") == "debug" {
		dbFactory.SetDebug()
	}

	err = provider.GetContainer().NamedSingleton("db-factory", func() database.Factory {
		return dbFactory
	})
	if err != nil {
		panic(err)
	}
}
