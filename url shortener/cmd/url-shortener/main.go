package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/slg"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoad()

	log := setUpLogger(cfg.Env)

	log.Info("starting", slog.String("env", cfg.Env))
	log.Debug("debug  msg are enable")

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to init storage", slg.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	_ = storage

}

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}
