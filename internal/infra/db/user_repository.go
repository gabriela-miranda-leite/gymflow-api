package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gabriela-miranda-leite/gymflow-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, phone, password_hash, ideal_time_enabled, occupancy_limit)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, user.ID, user.Name, user.Email, user.Phone, user.PasswordHash, user.IdealTimeEnabled, user.OccupancyLimit)
	return err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.GetContext(ctx, &user, `SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
