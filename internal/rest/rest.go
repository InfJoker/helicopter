package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"helicopter/internal/config"
	"helicopter/internal/core"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Rest struct {
	address string
	storage core.Storage
	router  *gin.Engine
}

func NewRest(cfg config.Config, storage core.Storage) (*Rest, error) {
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

	openapi_template_path := "./static"
	if cfg.OpenapiTemplate.Path != "" {
		openapi_template_path = cfg.OpenapiTemplate.Path
	}

	data, err := ioutil.ReadFile(openapi_template_path + "/openapi_template.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading file data: %v", err)
	}

	updatedData := strings.ReplaceAll(string(data), "{host}", host)
	updatedData = strings.ReplaceAll(updatedData, "{port}", strconv.Itoa(port))

	err = ioutil.WriteFile(openapi_template_path+"/openapi.yaml", []byte(updatedData), 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing to file: %v", err)
	}

	rest.router.StaticFile("/swagger-ui/doc.yaml", openapi_template_path+"/openapi.yaml")
	url := ginSwagger.URL("/swagger-ui/doc.yaml") // The url pointing to API definition
	rest.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return rest, nil
}

func (r *Rest) Run() error {
	return r.router.Run(r.address)
}

func (r *Rest) GetNodes(c *gin.Context) {
	root, last := c.Query("root"), c.Query("last")

	nodes, err := r.storage.GetSubTreeNodes(c.Request.Context(), root, last)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, nodes)
}

func (r *Rest) AddNode(c *gin.Context) {
	var requestBody core.Node

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body or parameters")
		return
	}

	node, err := r.storage.CreateNode(c.Request.Context(), requestBody.Parent, requestBody.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusCreated, node)
}
