package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/server"
	"github.com/we7coreteam/w7-rangine-go/v2/src/core/helper"
	"github.com/we7coreteam/w7-rangine-go/v2/src/http/response"
)

type Server struct {
	server.Server
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
	options := map[string]string{
		"Host": server.config.Host,
		"Mode": server.mode(),
	}
	if server.plainHTTPEnabled() {
		options["Port"] = server.config.Port
	}
	if server.httpsEnabled() {
		options["TLSPort"] = server.httpsPort()
	}

	return options
}

func (server *Server) Start() {
	err := server.validateConfig()
	if err != nil {
		panic(errors.New("http server config error, reason: " + err.Error()))
	}

	if server.plainHTTPEnabled() && server.httpsEnabled() {
		go server.startHTTP()
		server.startHTTPS()
		return
	}

	if server.httpsEnabled() {
		server.startHTTPS()
		return
	}

	server.startHTTP()
}

func (server *Server) startHTTP() {
	err := server.Engine.Run(server.httpAddress())
	if err != nil {
		panic(err)
	}
}

func (server *Server) startHTTPS() {
	err := server.Engine.RunTLS(server.httpsAddress(), server.config.TLS.CertFile, server.config.TLS.KeyFile)
	if err != nil {
		panic(err)
	}
}

func (server *Server) httpAddress() string {
	return server.config.Host + ":" + server.config.Port
}

func (server *Server) httpsAddress() string {
	return server.config.Host + ":" + server.httpsPort()
}

func (server *Server) httpsPort() string {
	if server.config.TLS.Port != "" {
		return server.config.TLS.Port
	}

	return server.config.Port
}

func (server *Server) plainHTTPEnabled() bool {
	if server.config.Port == "" {
		return false
	}
	if !server.config.TLS.Enable {
		return true
	}

	// When TLS is enabled without tls.port, port is kept as HTTPS-only for compatibility.
	return server.config.TLS.Port != ""
}

func (server *Server) httpsEnabled() bool {
	return server.config.TLS.Enable
}

func (server *Server) mode() string {
	if server.plainHTTPEnabled() && server.httpsEnabled() {
		return "http+https"
	}
	if server.httpsEnabled() {
		return "https"
	}
	if server.plainHTTPEnabled() {
		return "http"
	}

	return ""
}

func (server *Server) validateConfig() error {
	err := helper.ValidateConfig(server.config)
	if err != nil {
		return err
	}

	if !server.plainHTTPEnabled() && !server.httpsEnabled() {
		return errors.New("http port or tls must be enabled")
	}

	if server.httpsEnabled() {
		err = server.validateHTTPSConfig()
		if err != nil {
			return errors.New("tls config error, reason: " + err.Error())
		}
	}

	return nil
}

func (server *Server) validateHTTPSConfig() error {
	if server.httpsPort() == "" {
		return errors.New("tls port is required")
	}
	if server.config.TLS.Port != "" && server.config.TLS.Port == server.config.Port {
		return errors.New("tls port must be different from http port")
	}
	if server.config.TLS.CertFile == "" {
		return errors.New("tls cert_file is required")
	}
	if server.config.TLS.KeyFile == "" {
		return errors.New("tls key_file is required")
	}

	return nil
}
