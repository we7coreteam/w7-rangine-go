package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go-support/src/server"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
	httperf "github.com/we7coreteam/w7-rangine-go/src/http/error"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
	"net/http"
)

type Server struct {
	server.Server
	config Config

	Engine  *gin.Engine
	Session *session.Session
}

func NewHttpDefaultServer(config *viper.Viper) *Server {
	var sessionConfig session.SessionConf
	var cookieConfig session.Cookie
	var serverConfig Config
	err := config.UnmarshalKey("server.http", &serverConfig)
	if err != nil {
		panic(err)
	}
	err = config.UnmarshalKey("session", &sessionConfig)
	if err != nil {
		panic(err)
	}
	err = config.UnmarshalKey("cookie", &cookieConfig)
	if err != nil {
		panic(err)
	}

	httpServer := NewServer(serverConfig)
	httpServer.Session = session.NewSession(sessionConfig, cookieConfig)

	return httpServer
}

func NewServer(config Config) *Server {
	httpServer := &Server{
		config: config,
	}
	httpServer.initGinEngine()

	responseObj := response.Response{}
	httpServer.Engine.HandleMethodNotAllowed = true
	httpServer.Engine.NoRoute(func(context *gin.Context) {
		responseObj.JsonResponseWithError(context, httperf.NotFoundErr{
			Err: errorhandler.ResponseError{
				Msg: "Route not found, " + context.Request.URL.Path,
			},
		}, http.StatusNotFound)
	})
	httpServer.Engine.NoMethod(func(context *gin.Context) {
		responseObj.JsonResponseWithError(context, httperf.NotAllowErr{
			Err: errorhandler.ResponseError{
				Msg: "Route not allow, " + context.Request.URL.Path,
			},
		}, http.StatusMethodNotAllowed)
	})

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
		"Host": server.config.Host,
		"Port": server.config.Port,
	}
}

func (server *Server) Start() {
	if server.Session != nil {
		server.Session.Init()
	}

	err := server.Engine.Run(server.config.Host + ":" + server.config.Port)
	if err != nil {
		panic(err)
	}
}
