package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go-support/src/server"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
)

var GHttpServer *Server

func GetServer() *Server {
	return GHttpServer
}

func Use(middleware ...gin.HandlerFunc) gin.IRouter {
	GHttpServer.Engine.Use(middleware...)

	return GHttpServer.Engine
}

func RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(GHttpServer.Engine)
	return GHttpServer
}

func GetSession() *session.Session {
	return GHttpServer.Session
}

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

	newServer := NewServer(config)
	newServer.Session = session.NewSession(sessionConfig, cookieConfig)

	return newServer
}

func NewServer(config *viper.Viper) *Server {
	GHttpServer = &Server{
		config: config,
	}
	GHttpServer.initGinEngine()

	return GHttpServer
}

func (server *Server) initGinEngine() {
	gin.SetMode("release")
	server.Engine = gin.New()
	server.Engine.RedirectTrailingSlash = false
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
