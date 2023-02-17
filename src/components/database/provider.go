package database

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
)

type DatabaseProvider struct {
	provider.ProviderAbstract
}

func (databaseProvider *DatabaseProvider) Register() {
	var dbConfigMap map[string]Config
	err := databaseProvider.GetConfig().Unmarshal(&dbConfigMap)
	if err != nil {
		panic(err)
	}

	dbFactory := NewDatabaseFactory()
	logger, err := databaseProvider.GetLoggerFactory().Channel("default")
	if err == nil {
		dbFactory.SetLogger(logger)
	}
	dbFactory.Register(dbConfigMap)
	if databaseProvider.GetConfig().GetString("app.env") == "debug" {
		dbFactory.SetDebug()
	}

	err = databaseProvider.GetContainer().NamedSingleton("db-factory", func() *DatabaseFactory {
		return dbFactory
	})
	if err != nil {
		panic(err)
	}
}
