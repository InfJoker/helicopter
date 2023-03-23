package main

import (
	"log"
	"os"

	"helicopter/internal/mockstorage"
	"helicopter/internal/rest"
)

func main() {
	storage, err := mockstorage.NewStorage(0)
	if err != nil {
		log.Println("Error while creating mockstorage:", err)
		os.Exit(1)
	}

	rest := rest.NewRest(storage)
	err = rest.Run("localhost:8080")
	if err != nil {
		log.Println("Error while running http server:", err)
		os.Exit(1)
	}
}
