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

	config          *config.Config
	currencyHandler *handlers.Currency
	currencyService *services.Currency
}

func NewServer(config *config.Config) *Server {
	s := &Server{
		engine:          gin.Default(),
		config:          config,
		currencyService: services.NewCurrency(config.CurrencyURL),
	}

	log.Printf("[server] Starting server with config: %+v", config)

	s.currencyHandler = handlers.NewCurrency(s.currencyService)

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
	s.engine.GET("/currency/:code", s.currencyHandler.GetByCode)
	s.engine.GET("/currency/", s.currencyHandler.Get)
	s.engine.POST("/currency/update", s.currencyHandler.Update)

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
