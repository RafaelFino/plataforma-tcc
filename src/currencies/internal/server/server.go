package server

import (
	"currencies/internal/config"
	"currencies/internal/handlers"
	"currencies/internal/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	srv    *http.Server

	config  *config.Config
	handler *handlers.Currency
	service *services.Currency
}

func NewServer(config *config.Config) *Server {
	s := &Server{
		engine:  gin.Default(),
		config:  config,
		service: services.NewCurrency(config.CurrencyURL),
	}

	log.Printf("[server] Starting server with config: %+v", config)

	s.handler = handlers.NewCurrency(s.service)

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
	s.engine.GET("/currency/:code", s.handler.GetByCode)
	s.engine.GET("/currency/", s.handler.Get)
	s.engine.POST("/currency/", s.handler.Update)
	s.engine.PUT("/currency/", s.handler.Update)

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

	err := s.srv.Close()

	if err != nil {
		log.Printf("[server] error stopping server: %s", err)
	}

	return err
}

func (s *Server) makeAddress() string {
	return fmt.Sprintf("%s:%d", s.config.ServerAddress, s.config.ServerPort)
}
