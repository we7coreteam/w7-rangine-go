package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
	"strconv"
)

var GHttpServer *Server

type Server struct {
	config *viper.Viper

	GinEngine *gin.Engine
	Session   *session.Session
}

func NewHttpSerer(config *viper.Viper) *Server {
	server := &Server{
		config: config,
	}
	server.initGinEngine()
	GHttpServer = server

	return server
}

func (server *Server) initGinEngine() {
	gin.SetMode(server.config.GetString("app.env"))
	server.GinEngine = gin.Default()
}

func (server *Server) RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(server.GinEngine)
	return server
}

func (server *Server) Start() {
	var serverConfig Config
	err := server.config.UnmarshalKey("http_server", &serverConfig)
	if err != nil {
		panic(err)
	}
	err = server.GinEngine.Run(serverConfig.Host + ":" + strconv.Itoa(serverConfig.Port))

	panic(err)
}
