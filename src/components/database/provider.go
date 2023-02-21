package database

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	err := provider.GetContainer().NamedSingleton("db-factory", func() *Factory {
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

		return dbFactory
	})
	if err != nil {
		panic(err)
	}
}
