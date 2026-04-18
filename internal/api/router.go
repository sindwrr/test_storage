package api

import (
	"net/http"

	"github.com/sindwrr/test_storage/internal/api/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	return mux
}
