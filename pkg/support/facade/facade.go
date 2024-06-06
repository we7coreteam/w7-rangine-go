package facade

import (
	"github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/pkg/support/console"
	"github.com/we7coreteam/w7-rangine-go/pkg/support/database"
	"github.com/we7coreteam/w7-rangine-go/pkg/support/logger"
	"github.com/we7coreteam/w7-rangine-go/pkg/support/redis"
	"github.com/we7coreteam/w7-rangine-go/pkg/support/server"
)

var Container container.Container
var Config *viper.Viper
var Event EventBus.Bus
var LoggerFactory logger.Factory
var Validator = binding.Validator
var Console console.Console
var ServerManager server.Manager

func GetContainer() container.Container {
	return Container
}

func GetConfig() *viper.Viper {
	return Config
}

func GetEvent() EventBus.Bus {
	return Event
}

func GetLoggerFactory() logger.Factory {
	return LoggerFactory
}

func GetConsole() console.Console {
	return Console
}

func GetServerManager() server.Manager {
	return ServerManager
}

func GetRedisFactory() redis.Factory {
	var redisFactory redis.Factory
	_ = GetContainer().NamedResolve(&redisFactory, "redis-factory")

	return redisFactory
}

func GetDbFactory() database.Factory {
	var dbFactory database.Factory
	_ = GetContainer().NamedResolve(&dbFactory, "db-factory")

	return dbFactory
}

func GetTranslator() ut.Translator {
	var translator ut.Translator
	_ = GetContainer().NamedResolve(&translator, "translator")

	return translator
}

func GetValidator() binding.StructValidator {
	return Validator
}
