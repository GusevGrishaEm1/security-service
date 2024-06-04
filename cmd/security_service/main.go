package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/GusevGrishaEm1/security-service/internal/config"
	"github.com/GusevGrishaEm1/security-service/internal/server"
)

func main() {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := initLogger(config.Env)

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
