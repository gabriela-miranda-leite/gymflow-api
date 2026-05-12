package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gabriela-miranda-leite/gymflow-api/internal/usecase"
)

type AuthHandler struct {
	registerUser *usecase.RegisterUserUseCase
	loginUser    *usecase.LoginUseCase
}

func NewAuthHandler(registerUser *usecase.RegisterUserUseCase, loginUser *usecase.LoginUseCase) *AuthHandler {
	return &AuthHandler{registerUser: registerUser, loginUser: loginUser}
}

type registerUserResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     *string `json:"phone"`
	CreatedAt string  `json:"created_at"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input usecase.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	user, err := h.registerUser.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyInUse) {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "email already in use"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(registerUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input usecase.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	output, err := h.loginUser.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"access_token":  output.AccessToken,
		"refresh_token": output.RefreshToken,
		"user": map[string]string{
			"id":    output.UserID,
			"name":  output.Name,
			"email": output.Email,
		},
	})
}
