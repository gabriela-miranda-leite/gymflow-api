package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gabriela-miranda-leite/gymflow-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type refreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	_, err := r.db.ExecContext(ctx, `
	INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`, token.UserID, token.TokenHash, token.ExpiresAt)
	return err
}

func (r *refreshTokenRepository) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	err := r.db.GetContext(ctx, &token, `SELECT * FROM refresh_tokens WHERE token_hash = $1`, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked = true WHERE id = $1`, id)
	return err
}
