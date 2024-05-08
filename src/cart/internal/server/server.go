package server

import (
	"cart/internal/config"
	"cart/internal/handlers"
	"cart/internal/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	srv    *http.Server

	config  *config.Config
	handler *handlers.Cart
	service *services.Cart
}

func NewServer(config *config.Config) *Server {
	s := &Server{
		engine:  gin.Default(),
		config:  config,
		service: services.NewCart(config),
	}

	log.Printf("[server] Starting server with config: %+v", config)

	s.handler = handlers.NewCart(config)

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

	s.engine.POST("/cart/", s.handler.CreateCart)
	s.engine.GET("/cart/:cart_id", s.handler.Get)
	s.engine.DELETE("/cart/:cart_id", s.handler.DeleteCart)

	s.engine.POST("/cart/:cart_id", s.handler.AddProduct)
	s.engine.DELETE("/cart/:cart_id/:product_id", s.handler.RemoveProduct)

	s.engine.GET("/cart/client/:client_id", s.handler.GetByClient)
	s.engine.PUT("/cart/:cart_id", s.handler.Checkout)

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
