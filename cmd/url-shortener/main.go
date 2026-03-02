package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/postgres"
	"url-shortener/internal/storage/memory"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	// init config: cleanenv
	cfg := config.MustLoad()
	//fmt.Println(cfg)

	// init logger: slog
	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init storage
	var st storage.Storage
	switch cfg.Storage.Type {
	case "memory":
		st = memory.New()
	case "postgres":
		dsn := cfg.Postgres.DSN()
		pg, err := postgres.New(dsn)
		if err != nil {
			log.Error("fail init postgres storage", slog.String("dsn", dsn), slog.Any("error", err))
			os.Exit(1)
		}
		st = pg
	}
	_ = st

	// init router: chi
	router := chi.NewRouter()
	// middleware
	router.Use(middleware.RequestID)	// добавляем к каждому запросу id, исп. для логгирования
	router.Use(middleware.Logger)		// логгируем каждый запрос (*исп-ет свой логгер)
	router.Use(middleware.Recoverer)	// если на сервере паника, восстанавливаем
	router.Use(middleware.URLFormat)	// парсер urlов поступающих запросов
	
	// TODO: run server
}

// Логгирование, так как его установка зависит от пар-ра env
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