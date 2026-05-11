package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("[INIT] Connecting to database...")
	if err := initDB(); err != nil {
		log.Fatalf("[FATAL] Failed to initialize database: %v", err)
	}
	defer closeDB()

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[INFO] Starting server on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), setupRouter(db)); err != nil {
		log.Fatalf("[FATAL] Server failed: %v", err)
	}
}
