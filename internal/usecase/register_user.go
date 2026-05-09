package usecase

import (
	"context"
	"errors"

	"github.com/gabriela-miranda-leite/gymflow-api/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyInUse = errors.New("email already in use")

type RegisterUserInput struct {
	Name     string
	Email    string
	Password string
	Phone    *string
}

type RegisterUserUseCase struct {
	repo domain.UserRepository
}

func NewRegisterUserUseCase(repo domain.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{repo: repo}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, input RegisterUserInput) (*domain.User, error) {
	if len(input.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	existingUser, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrEmailAlreadyInUse
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return nil, err
	}

	id := uuid.NewString()

	user := domain.User{
		ID:               id,
		Name:             input.Name,
		Email:            input.Email,
		Phone:            input.Phone,
		PasswordHash:     string(hash),
		IdealTimeEnabled: false,
	}

	if err := uc.repo.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
