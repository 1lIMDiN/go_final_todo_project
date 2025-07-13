package main

import (
	"log"
	"os"
	"strconv"

	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	dbFile := "scheduler.db"
	if eDBFile := os.Getenv("TODO_DBFILE"); eDBFile != "" {
		dbFile = eDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.DB.Close()

	port := 7540
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err != nil {
			port = p
		}
	}

	if err := server.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
