package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetSessionMiddleware(config *viper.Viper, storeResolver func(*viper.Viper, ...[]byte) sessions.Store, keyPairs ...[]byte) gin.HandlerFunc {
	config.SetDefault("session.name", "SESSIONID")

	return sessions.Sessions(config.GetString("session.name"), storeResolver(config, keyPairs...))
}
