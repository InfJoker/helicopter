package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"helicopter/internal/grpc"
	"helicopter/internal/mockstorage"
	"helicopter/internal/rest"

	"golang.org/x/sync/errgroup"
)

func main() {
	storage, err := mockstorage.NewStorage(0)
	if err != nil {
		log.Println("Error while creating mockstorage:", err)
		os.Exit(1)
	}

	g, _ := errgroup.WithContext(context.Background())

	rest := rest.NewRest(storage)
	g.Go(func() error {
		if err := rest.Run("0.0.0.0:8080"); err != nil {
			return fmt.Errorf("error while running rest server: %v", err)
		}
		return nil
	})
	rpc := grpc.NewGrpc(storage)
	g.Go(func() error {
		if err := rpc.Run("0.0.0.0:1337"); err != nil {
			return fmt.Errorf("error while running rpc server: %v", err)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
