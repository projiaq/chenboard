package main

import (
	"flag"
	"fmt"
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
	pgPort := flag.String("pg-port", "15432", "PostgreSQL protocol port")
	flag.Parse()

	log.Println("=== TSDB Time-Series Database ===")
	log.Printf("Data directory: %s", *dataDir)
	log.Printf("HTTP API port: %s", *httpPort)
	log.Printf("PostgreSQL protocol port: %s", *pgPort)
	log.Println()

	// Initialize storage engine
	log.Println("Initializing storage engine...")
	engine, err := storage.NewEngine(*dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize storage engine: %v", err)
	}
	defer engine.Close()
	log.Println("✓ Storage engine initialized")

	// Start HTTP API server
	httpServer := api.NewServer(engine, *httpPort)
	httpErrChan := make(chan error, 1)
	go func() {
		log.Printf("Starting HTTP API server on http://localhost:%s", *httpPort)
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			httpErrChan <- fmt.Errorf("HTTP server error: %v", err)
		}
	}()

	// Start PostgreSQL protocol server
	pgServer := pgserver.NewPgServer(engine, *pgPort)
	pgErrChan := make(chan error, 1)
	go func() {
		log.Printf("Starting PostgreSQL protocol server on localhost:%s", *pgPort)
		if err := pgServer.Start(); err != nil {
			pgErrChan <- fmt.Errorf("PostgreSQL server error: %v", err)
		}
	}()

	// Wait for startup errors or shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println()
	log.Println("✓ All servers started successfully")
	log.Println("Press Ctrl+C to shutdown")
	log.Println()

	select {
	case err := <-httpErrChan:
		log.Fatalf("HTTP server failed: %v", err)
	case err := <-pgErrChan:
		log.Fatalf("PostgreSQL server failed: %v", err)
	case <-sigChan:
		log.Println()
		log.Println("Shutting down servers...")
	}

	httpServer.Stop()
	pgServer.Stop()
	log.Println("✓ Shutdown complete")
}
