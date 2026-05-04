---
name: sql-patterns
description: Padrões de query SQL com sqlx e pgx no gymflow-api
---

## Setup da conexão

```go
// internal/infra/db/postgres.go
func NewPostgres(databaseURL string) (*sqlx.DB, error) {
    db, err := sqlx.Connect("pgx", databaseURL)
    if err != nil {
        return nil, fmt.Errorf("connecting to postgres: %w", err)
    }
    return db, nil
}
```

## Padrão de repositório

```go
type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
    return &userRepository{db: db}
}
```

## Insert

```go
func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
    query := `
        INSERT INTO users (id, name, email, created_at)
        VALUES (:id, :name, :email, :created_at)
    `
    _, err := r.db.NamedExecContext(ctx, query, u)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return domain.ErrAlreadyExists
        }
        return fmt.Errorf("inserting user: %w", err)
    }
    return nil
}
```

## Select por ID

```go
func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    var u domain.User
    err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE id = $1`, id)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, domain.ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("finding user: %w", err)
    }
    return &u, nil
}
```

## Select lista

```go
func (r *userRepository) List(ctx context.Context) ([]*domain.User, error) {
    var users []*domain.User
    err := r.db.SelectContext(ctx, &users, `SELECT * FROM users ORDER BY created_at DESC`)
    if err != nil {
        return nil, fmt.Errorf("listing users: %w", err)
    }
    return users, nil
}
```

## Migrations

Sempre criar par `.up.sql` e `.down.sql`:

```sql
-- migrations/001_create_users.up.sql
CREATE TABLE users (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- migrations/001_create_users.down.sql
DROP TABLE IF EXISTS users;
```

## Boas práticas

- Use `NamedExecContext` para inserts/updates com structs — evita erros de ordem de parâmetros
- Use `GetContext` para buscar um registro — já retorna `sql.ErrNoRows` quando não encontrado
- Use `SelectContext` para listas — nunca retorna `ErrNoRows`, apenas slice vazio
- Sempre mapeie `sql.ErrNoRows` para `domain.ErrNotFound` no repositório
- Sempre mapeie `pgErr.Code == "23505"` (unique violation) para `domain.ErrAlreadyExists`
- Wrapeie erros com `fmt.Errorf("contexto: %w", err)` para preservar a cadeia
