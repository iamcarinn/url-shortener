package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/storage"
	"url-shortener/internal/storage/memory"
	"url-shortener/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	// init config
	cfg := config.MustLoad()
	//fmt.Println(cfg)

	// init logger
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

	// init router
	router := chi.NewRouter()
	// middleware
	router.Use(middleware.RequestID)	// добавляем к каждому запросу id, исп. для логгирования
	router.Use(middleware.Logger)		// логгируем каждый запрос (*исп-ет свой логгер)
	router.Use(middleware.Recoverer)	// если на сервере паника, восстанавливаем
	router.Use(middleware.URLFormat)	// парсер urlов поступающих запросов
	
	router.Post("/url", save.New(log, st))
	log.Info("starting server", slog.String("address", cfg.Storage.Type))

	// init server
	srv := &http.Server{
		Addr: cfg.HTTPServer.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}
	// run server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("fail start server")
	}

	log.Error("server stopped")

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