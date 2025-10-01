package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"project/internal/config"
	"project/internal/lib/logger/sl"

	"project/internal/storage/sqlite"

	"project/internal/http-server/handlers/delete"
	"project/internal/http-server/handlers/redirect"
	"project/internal/http-server/handlers/url/save"
	"project/internal/http-server/middleware"
	"project/internal/http-server/middleware/mwlogger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)

	log.Info("starting url-shortner", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		slog.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router: gin
	router := gin.New()

	// middleware
	router.Use(middleware.RequestID())
	router.Use(mwlogger.Logger())
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recovery())
	router.Use(middleware.URLFormat())

	router.POST("/url", save.New(log, storage))
	router.GET("/{alias}", redirect.New(log, storage))
	router.DELETE("/url/{alias}", delete.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
