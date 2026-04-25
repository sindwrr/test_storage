package api

import (
	"database/sql"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/sindwrr/test_storage/docs"
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
	mux.HandleFunc("/health/alive", handlers.AliveHandler(healthSvc))
	mux.HandleFunc("/health/ready", handlers.ReadyHandler(healthSvc))
	mux.HandleFunc("/", middleware.RequireAuth(handlers.HelloHandler))

	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	return middleware.CorsMiddleware(mux)
}
