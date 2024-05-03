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

	config         *config.Config
	clientsHandler *handlers.Clients
	clientsService *services.Clients
}

func NewServer(config *config.Config) *Server {
	s := &Server{
		engine:         gin.Default(),
		config:         config,
		clientsService: services.NewClients(),
	}

	log.Printf("[server] Starting server with config: %+v", config)

	s.clientsHandler = handlers.NewClients()

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
	s.engine.GET("/clients/:id", s.clientsHandler.GetById)
	s.engine.GET("/clients/", s.clientsHandler.Get)
	s.engine.POST("/clients/", s.clientsHandler.Insert)
	s.engine.PUT("/clients/:id", s.clientsHandler.Update)

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
