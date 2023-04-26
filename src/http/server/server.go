package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go-support/src/server"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
)

var DefaultHttpServer *Server

func GetServer() *Server {
	return DefaultHttpServer
}

func Use(middleware ...gin.HandlerFunc) gin.IRouter {
	DefaultHttpServer.Engine.Use(middleware...)

	return DefaultHttpServer.Engine
}

func RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(DefaultHttpServer.Engine)
	return DefaultHttpServer
}

func GetSession() *session.Session {
	return DefaultHttpServer.Session
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

	DefaultHttpServer = NewServer(config)
	DefaultHttpServer.Session = session.NewSession(sessionConfig, cookieConfig)

	return DefaultHttpServer
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
