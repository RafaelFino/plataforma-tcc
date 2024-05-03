package server

import (
	"clients/internal/config"
	"clients/internal/handlers"
	"clients/internal/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	srv    *http.Server

	config  *config.Config
	handler *handlers.Client
	service *services.Client
}

func NewServer(config *config.Config) *Server {
	s := &Server{
		engine:  gin.Default(),
		config:  config,
		service: services.NewProduct(config),
	}

	log.Printf("[server] Starting server with config: %+v", config)

	s.handler = handlers.NewProduct(config)

	gin.ForceConsoleColor()
	gin.DefaultWriter = log.Writer()
	gin.DefaultErrorWriter = log.Writer()

	if s.config.Debug {
		log.Printf("[server] Debug mode enabled")
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine = gin.Default()
	s.engine.GET("/products/:id", s.handler.GetById)
	s.engine.GET("/products/", s.handler.Get)
	s.engine.POST("/products/", s.handler.Insert)
	s.engine.PUT("/products/", s.handler.Update)

	s.srv = &http.Server{
		Addr:    s.makeAddress(),
		Handler: s.engine,
	}

	return s
}

func (s *Server) Run() {
	log.Printf("[server] starting server on %s", s.makeAddress())
	err := s.srv.ListenAndServe()
	if err != nil {
		log.Printf("[server] error starting server: %s", err)
		panic(err)
	}
}

func (s *Server) Stop() error {
	return s.srv.Close()
}

func (s *Server) makeAddress() string {
	return fmt.Sprintf("%s:%d", s.config.ServerAddress, s.config.ServerPort)
}
