package auth

import (
	"context"
	"log/slog"
	"testing"

	security_servicev1 "github.com/GusevGrishaEm1/protos/gen/go/security_service"
	"github.com/GusevGrishaEm1/security-service/internal/config"
	"github.com/GusevGrishaEm1/security-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type AuthStorageMock struct {
	mock.Mock
}

func (s *AuthStorageMock) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	args := s.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}

func (s *AuthStorageMock) SaveUser(ctx context.Context, user model.User) error {
	args := s.Called(ctx, user)
	return args.Error(0)
}

func TestLogin(t *testing.T) {
	config := &config.Config{}
	config.SecretKey = "test"
	config.TokenTTL = 1
	hashpassword, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	storage := &AuthStorageMock{}
	storage.On("FindUserByEmail", mock.Anything, "test").Return(model.User{
		Email:    "test",
		Password: string(hashpassword),
	}, nil)
	service := NewAuthService(config, storage, slog.Default())

	res, err := service.Login(context.Background(), &security_servicev1.LoginRequest{
		Email:    "test",
		Password: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, res.Token)
}

func TestRegister(t *testing.T) {
	config := &config.Config{}
	config.SecretKey = "test"
	config.TokenTTL = 1
	storage := &AuthStorageMock{}
	storage.On("SaveUser", mock.Anything, mock.Anything).Return(nil)
	service := NewAuthService(config, storage, slog.Default())

	res, err := service.Register(context.Background(), &security_servicev1.RegisterRequest{
		Email:    "test",
		Password: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, res.Token)
}
