package api

import (
	"net/http"

	"github.com/sindwrr/test_storage/internal/api/handlers"
	"github.com/sindwrr/test_storage/internal/auth"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	authService := auth.NewService()
	loginHandler := handlers.NewLoginHandler(authService)

	mux.HandleFunc("/login", loginHandler.Handle)
	mux.HandleFunc("/", handlers.HelloHandler)

	return mux
}
