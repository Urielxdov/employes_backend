package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		"employees",
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Retry connection for 30s (Docker startup timing)
	for i := 0; i < 60; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("[DB] Connected successfully")
			return nil
		}
		log.Printf("[DB] Connection attempt %d failed, retrying...\n", i+1)
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("failed to connect to database after retries")
}

func closeDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
