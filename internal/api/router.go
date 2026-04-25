package api

import (
	"database/sql"
	"net/http"

	"github.com/sindwrr/test_storage/internal/api/handlers"
	"github.com/sindwrr/test_storage/internal/api/middleware"
	"github.com/sindwrr/test_storage/internal/auth"
	"github.com/sindwrr/test_storage/internal/health"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	healthSvc := health.NewService(db)
	authSvc := auth.NewService()
	loginHandler := handlers.NewLoginHandler(authSvc)

	mux.HandleFunc("/login", loginHandler.Handle)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	mux.HandleFunc("/", middleware.RequireAuth(handlers.HelloHandler))
	mux.HandleFunc("/health/alive", handlers.AliveHandler(healthSvc))
	mux.HandleFunc("/health/ready", handlers.ReadyHandler(healthSvc))

	return mux
}
