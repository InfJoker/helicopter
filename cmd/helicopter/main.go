package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"helicopter/internal/config"
	"helicopter/internal/grpc"
	"helicopter/internal/mockstorage"
	"helicopter/internal/rest"

	"golang.org/x/sync/errgroup"
)

// Print a help message
func printHelp() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s [-c|--config <file>] [-h|--help]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	configPtr := flag.String("config", "", "the configuration file")
	helpPtr := flag.Bool("help", false, "print help message")

	flag.Parse()

	if *helpPtr || *configPtr == "" {
		printHelp()
		return
	}
	cfg, err := config.NewConfig(*configPtr)

	if err != nil {
		log.Println("Error parsing config: ", err)
		os.Exit(1)
	}

	storage, err := mockstorage.NewStorage(0)
	if err != nil {
		log.Println("Error while creating mockstorage:", err)
		os.Exit(1)
	}

	g, _ := errgroup.WithContext(context.Background())

	rest := rest.NewRest(cfg, storage)
	g.Go(func() error {
		if err := rest.Run(); err != nil {
			return fmt.Errorf("error while running rest server: %v", err)
		}
		return nil
	})
	rpc := grpc.NewGrpc(cfg, storage)
	g.Go(func() error {
		if err := rpc.Run(); err != nil {
			return fmt.Errorf("error while running rpc server: %v", err)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
