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
	rest := new(Rest)
	rest.storage = storage
	return rest
}

func (r *Rest) GetNodeChildren(c *gin.Context) {
	lseq, err := strconv.ParseInt(c.Param("lseq"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid query parameters")
		return
	}

	if lseq == 0 {
		c.JSON(http.StatusForbidden, "Cannot retrieve children of the 0 node")
		return
	}

	children, _ := r.storage.GetSubTreeNodes(lseq, lseq)
	c.JSON(http.StatusOK, children)
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

type PostNodesRequestBody struct {
	Ref     int64  `json:"ref"`
	Content []byte `json:"content"`
}

func (r *Rest) AddNode(c *gin.Context) {
	var requestBody PostNodesRequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body or parameters")
	}

	node, _ := r.storage.CreateNode(requestBody.Ref, requestBody.Content)

	c.JSON(http.StatusCreated, node)
}
