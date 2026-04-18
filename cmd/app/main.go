package main

import (
	"log"
	"net/http"

	"github.com/sindwrr/test_storage/internal/api"
)

func main() {
	router := api.NewRouter()

	log.Println("Server starting on :8000")

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
