package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gabriela-miranda-leite/gymflow-api/internal/infra/http/middleware"
	pkgjwt "github.com/gabriela-miranda-leite/gymflow-api/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	_ = os.Setenv("JWT_SECRET", "test-secret")
}

func nextHandler(t *testing.T, expectUserID string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserIDFromContext(r.Context())
		assert.Equal(t, expectUserID, userID)
		w.WriteHeader(http.StatusOK)
	})
}

func TestJWTMiddleware_ValidToken(t *testing.T) {
	// Arrange
	userID := "user-id-123"
	token, err := pkgjwt.Generate(userID)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	// Act
	middleware.JWT(nextHandler(t, userID)).ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestJWTMiddleware_MissingToken(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	// Act
	middleware.JWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestJWTMiddleware_MalformedToken(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer token.invalido.aqui")
	rec := httptest.NewRecorder()

	// Act
	middleware.JWT(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
