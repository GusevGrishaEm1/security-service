package server

import (
	"log/slog"
	"net"
	"strconv"

	securityservicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
	"github.com/GusevGrishaEm1/security-service/internal/config"
	storage "github.com/GusevGrishaEm1/security-service/internal/storage/sqllite"
	"github.com/GusevGrishaEm1/security-service/internal/usecase/auth"
	"google.golang.org/grpc"
)

func StartServer(logger *slog.Logger, config *config.Config) error {
	// init listener
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
	if err != nil {
		return err
	}

	// init authStorage
	logger.Info("Initializing authStorage")
	authStorage, err := storage.NewAuthStorage(config)
	if err != nil {
		return err
	}

	// init service
	logger.Info("Initializing service")
	service := auth.NewAuthService(config, authStorage, logger)

	// init server
	logger.Info("Initializing server")
	s := grpc.NewServer()

	// register service
	securityservicev1.RegisterAuthServer(s, service)

	// start server
	logger.Info("Start server")
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
