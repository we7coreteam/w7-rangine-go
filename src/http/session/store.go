package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/gorm"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/v2/src/components/database"
	rf "github.com/we7coreteam/w7-rangine-go/v2/src/components/redis"
	"net/http"
	"strconv"
)

func BuildOptions(config *viper.Viper) sessions.Options {
	var cookieConfig Cookie
	err := config.UnmarshalKey("cookie", &cookieConfig)
	if err != nil {
		panic(err)
	}

	return sessions.Options{
		MaxAge:   cookieConfig.MaxAge,
		SameSite: http.SameSite(cookieConfig.SameSite),
		Domain:   cookieConfig.Domain,
		Path:     cookieConfig.Path,
		Secure:   cookieConfig.Secure,
		HttpOnly: cookieConfig.HttpOnly,
	}
}

func GetMemoryStore(config *viper.Viper, keyPairs ...[]byte) sessions.Store {
	store := memstore.NewStore(keyPairs...)
	store.Options(BuildOptions(config))

	return store
}

func GetGormStore(config *viper.Viper, dbFactory *database.Factory, keyPairs ...[]byte) sessions.Store {
	sessionDb := config.GetString("session.db")
	if sessionDb == "" {
		sessionDb = "default"
	}
	db, err := dbFactory.Channel(sessionDb)
	if err != nil {
		panic(err)
	}

	store := gorm.NewStore(db, true, keyPairs...)
	store.Options(BuildOptions(config))

	return store
}

func GetRedisStore(config *viper.Viper, keyPairs ...[]byte) sessions.Store {
	sessionDb := config.GetString("session.db")
	if sessionDb == "" {
		sessionDb = "default"
	}
	var redisConfig rf.Config
	err := config.UnmarshalKey("redis."+sessionDb, &redisConfig)
	if err != nil {
		panic(err)
	}

	store, err := redis.NewStoreWithDB(int(redisConfig.PoolSize), "tcp", redisConfig.Host+":"+strconv.Itoa(int(redisConfig.Port)), redisConfig.Password, strconv.Itoa(int(redisConfig.Db)), keyPairs...)
	if err != nil {
		panic(err)
	}
	store.Options(BuildOptions(config))

	return store
}
