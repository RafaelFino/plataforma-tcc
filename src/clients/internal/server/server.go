package server

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/rafaelfino/plataforma-tcc/clients/internal/config"
	handlers "github.com/rafaelfino/plataforma-tcc/clients/internal/handlers"
	services "github.com/rafaelfino/plataforma-tcc/clients/internal/services"

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
		service: services.NewClient(config),
	}

	log.Printf("[server] Starting server with config: %+v", config)

	s.handler = handlers.NewClient(config)

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
	s.engine.GET("/clients/:id", s.handler.GetById)
	s.engine.GET("/clients/", s.handler.Get)
	s.engine.POST("/clients/", s.handler.Insert)
	s.engine.PUT("/clients/", s.handler.Update)
	s.engine.DELETE("/clients/", s.handler.Delete)

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
	log.Printf("[server] stopping service")
	err := s.service.Close()

	if err != nil {
		log.Printf("[server] error stopping server: %s", err)
	}

	err = s.srv.Close()

	if err != nil {
		log.Printf("[server] error stopping server: %s", err)
	}

	return err
}

func (s *Server) makeAddress() string {
	return fmt.Sprintf("%s:%d", s.config.ServerAddress, s.config.ServerPort)
}
