package usecase_test

import (
	"context"
	"testing"

	"github.com/gabriela-miranda-leite/gymflow-api/internal/domain"
	"github.com/gabriela-miranda-leite/gymflow-api/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepository struct {
	findByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	createFn      func(ctx context.Context, user *domain.User) error
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.findByEmailFn(ctx, email)
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.createFn(ctx, user)
}

func TestRegisterUser_Success(t *testing.T) {
	// Arrange
	repo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		},
		createFn: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}
	uc := usecase.NewRegisterUserUseCase(repo)
	input := usecase.RegisterUserInput{
		Name:     "Gabriela",
		Email:    "gabriela@email.com",
		Password: "senha123",
	}

	// Act
	_, err := uc.Execute(context.Background(), input)

	// Assert
	require.NoError(t, err)
}

func TestRegisterUser_EmailAlreadyInUse(t *testing.T) {
	// Arrange
	repo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return &domain.User{}, nil
		},
		createFn: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}
	uc := usecase.NewRegisterUserUseCase(repo)
	input := usecase.RegisterUserInput{
		Name:     "Gabriela",
		Email:    "gabriela@email.com",
		Password: "senha123",
	}

	// Act
	_, err := uc.Execute(context.Background(), input)

	// Assert
	assert.ErrorIs(t, err, usecase.ErrEmailAlreadyInUse)
}

func TestRegisterUser_WeakPassword(t *testing.T) {
	// Arrange
	repo := &mockUserRepository{
		findByEmailFn: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, nil
		},
		createFn: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}
	uc := usecase.NewRegisterUserUseCase(repo)
	input := usecase.RegisterUserInput{
		Name:     "Gabriela",
		Email:    "gabriela@email.com",
		Password: "123",
	}

	// Act
	_, err := uc.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.EqualError(t, err, "password must be at least 6 characters")

}
