package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger/v2"
	"log/slog"
	"net/http"
	"os"
	_ "url-shortener/docs"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/storage/postgres"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is url-shortener project
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		url-shortener.swagger.io
//	@BasePath	/v2

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// Init config
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("Starting URL-shortener")

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		logger.Error("Failed to init database")
		logger.Error(err.Error())
		os.Exit(1)
	}

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Only for chi
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(logger, storage))
	router.Get("/url/{alias}", redirect.New(logger, storage))
	router.Delete("/url/{alias}", delete.New(logger, storage))

	// Docs
	router.Mount("/docs", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	logger.Info("Started URL-shortener")

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server")
	}

}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {

	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}

	return logger
}
