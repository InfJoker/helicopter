package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"helicopter/internal/core"
)

type Rest struct {
	storage core.Storage
	router  *gin.Engine
}

func NewRest(storage core.Storage) *Rest {
	rest := &Rest{
		storage: storage,
		router:  gin.Default(),
	}
	rest.router.GET("/nodes", rest.GetNodes)
	rest.router.POST("/nodes", rest.AddNode)

	return rest
}

func (r *Rest) Run(host string) error {
	return r.router.Run(host)
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
