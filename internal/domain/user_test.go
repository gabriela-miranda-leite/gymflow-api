package domain_test

import (
	"testing"

	"github.com/gabriela-miranda-leite/gymflow-api/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewUser_Success(t *testing.T) {
	// Arrange
	id, name, email := "1", "Alice", "alice@example.com"

	// Act
	user, err := domain.NewUser(id, name, email)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, id, user.ID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, email, user.Email)
}

func TestNewUser_MissingName(t *testing.T) {
	// Arrange / Act
	user, err := domain.NewUser("1", "", "alice@example.com")

	// Assert
	assert.Nil(t, user)
	assert.EqualError(t, err, "name is required")
}

func TestNewUser_MissingEmail(t *testing.T) {
	// Arrange / Act
	user, err := domain.NewUser("1", "Alice", "")

	// Assert
	assert.Nil(t, user)
	assert.EqualError(t, err, "email is required")
}
