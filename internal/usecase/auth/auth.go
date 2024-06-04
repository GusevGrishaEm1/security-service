package auth

import (
	"context"
	"time"

	security_servicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
	"github.com/GusevGrishaEm1/security-service/internal/config"
	"github.com/GusevGrishaEm1/security-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	config  *config.Config
	storage AuthStorage
	security_servicev1.UnimplementedAuthServer
}

//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=auth
type AuthStorage interface {
	FindUserByEmail(ctx context.Context, email string) (model.User, error)
	SaveUser(ctx context.Context, user model.User) error
}

func NewAuthService(ctx *config.Config, storage AuthStorage) *AuthService {
	return &AuthService{
		config:  ctx,
		storage: storage,
	}
}

func (s *AuthService) Login(ctx context.Context, req *security_servicev1.LoginRequest) (*security_servicev1.LoginResponse, error) {
	user, err := s.storage.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}
	token, err := s.createToken(user.Email)
	if err != nil {
		return nil, err
	}
	return &security_servicev1.LoginResponse{
		Token: token,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *security_servicev1.RegisterRequest) (*security_servicev1.RegisterResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := model.User{
		Email:    req.Email,
		Password: string(passwordHash),
	}
	err = s.storage.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	token, err := s.createToken(user.Email)
	if err != nil {
		return nil, err
	}
	return &security_servicev1.RegisterResponse{
		Token: token,
	}, nil
}

func (s *AuthService) Logout(context.Context, *security_servicev1.LogoutRequest) (*security_servicev1.LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (s *AuthService) Refresh(context.Context, *security_servicev1.RefreshRequest) (*security_servicev1.RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}

func (s *AuthService) createToken(email string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": email,
			"exp":   time.Now().Add(time.Hour * time.Duration(s.config.TokenTTL)).Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(s.config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
