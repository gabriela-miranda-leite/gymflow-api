package domain

import "errors"

type User struct {
	ID    string
	Name  string
	Email string
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
