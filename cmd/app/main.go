package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/sindwrr/test_storage/internal/api"
	"github.com/sindwrr/test_storage/internal/config"
)

// @title           Test Storage API
// @version         1.0
// @description     API для системы хранения артефактов автотестов.
// @host            localhost:8000
// @BasePath        /
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

	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to setup migrations! Err: %s", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations! Err: %s", err)
	}
	log.Print("Applied migrations successfully!")

	router := api.NewRouter(db, cfg)
	log.Println("Server starting on :8000")
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
