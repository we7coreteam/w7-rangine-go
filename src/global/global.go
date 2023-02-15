package global

import (
	"github.com/asaskevich/EventBus"
	ut "github.com/go-playground/universal-translator"
	"github.com/we7coreteam/w7-rangine-go/src/app"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/logger"
	"github.com/we7coreteam/w7-rangine-go/src/components/redis"
)

type RGlobal struct {
	App             *app.App
	LoggerFactory   *logger.LoggerFactory
	Event           EventBus.Bus
	RedisFactory    *redis.RedisFactory
	DatabaseFactory *database.DatabaseFactory
	Translator      ut.Translator
}

var G = new(RGlobal)
