package server

import (
	"currencies/internal/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine  *gin.Engine
	srv     *http.Server
	address string
	port    int

	config          *Config
	currencyHandler *handlers.Currency
}

func NewServer() *Server {
	return &Server{
		engine: gin.Default(),
	}
}

func (s *Server) Run() error {
	s.srv = &http.Server{
		Addr:    s.makeAddress(),
		Handler: s.engine,
	}

	log.Printf("[server] starting server on %s", s.makeAddress())

	gin.ForceConsoleColor()
	gin.DefaultWriter = log.Writer()
	gin.DefaultErrorWriter = log.Writer()
	if s.config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s.engine = gin.Default()
	s.engine.GET("/currency/:code", s.currencyHandler.GetByCode)
	s.engine.GET("/currency/", s.currencyHandler.Get)
	s.engine.POST("/currency/update", s.currencyHandler.Update)

	log.Print("Router started")

	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.srv.Close()
}

func (s *Server) makeAddress() string {
	return fmt.Sprintf("%s:%d", s.address, s.port)
}
