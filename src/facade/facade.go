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

func GetContainer() container.Container {
	return app.GApp.GetContainer()
}

func GetConfig() *viper.Viper {
	return app.GApp.GetConfig()
}

func GetEvent() EventBus.Bus {
	return app.GApp.GetEvent()
}

func GetLoggerFactory() *logger.Factory {
	return app.GApp.GetLoggerFactory()
}

func GetRedisFactory() *redis.Factory {
	var redisFactory *redis.Factory
	_ = app.GApp.GetContainer().NamedResolve(&redisFactory, "redis-factory")

	return redisFactory
}

func GetDbFactory() *database.Factory {
	var dbFactory *database.Factory
	_ = app.GApp.GetContainer().NamedResolve(&dbFactory, "db-factory")

	return dbFactory
}

func GetTranslator() ut.Translator {
	var translator ut.Translator
	_ = app.GApp.GetContainer().NamedResolve(&translator, "translator")

	return translator
}

func GetHttpServer() *httpserver.Server {
	var server *httpserver.Server
	_ = app.GApp.GetContainer().NamedResolve(&server, "http-server")

	return server
}
