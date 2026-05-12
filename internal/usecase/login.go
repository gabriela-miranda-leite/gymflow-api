package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/gabriela-miranda-leite/gymflow-api/internal/domain"
	"github.com/gabriela-miranda-leite/gymflow-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

const refreshTokenExpiration = 7 * 24 * time.Hour

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	Name         string
	Email        string
}

type LoginUseCase struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
}

func NewLoginUseCase(userRepo domain.UserRepository, refreshTokenRepo domain.RefreshTokenRepository) *LoginUseCase {
	return &LoginUseCase{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := jwt.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return nil, err
	}
	rawTokenStr := hex.EncodeToString(rawToken)

	hash := sha256.Sum256([]byte(rawTokenStr))
	tokenHash := hex.EncodeToString(hash[:])

	refreshToken := &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(refreshTokenExpiration),
	}
	if err := uc.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: rawTokenStr,
		UserID:       user.ID,
		Name:         user.Name,
		Email:        user.Email,
	}, nil
}
