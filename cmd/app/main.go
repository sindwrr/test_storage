package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/sindwrr/test_storage/internal/api"
	"github.com/sindwrr/test_storage/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env found. Using system env instead.")
	}

	cfg := config.Load()
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open DB! Err: %s", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB! Err: %s", err)
	}

	log.Print("Connection with DB established!")

	router := api.NewRouter(cfg)
	log.Println("Server starting on :8000")
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
