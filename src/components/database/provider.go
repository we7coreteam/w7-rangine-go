package database

import (
	"github.com/we7coreteam/w7-rangine-go-support/src/database"
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	"github.com/we7coreteam/w7-rangine-go-support/src/provider"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	var dbConfigMap map[string]Config
	err := facade.GetConfig().UnmarshalKey("database", &dbConfigMap)
	if err != nil {
		panic(err)
	}

	dbFactory := NewDatabaseFactory()
	logger, err := facade.GetLoggerFactory().Channel("default")
	if err == nil {
		dbFactory.SetLogger(logger)
	}
	dbFactory.Register(dbConfigMap)
	if facade.GetConfig().GetString("app.env") == "debug" {
		dbFactory.SetDebug()
	}

	err = facade.GetContainer().NamedSingleton("db-factory", func() database.Factory {
		return dbFactory
	})
	if err != nil {
		panic(err)
	}
}
