package handlers

import (
	"log"
	"net/http"
	"products/internal/config"
	"products/internal/services"
	"products/pkg/products/domain"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	service *services.Product
}

func NewProduct(config *config.Config) *Product {
	return &Product{
		service: services.NewProduct(config),
	}
}

func (p *Product) GetById(ctx *gin.Context) {
	start := time.Now()

	id := ctx.Param("id")

	log.Printf("[handlers.Product] Getting product with ID: %s", id)

	product, err := p.service.GetById(id)

	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"product":   product,
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (p *Product) Get(ctx *gin.Context) {
	start := time.Now()

	log.Printf("[handlers.Product] Getting all products")

	products, err := p.service.Get()

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"products":  products,
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (p *Product) Insert(ctx *gin.Context) {
	start := time.Now()

	var product *domain.Product

	log.Printf("[handlers.Product] Inserting product")

	err := ctx.BindJSON(&product)

	if err != nil {
		body := ctx.Request.Body
		defer body.Close()

		log.Printf("[handlers.Product] Error binding JSON: %s -> Body: %s", err, body)

		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	id, err := p.service.Insert(product)

	if err != nil {
		log.Printf("[handlers.Product] Error inserting product: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	log.Printf("[handlers.Product] Product inserted with ID: %s -> Product: %+v", id, product)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"id":        id,
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (p *Product) Update(ctx *gin.Context) {
	start := time.Now()

	log.Printf("[handlers.Product] Updating product")

	var product *domain.Product

	err := ctx.BindJSON(&product)

	if err != nil {
		body := ctx.Request.Body
		defer body.Close()

		log.Printf("[handlers.Product] Error binding JSON: %s -> Body: %s", err, body)

		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	err = p.service.Update(product)

	if err != nil {
		log.Printf("[handlers.Product] Error updating product: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	log.Printf("[handlers.Product] Product updated with ID: %s -> Product: %+v", product.ID, product)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (p *Product) Delete(ctx *gin.Context) {
	start := time.Now()

	id := ctx.Param("id")

	log.Printf("[handlers.Product] Deleting product with ID: %s", id)

	err := p.service.Delete(id)

	if err != nil {
		log.Printf("[handlers.Product] Error deleting product: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
