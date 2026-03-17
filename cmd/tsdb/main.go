package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/chenboard/tsdb/internal/api"
	"github.com/chenboard/tsdb/internal/pgserver"
	"github.com/chenboard/tsdb/internal/storage"
)

func main() {
	dataDir := flag.String("data", "./data", "Data directory")
	httpPort := flag.String("http-port", "6041", "HTTP API port")
	pgPort := flag.String("pg-port", "5432", "PostgreSQL protocol port")
	flag.Parse()

	engine, err := storage.NewEngine(*dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize storage engine: %v", err)
	}
	defer engine.Close()

	// Start HTTP API server
	httpServer := api.NewServer(engine, *httpPort)
	go func() {
		log.Printf("Starting HTTP API server on port %s", *httpPort)
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start PostgreSQL protocol server
	pgServer := pgserver.NewPgServer(engine, *pgPort)
	go func() {
		log.Printf("Starting PostgreSQL protocol server on port %s", *pgPort)
		if err := pgServer.Start(); err != nil {
			log.Fatalf("PostgreSQL server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down servers...")
	httpServer.Stop()
	pgServer.Stop()
}
