package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go-support/src/server"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
)

type Server struct {
	server.Server
	config *viper.Viper

	Engine  *gin.Engine
	Session *session.Session
}

func NewHttpDefaultServer(config *viper.Viper) *Server {
	var sessionConfig session.SessionConf
	var cookieConfig session.Cookie
	err := config.UnmarshalKey("session", &sessionConfig)
	if err != nil {
		panic(err)
	}
	err = config.UnmarshalKey("cookie", &cookieConfig)
	if err != nil {
		panic(err)
	}

	httpServer := NewServer(config)
	httpServer.Session = session.NewSession(sessionConfig, cookieConfig)

	return httpServer
}

func NewServer(config *viper.Viper) *Server {
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
}

func (server *Server) Use(middleware ...gin.HandlerFunc) gin.IRouter {
	server.Engine.Use(middleware...)

	return server.Engine
}

func (server *Server) RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(server.Engine)
	return server
}

func (server *Server) GetSession() *session.Session {
	return server.Session
}

func (server *Server) GetServerName() string {
	return "http"
}

func (server *Server) GetOptions() map[string]string {
	return map[string]string{
		"Host": server.config.GetString("server.http.host"),
		"Port": server.config.GetString("server.http.port"),
	}
}

func (server *Server) Start() {
	if server.Session != nil {
		server.Session.Init()
	}

	err := server.Engine.Run(server.config.GetString("server.http.host") + ":" + server.config.GetString("server.http.port"))
	if err != nil {
		panic(err)
	}
}
