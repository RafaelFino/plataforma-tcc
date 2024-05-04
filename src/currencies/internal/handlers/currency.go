package handlers

import (
	"currencies/internal/domain"
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

	data := c.service.Get()
	ret := make([]*domain.CurrencyData, 0)

	for _, d := range data {
		item := d.GetData()

		if item != nil {
			ret = append(ret, item)
		}
	}

	log.Printf("[handlers.Currency] All currencies: %+v", ret)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"currencies": ret,
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
		"currency":  ret.GetData(),
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
