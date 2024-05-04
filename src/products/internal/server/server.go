package server

import (
	"fmt"
	"log"
	"net/http"
	"products/internal/config"
	"products/internal/handlers"
	"products/internal/services"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	srv    *http.Server

	config  *config.Config
	handler *handlers.Product
	service *services.Product
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
	s.engine.DELETE("/products/", s.handler.Delete)

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
