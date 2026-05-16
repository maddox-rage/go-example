package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
)

const (
	envLocal = "local"
	envProd = "prod"
	envDev = "dev"
)

func main() {
	cfg := config.MustLoad()
	
	log := setUpLogger(cfg.Env)

	log.Info("starting", slog.String("env", cfg.Env))
	log.Debug("debug  msg are enable")
}

func setUpLogger(env string) *slog.Logger{
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}