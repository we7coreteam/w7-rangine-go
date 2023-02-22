package facade

import (
	"github.com/asaskevich/EventBus"
	ut "github.com/go-playground/universal-translator"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	app "github.com/we7coreteam/w7-rangine-go/src"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/redis"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
	httpserver "github.com/we7coreteam/w7-rangine-go/src/http/server"
)

type Facades struct {
}

func (facade Facades) GetContainer() container.Container {
	return app.GApp.GetContainer()
}

func (facade Facades) GetConfig() *viper.Viper {
	return app.GApp.GetConfig()
}

func (facade Facades) GetEvent() EventBus.Bus {
	return app.GApp.GetEvent()
}

func (facade Facades) GetLoggerFactory() *logger.Factory {
	return app.GApp.GetLoggerFactory()
}

func (facade Facades) GetRedisFactory() *redis.Factory {
	var redisFactory *redis.Factory
	_ = app.GApp.GetContainer().NamedResolve(&redisFactory, "redis-factory")

	return redisFactory
}

func (facade Facades) GetDbFactory() *database.Factory {
	var dbFactory *database.Factory
	_ = app.GApp.GetContainer().NamedResolve(&dbFactory, "db-factory")

	return dbFactory
}

func (facade Facades) GetTranslator() ut.Translator {
	var translator ut.Translator
	_ = app.GApp.GetContainer().NamedResolve(&translator, "translator")

	return translator
}

func (facade Facades) GetHttpServer() *httpserver.Server {
	var server *httpserver.Server
	_ = app.GApp.GetContainer().NamedResolve(&server, "http-server")

	return server
}

var Facade = new(Facades)
