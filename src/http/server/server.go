package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	app "github.com/we7coreteam/w7-rangine-go/src"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
	"strconv"
)

var GHttpServer *Server

type Server struct {
	config *viper.Viper

	GinEngine *gin.Engine
	Session   *session.Session
}

func NewHttpDefaultServer(app *app.App) *Server {
	var sessionConfig session.SessionConf
	var cookieConfig session.Cookie
	err := app.GetConfig().UnmarshalKey("session", &sessionConfig)
	if err != nil {
		panic(err)
	}
	err = app.GetConfig().UnmarshalKey("cookie", &cookieConfig)
	if err != nil {
		panic(err)
	}

	server := NewServer(app)
	server.Session = session.NewSession(sessionConfig, cookieConfig)

	return server
}

func NewServer(app *app.App) *Server {
	server := &Server{
		config: app.GetConfig(),
	}
	server.initGinEngine()
	GHttpServer = server

	return server
}

func (server *Server) initGinEngine() {
	gin.SetMode(server.config.GetString("app.env"))
	server.GinEngine = gin.New()
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
