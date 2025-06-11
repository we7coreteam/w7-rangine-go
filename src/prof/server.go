package prof

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/server"
	"github.com/we7coreteam/w7-rangine-go/v3/src/core/helper"
	"net/http"
	"net/http/pprof"
)

type Server struct {
	server.ServerInterface

	config Config
	server *http.ServeMux
	routes []string
}

func NewProfServer(config Config) *Server {
	return &Server{
		server: http.NewServeMux(),
		config: config,
	}
}

func (server *Server) handleFunc(pattern string, handler http.HandlerFunc) {
	server.server.HandleFunc(pattern, handler)
	server.routes = append(server.routes, pattern)
}

func (server *Server) registerRoutes() {
	server.handleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(server.routes)
	})
	server.handleFunc("/debug/pprof/", pprof.Index)
	server.handleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	server.handleFunc("/debug/pprof/profile", pprof.Profile)
	server.handleFunc("/debug/pprof/symbol", pprof.Symbol)
	server.handleFunc("/debug/pprof/trace", pprof.Trace)
}

func (server *Server) GetServerName() string {
	return "prof"
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
		panic(errors.New("prof server config error, reason: " + err.Error()))
	}

	server.registerRoutes()

	err = http.ListenAndServe(fmt.Sprintf("%s:%s", server.config.Host, server.config.Port), server.server)
	if err != nil {
		panic(err)
	}
}
