package auth

import (
	"context"
	"log/slog"
	"time"

	securityservicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
	"github.com/GusevGrishaEm1/security-service/internal/config"
	"github.com/GusevGrishaEm1/security-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	config  *config.Config
	storage Storage
	securityservicev1.UnimplementedAuthServer
	logger *slog.Logger
}

//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=auth
type Storage interface {
	FindUserByEmail(ctx context.Context, email string) (model.User, error)
	SaveUser(ctx context.Context, user model.User) error
}

func NewAuthService(ctx *config.Config, storage Storage, logger *slog.Logger) *Service {
	return &Service{
		config:  ctx,
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) Login(ctx context.Context, req *securityservicev1.LoginRequest) (*securityservicev1.LoginResponse, error) {
	s.logger.Info(req.Email)
	user, err := s.storage.FindUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}
	token, err := s.createToken(user.Email)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &securityservicev1.LoginResponse{
		Token: token,
	}, nil
}

func (s *Service) Register(ctx context.Context, req *securityservicev1.RegisterRequest) (*securityservicev1.RegisterResponse, error) {
	s.logger.Info(req.Email)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	user := model.User{
		Email:    req.Email,
		Password: string(passwordHash),
	}
	err = s.storage.SaveUser(ctx, user)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	token, err := s.createToken(user.Email)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	return &securityservicev1.RegisterResponse{
		Token: token,
	}, nil
}

func (s *Service) Logout(context.Context, *securityservicev1.LogoutRequest) (*securityservicev1.LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (s *Service) Refresh(context.Context, *securityservicev1.RefreshRequest) (*securityservicev1.RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}

func (s *Service) createToken(email string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(s.config.TokenTTL).Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
