package server

import (
	"github.com/gin-gonic/gin"
	app "github.com/we7coreteam/w7-rangine-go/src"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
	"strconv"
)

type Server struct {
	App *app.App

	GinEngine *gin.Engine
	Session   *session.Session
}

func NewHttpSerer(app *app.App) *Server {
	server := &Server{
		App: app,
	}
	server.initGinEngine()

	return server
}

func (server *Server) initGinEngine() {
	server.App.GetConfig().SetDefault("app.env", "release")
	gin.SetMode(server.App.GetConfig().GetString("app.env"))
	server.GinEngine = gin.Default()
}

func (server *Server) RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(server.GinEngine)
	return server
}

func (server *Server) Start() {
	var serverConfig Config
	err := server.App.GetConfig().Unmarshal(&serverConfig)
	if err != nil {
		panic(err)
	}

	err = server.GinEngine.Run(serverConfig.Host + ":" + strconv.Itoa(serverConfig.Port))
	if err != nil {
		panic(err)
	}
}
