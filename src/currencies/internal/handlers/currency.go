package handlers

import (
	"currencies/internal/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Currency struct {
	service *services.Currency
}

func NewCurrency(service *services.Currency) *Currency {
	return &Currency{
		service: service,
	}
}

func (c *Currency) Update(ctx *gin.Context) {
	start := time.Now()

	log.Printf("[handlers.Currency] Updating...")

	err := c.service.Update()

	if err != nil {
		log.Printf("[handlers.Currency] Error updating currency: %s", err.Error())
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"last_update": c.service.LastUpdate().Format(time.RFC3339),
		"elapsed":     time.Since(start).Milliseconds(),
		"timestamp":   time.Now().Format(time.RFC3339),
	})
}

func (c *Currency) Get(ctx *gin.Context) {
	start := time.Now()

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"currencies": c.service.Get(),
		"elapsed":    time.Since(start).Milliseconds(),
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

func (c *Currency) GetByCode(ctx *gin.Context) {
	start := time.Now()

	code := ctx.Param("code")

	ret, err := c.service.GetByCode(code)

	if err != nil {
		log.Printf("[handlers.Currency] Error getting currency: %s (%s)", err.Error(), code)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	log.Printf("[handlers.Currency] currency: %s -> %s", code, ret.ToJSON())

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"currency":  ret,
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
