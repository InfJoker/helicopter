package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"helicopter/internal/config"
	"helicopter/internal/core"
)

type Rest struct {
	address string
	storage core.Storage
	router  *gin.Engine
}

func NewRest(cfg config.Config, storage core.Storage) *Rest {
	host, port := "localhost", 8080
	if cfg.HttpServer.Host != "" {
		host = cfg.HttpServer.Host
	}
	if cfg.HttpServer.Port != 0 {
		port = cfg.HttpServer.Port
	}
	address := fmt.Sprintf("%s:%d", host, port)

	rest := &Rest{
		address: address,
		storage: storage,
		router:  gin.Default(),
	}
	rest.router.GET("/nodes", rest.GetNodes)
	rest.router.POST("/nodes", rest.AddNode)

	return rest
}

func (r *Rest) Run() error {
	return r.router.Run(r.address)
}

func (r *Rest) GetNodes(c *gin.Context) {
	root, last := c.Param("root"), c.Param("last")

	nodes, _ := r.storage.GetSubTreeNodes(c.Request.Context(), root, last)
	c.JSON(http.StatusOK, nodes)
}

func (r *Rest) AddNode(c *gin.Context) {
	var requestBody core.Node

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body or parameters")
		return
	}

	node, _ := r.storage.CreateNode(c.Request.Context(), requestBody.Parent, requestBody.Value)

	c.JSON(http.StatusCreated, node)
}
