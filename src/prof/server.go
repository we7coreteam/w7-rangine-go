package prof

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/we7coreteam/w7-rangine-go-support/src/server"
	"net/http"
	"net/http/pprof"
)

type Server struct {
	server.Server

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
	err := binding.Validator.ValidateStruct(server.config)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errStr := "prof server config error, fields: "
			for _, e := range validationErrors {
				errStr += e.Field() + ";"
			}
			panic(errStr)
		} else {
			panic(err)
		}
	}

	server.registerRoutes()

	err = http.ListenAndServe(fmt.Sprintf("%s:%s", server.config.Host, server.config.Port), server.server)
	if err != nil {
		panic(err)
	}
}
