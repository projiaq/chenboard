package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chenboard/tsdb/internal/api"
	"github.com/chenboard/tsdb/internal/storage"
)

func main() {
	dataDir := flag.String("data", "./data", "Data directory")
	port := flag.String("port", "6041", "HTTP port")
	flag.Parse()

	engine, err := storage.NewEngine(*dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize storage engine: %v", err)
	}
	defer engine.Close()

	server := api.NewServer(engine, *port)

	go func() {
		log.Printf("Starting TSDB server on port %s", *port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	server.Stop()
}
