package facade

import (
	ut "github.com/go-playground/universal-translator"
	app "github.com/we7coreteam/w7-rangine-go/src"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/redis"
)

type Facades struct {
}

func (facade Facades) GetRedisFactory() *redis.RedisFactory {
	var redisFactory *redis.RedisFactory
	_ = app.GApp.GetContainer().NamedResolve(&redisFactory, "redis-factory")

	return redisFactory
}

func (facade Facades) GetDbFactory() *database.DatabaseFactory {
	var dbFactory *database.DatabaseFactory
	_ = app.GApp.GetContainer().NamedResolve(&dbFactory, "db-factory")

	return dbFactory
}

func (facade Facades) GetTranslator() ut.Translator {
	var translator ut.Translator
	_ = app.GApp.GetContainer().NamedResolve(&translator, "translator")

	return translator
}

var Facade = new(Facades)
