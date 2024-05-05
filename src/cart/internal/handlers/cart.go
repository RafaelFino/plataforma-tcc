package handlers

import (
	"cart/internal/config"
	"cart/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Cart struct {
	service *services.Cart
}

func NewCart(config *config.Config) *Cart {
	return &Cart{
		service: services.NewCart(config),
	}
}

func (c *Cart) CreateCart(ctx *gin.Context) {
	start := time.Now()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Cart) Get(ctx *gin.Context) {
	start := time.Now()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Cart) DeleteCart(ctx *gin.Context) {
	start := time.Now()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Cart) AddProduct(ctx *gin.Context) {
	start := time.Now()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Cart) RemoveProduct(ctx *gin.Context) {
	start := time.Now()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Cart) GetByClient(ctx *gin.Context) {
	start := time.Now()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Cart) Checkout(ctx *gin.Context) {
	start := time.Now()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
