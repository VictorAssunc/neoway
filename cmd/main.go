package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"neoway/pkg/handler"
	"neoway/pkg/lib/file"
	"neoway/pkg/repository"
	"neoway/pkg/service"
)

const (
	maxRetries    = 10
	retryInterval = time.Second
)

func main() {
	ctx := context.Background()

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DB")
	password := os.Getenv("POSTGRES_PASSWORD")

	connInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatalf("[ERROR] Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Constant backoff retry
	for _ = range maxRetries {
		err = db.Ping()
		if err == nil {
			log.Println("[INFO] Connected to database :D")
			break
		}

		log.Println("[INFO] Failed to connect to database. Retrying...")
		time.Sleep(retryInterval)
	}
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect to database: %v", err)
	}

	log.Println("[INFO] Processing...")

	clientRepo := repository.NewClient(db)
	clientService := service.NewClient(clientRepo)
	clientHandler := handler.NewClient(clientService)

	totalStart := time.Now()
	filename := "base.txt"
	baseFile := file.New(filename)
	if err = baseFile.Open(); err != nil {
		log.Fatalf("[ERROR] Failed to open file: %v", err)
	}
	defer baseFile.Close()

	readStart := time.Now()
	lines, err := baseFile.Read()
	if err != nil {
		log.Fatalf("[ERROR] Failed to read file: %v", err)
	}
	readDuration := time.Since(readStart)

	storeStart := time.Now()
	if err = clientHandler.CreateClientsFromLines(ctx, lines); err != nil {
		log.Fatalf("[ERROR] Failed to create clients: %v", err)
	}
	storeDuration := time.Since(storeStart)

	normalizeStart := time.Now()
	if err = clientHandler.NormalizeClients(ctx); err != nil {
		log.Fatalf("[ERROR] Failed to normalize clients: %v", err)
	}
	normalizeDuration := time.Since(normalizeStart)
	totalDuration := time.Since(totalStart)

	log.Printf(`[INFO] Took %s to read, store and normalize %d clients
  - Read time: %s
  - Store time: %s
  - Normalize time: %s
`, totalDuration.String(), len(lines), readDuration.String(), storeDuration.String(), normalizeDuration.String())
}
