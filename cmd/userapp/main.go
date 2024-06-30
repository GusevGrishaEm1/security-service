package main

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/GusevGrishaEm1/security-service/internal/config"
	"github.com/GusevGrishaEm1/security-service/internal/server"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

func main() {
	c, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := initLogger()

	db, err := sql.Open("sqlite3", c.StoragePath)
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

	logger.Info("Starting service", slog.Any("config", c))

	if err := server.StartServer(logger, c); err != nil {
		panic(err)
	}
}

func initLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
