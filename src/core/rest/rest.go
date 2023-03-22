package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"helicopter/core"
)

type Rest struct {
	storage core.Storage
}

func NewRest(storage core.Storage) *Rest {
	return &Rest{storage: storage}
}

func (r *Rest) GetNodes(c *gin.Context) {
	root, err := strconv.ParseInt(c.Param("root"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid query parameters")
		return
	}
	last, err := strconv.ParseInt(c.Param("last"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid query parameters")
		return
	}

	nodes, _ := r.storage.GetSubTreeNodes(root, last)
	c.JSON(http.StatusOK, nodes)
}

func (r *Rest) AddNode(c *gin.Context) {
	var requestBody core.Node

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body or parameters")
		return
	}

	node, _ := r.storage.CreateNode(requestBody.Parent, requestBody.Value)

	c.JSON(http.StatusCreated, node)
}
