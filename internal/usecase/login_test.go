package usecase_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gabriela-miranda-leite/gymflow-api/internal/domain"
	"github.com/gabriela-miranda-leite/gymflow-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockRefreshTokenRepository struct {
	createFn     func(ctx context.Context, token *domain.RefreshToken) error
	findByHashFn func(ctx context.Context, hash string) (*domain.RefreshToken, error)
	revokeFn     func(ctx context.Context, id string) error
}

func (m *mockRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	return m.createFn(ctx, token)
}

func (m *mockRefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	return m.findByHashFn(ctx, hash)
}

func (m *mockRefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	return m.revokeFn(ctx, id)
}

func init() {
	_ = os.Setenv("JWT_SECRET", "test-secret")
}

func validPasswordHash() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), 12)
	return string(hash)
}

func TestLogin_Success(t *testing.T) {
	// Arrange
	passwordHash := validPasswordHash()
	userRepo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:           "user-id-123",
				Name:         "Gabriela",
				Email:        email,
				PasswordHash: passwordHash,
				CreatedAt:    time.Now(),
			}, nil
		},
		createFn: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}
	refreshTokenRepo := &mockRefreshTokenRepository{
		createFn: func(ctx context.Context, token *domain.RefreshToken) error {
			return nil
		},
		findByHashFn: func(ctx context.Context, hash string) (*domain.RefreshToken, error) {
			return nil, nil
		},
		revokeFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
	uc := usecase.NewLoginUseCase(userRepo, refreshTokenRepo)
	input := usecase.LoginInput{
		Email:    "gabriela@email.com",
		Password: "senha123",
	}

	// Act
	output, err := uc.Execute(context.Background(), input)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, output.AccessToken)
	assert.NotEmpty(t, output.RefreshToken)
	assert.Equal(t, "user-id-123", output.UserID)
	assert.Equal(t, "Gabriela", output.Name)
}

func TestLogin_EmailNotFound(t *testing.T) {
	// Arrange
	userRepo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		},
		createFn: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}
	refreshTokenRepo := &mockRefreshTokenRepository{
		createFn: func(ctx context.Context, token *domain.RefreshToken) error {
			return nil
		},
		findByHashFn: func(ctx context.Context, hash string) (*domain.RefreshToken, error) {
			return nil, nil
		},
		revokeFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
	uc := usecase.NewLoginUseCase(userRepo, refreshTokenRepo)
	input := usecase.LoginInput{
		Email:    "naoexiste@email.com",
		Password: "senha123",
	}

	// Act
	_, err := uc.Execute(context.Background(), input)

	// Assert
	assert.ErrorIs(t, err, usecase.ErrInvalidCredentials)
}

func TestLogin_WrongPassword(t *testing.T) {
	// Arrange
	passwordHash := validPasswordHash()
	userRepo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{
				ID:           "user-id-123",
				Name:         "Gabriela",
				Email:        email,
				PasswordHash: passwordHash,
				CreatedAt:    time.Now(),
			}, nil
		},
		createFn: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}
	refreshTokenRepo := &mockRefreshTokenRepository{
		createFn: func(ctx context.Context, token *domain.RefreshToken) error {
			return nil
		},
		findByHashFn: func(ctx context.Context, hash string) (*domain.RefreshToken, error) {
			return nil, nil
		},
		revokeFn: func(ctx context.Context, id string) error {
			return nil
		},
	}
	uc := usecase.NewLoginUseCase(userRepo, refreshTokenRepo)
	input := usecase.LoginInput{
		Email:    "gabriela@email.com",
		Password: "senhaerrada",
	}

	// Act
	_, err := uc.Execute(context.Background(), input)

	// Assert
	assert.ErrorIs(t, err, usecase.ErrInvalidCredentials)
}
