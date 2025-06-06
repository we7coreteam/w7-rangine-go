package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/server"
	"github.com/we7coreteam/w7-rangine-go/v3/src/core/helper"
	"github.com/we7coreteam/w7-rangine-go/v3/src/http/response"
	"net/http"
)

type Server struct {
	server.IServer
	config Config

	Engine *gin.Engine
}

func NewHttpDefaultServer(config *viper.Viper) *Server {
	var serverConfig Config
	err := config.UnmarshalKey("server.http", &serverConfig)
	if err != nil {
		panic(err)
	}
	httpServer := NewServer(serverConfig)

	return httpServer
}

func NewServer(config Config) *Server {
	httpServer := &Server{
		config: config,
	}
	httpServer.initGinEngine()

	return httpServer
}

func (server *Server) initGinEngine() {
	gin.SetMode("release")
	server.Engine = gin.New()
	server.Engine.RedirectTrailingSlash = false
	if server.config.MaxBodySize > 0 {
		server.Engine.MaxMultipartMemory = server.config.MaxBodySize
	}

	responseObj := response.Response{}
	server.Engine.HandleMethodNotAllowed = true
	server.Engine.NoRoute(func(ctx *gin.Context) {
		responseObj.JsonResponseWithError(ctx, errors.New("Route not found, "+ctx.Request.URL.Path), http.StatusNotFound)
	})
	server.Engine.NoMethod(func(ctx *gin.Context) {
		responseObj.JsonResponseWithError(ctx, errors.New("Route not allow, "+ctx.Request.URL.Path), http.StatusMethodNotAllowed)
	})
}

func (server *Server) Use(middleware ...gin.HandlerFunc) gin.IRouter {
	server.Engine.Use(middleware...)

	return server.Engine
}

func (server *Server) RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(server.Engine)
	return server
}

func (server *Server) GetServerName() string {
	return "http"
}

func (server *Server) GetOptions() map[string]string {
	return map[string]string{
		"Host": server.config.Host,
		"Port": server.config.Port,
	}
}

func (server *Server) Start() {
	err := helper.ValidateConfig(server.config)
	if err != nil {
		panic(errors.New("http server config error, reason: " + err.Error()))
	}

	err = server.Engine.Run(server.config.Host + ":" + server.config.Port)
	if err != nil {
		panic(err)
	}
}
