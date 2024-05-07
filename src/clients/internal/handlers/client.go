package handlers

import (
	config "github.com/rafaelfino/plataforma-tcc/src/clients/internal/config"
	services "github.com/rafaelfino/plataforma-tcc/src/clients/internal/services"
	domain "github.com/rafaelfino/plataforma-tcc/src/clients/pkg/domain"

	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Client struct {
	service *services.Client
}

func NewClient(config *config.Config) *Client {
	return &Client{
		service: services.NewClient(config),
	}
}

func (c *Client) GetById(ctx *gin.Context) {
	start := time.Now()

	id := ctx.Param("id")

	log.Printf("[handlers.Client] Getting client with ID: %s", id)

	client, err := c.service.GetById(id)

	if err != nil {
		log.Printf("[handlers.Client] Error getting client: %s", err)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"client":    client,
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Client) Get(ctx *gin.Context) {
	start := time.Now()

	log.Printf("[handlers.Client] Getting all clients")

	clients, err := c.service.Get()

	if err != nil {
		log.Printf("[handlers.Client] Error getting clients: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})

		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"clients":   clients,
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Client) Insert(ctx *gin.Context) {
	start := time.Now()

	log.Printf("[handlers.Client] Inserting client")

	client := &domain.Client{}

	err := ctx.BindJSON(client)

	if err != nil {
		body := ctx.Request.Body
		defer body.Close()

		log.Printf("[handlers.Client] Error binding JSON: %s -> Body: %s", err, body)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	id, err := c.service.Insert(client)

	if err != nil {
		log.Printf("[handlers.Client] Error inserting client: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	log.Printf("[handlers.Client] Client inserted with ID: %s -> Client: %+v", id, client)

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"id":        id,
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Client) Update(ctx *gin.Context) {
	start := time.Now()

	log.Printf("[handlers.Client] Updating client")

	client := &domain.Client{}

	err := ctx.BindJSON(client)

	if err != nil {
		body := ctx.Request.Body
		defer body.Close()

		log.Printf("[handlers.Client] Error binding JSON: %s -> Body: %s", err, body)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})

		return
	}

	err = c.service.Update(client)

	if err != nil {
		log.Printf("[handlers.Client] Error updating client: %s", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
			"elapsed":   time.Since(start).Milliseconds(),
			"error":     err.Error(),
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return

	}

	log.Printf("[handlers.Client] Client updated with ID: %s -> Client: %+v", client.ID, client)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"elapsed":   time.Since(start).Milliseconds(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func (c *Client) Delete(ctx *gin.Context) {
	start := time.Now()

	id := ctx.Param("id")

	log.Printf("[handlers.Client] Deleting client with ID: %s", id)

	err := c.service.Delete(id)

	if err != nil {
		log.Printf("[handlers.Client] Error deleting client: %s", err)
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
