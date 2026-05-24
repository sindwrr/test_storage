package api

import (
	"database/sql"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/sindwrr/test_storage/docs"
	"github.com/sindwrr/test_storage/internal/analytics"
	"github.com/sindwrr/test_storage/internal/api/handlers"
	"github.com/sindwrr/test_storage/internal/api/middleware"
	"github.com/sindwrr/test_storage/internal/auth"
	"github.com/sindwrr/test_storage/internal/config"
	"github.com/sindwrr/test_storage/internal/health"
	"github.com/sindwrr/test_storage/internal/metadata"
	"github.com/sindwrr/test_storage/internal/preview"
	"github.com/sindwrr/test_storage/internal/storage"
)

func NewRouter(db *sql.DB, cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	healthSvc := health.NewService(db)
	authSvc := auth.NewService(cfg.LDAPAddr, cfg.LDAPBaseDN, cfg.LDAPUser, cfg.LDAPPassword, db)

	storageSvc, err := storage.NewStorageService(cfg.ArtifactVolume)
	if err != nil {
		log.Fatalf("Router: failed to create storage service! Err: %s", err)
	}

	metadataSvc := metadata.NewMetadataService(db)
	analyticsSvc := analytics.NewService(db)
	previewSvc := preview.NewService(metadataSvc, storageSvc)

	loginHandler := handlers.NewLoginHandler(authSvc)
	uploadHandler := handlers.NewUploadHandler(storageSvc, metadataSvc, cfg.MaxFileBytes)
	indexHandler := handlers.NewIndexHandler(metadataSvc)
	downloadHandler := handlers.NewDownloadHandler(metadataSvc, storageSvc)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsSvc)
	artifactsHandler := handlers.NewArtifactsHandler(metadataSvc)
	previewHandler := handlers.NewPreviewHandler(previewSvc)
	logoutHandler := handlers.NewLogoutHandler(authSvc)

	mux.HandleFunc("/login", loginHandler.Handle)
	mux.HandleFunc("/logout", logoutHandler.Handle)
	mux.HandleFunc("/health/alive", handlers.AliveHandler(healthSvc))
	mux.HandleFunc("/health/ready", handlers.ReadyHandler(healthSvc))
	mux.HandleFunc("/upload", middleware.RequireAuth(middleware.RequireUpload(authSvc)(uploadHandler.Handle)))

	mux.HandleFunc("/", middleware.RequireAuth(indexHandler.Handle))
	mux.HandleFunc("/artifact/download/{id}", middleware.RequireAuth(downloadHandler.Handle))
	mux.HandleFunc("/artifacts", middleware.RequireAuth(artifactsHandler.List))
	mux.HandleFunc("/analytics/artifacts-per-day", middleware.RequireAuth(analyticsHandler.ArtifactsPerDay))
	mux.HandleFunc("/analytics/status-distribution", middleware.RequireAuth(analyticsHandler.StatusDistribution))
	mux.HandleFunc("/analytics", middleware.RequireAuth(handlers.AnalyticsPageHandler))
	mux.HandleFunc("/preview", middleware.RequireAuth(previewHandler.ServePreview))

	swaggerHandler := httpSwagger.WrapHandler
	adminOnly := middleware.RequireAdmin(authSvc)
	mux.Handle("/docs/", middleware.RequireAuth(adminOnly(swaggerHandler.ServeHTTP)))
	mux.Handle("/docs", middleware.RequireAuth(adminOnly(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	})))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	return middleware.CorsMiddleware(mux)
}
