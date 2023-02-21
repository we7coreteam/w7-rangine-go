package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
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
	err := server.config.Unmarshal(&serverConfig)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    serverConfig.Host + ":" + strconv.Itoa(serverConfig.Port),
		Handler: server.GinEngine.Handler(),
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	//设置 5 秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
