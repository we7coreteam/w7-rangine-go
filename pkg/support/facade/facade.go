package facade

import (
	"github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/console"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/database"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/logger"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/redis"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/server"
)

var Container container.Container
var Config *viper.Viper
var Event EventBus.Bus
var LoggerFactory logger.FactoryInterface
var Validator = binding.Validator
var Console console.ConsoleInterface
var ServerManager server.ManagerInterface

func GetContainer() container.Container {
	return Container
}

func GetConfig() *viper.Viper {
	return Config
}

func GetEvent() EventBus.Bus {
	return Event
}

func GetLoggerFactory() logger.FactoryInterface {
	return LoggerFactory
}

func GetConsole() console.ConsoleInterface {
	return Console
}

func GetServerManager() server.ManagerInterface {
	return ServerManager
}

func GetRedisFactory() redis.FactoryInterface {
	var redisFactory redis.FactoryInterface
	_ = GetContainer().NamedResolve(&redisFactory, "redis-factory")

	return redisFactory
}

func GetDbFactory() database.FactoryInterface {
	var dbFactory database.FactoryInterface
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
