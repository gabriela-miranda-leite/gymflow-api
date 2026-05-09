package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	db_infra "github.com/gabriela-miranda-leite/gymflow-api/internal/infra/db"
	http_infra "github.com/gabriela-miranda-leite/gymflow-api/internal/infra/http"
	"github.com/gabriela-miranda-leite/gymflow-api/internal/usecase"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	userRepo := db_infra.NewUserRepository(db)
	registerUC := usecase.NewRegisterUserUseCase(userRepo)
	authHandler := http_infra.NewAuthHandler(registerUC)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("POST /auth/register", authHandler.Register)

	log.Printf("server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
