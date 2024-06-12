package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/GusevGrishaEm1/security-service/internal/config"
	"github.com/GusevGrishaEm1/security-service/internal/server"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func main() {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := initLogger(config.Env)

	db, err := sql.Open("sqlite3", config.StoragePath)
	if err != nil {
		panic(err)
	}
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}
	err = goose.Up(db, "migrations")
	if err != nil {
		panic(err)
	}

	logger.Info("Starting service", slog.Any("config", config))

	if err := server.StartServer(context.Background(), logger, config); err != nil {
		panic(err)
	}
}

func initLogger(env string) *slog.Logger {
	switch env {
	case "dev":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return nil
	}
}
