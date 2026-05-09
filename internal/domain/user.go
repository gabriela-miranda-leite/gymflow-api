package domain

import (
	"errors"
	"time"
)

type User struct {
	ID               string    `db:"id"`
	Name             string    `db:"name"`
	Email            string    `db:"email"`
	Phone            *string   `db:"phone"`
	PasswordHash     string    `db:"password_hash"`
	IdealTimeEnabled bool      `db:"ideal_time_enabled"`
	OccupancyLimit   *string   `db:"occupancy_limit"`
	CreatedAt        time.Time `db:"created_at"`
}

func NewUser(id, name, email string) (*User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	return &User{ID: id, Name: name, Email: email}, nil
}
