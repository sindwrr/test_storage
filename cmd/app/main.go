package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/sindwrr/test_storage/internal/api"
	"github.com/sindwrr/test_storage/internal/config"
	"github.com/sindwrr/test_storage/internal/worker"
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

	workerPool := worker.NewPool(4)
	workerPool.Start()

	go func() {
		for {
			time.Sleep(7 * 24 * time.Hour)
			workerPool.Submit(worker.IntegrityCheckTask{
				DB:       db,
				BasePath: cfg.ArtifactVolume,
			})
		}
	}()

	router := api.NewRouter(db, cfg)
	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		log.Println("Server starting on :8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	workerPool.Shutdown()
	log.Println("Server exited")
}
