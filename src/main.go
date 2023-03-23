package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"helicopter/core/mockstorage"
	"helicopter/core/rest"
)

func main() {
	storage, err := mockstorage.NewStorage(0)
	if err != nil {
		log.Println("Error while creating mockstorage:", err)
		os.Exit(1)
	}
	rest := rest.NewRest(storage)

	router := gin.Default()
	router.GET("/nodes", rest.GetNodes)
	router.POST("/nodes", rest.AddNode)

	router.Run("localhost:8080")
}
