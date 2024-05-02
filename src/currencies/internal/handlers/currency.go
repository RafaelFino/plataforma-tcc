package handlers

import (
	"currency/internal/services"
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

	log.Printf("[handlers.Currency] Update")

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message":   "Updated",
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Currency) Get(ctx *gin.Context) {
	start := time.Now()

	ret, err := c.service.Get()

	if err != nil {
		log.Printf("[handlers.Currency] Error getting currencies: %s", err.Error())
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	body := make(map[string]string)

	for k, v := range ret {
		body[k] = v.ToJSON()
	}

	log.Printf("[handlers.Currency] currencies: %v", body)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"currencies": body,
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
		"currency":  ret.ToJSON(),
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
