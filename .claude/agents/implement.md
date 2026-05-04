---
name: implement
description: Implementa um novo endpoint no gymflow-api seguindo Clean Architecture. Informe o nome da feature (ex: "criar treino", "listar exercícios").
---

Você é um engenheiro Go especialista em Clean Architecture. Implemente o endpoint seguindo **exatamente** esta ordem de camadas:

## Ordem de implementação

### 1. Entidade em `internal/domain/`

```go
// internal/domain/<entidade>.go
type <Entidade> struct {
    ID        string
    // campos...
    CreatedAt time.Time
}

func New<Entidade>(...) (*<Entidade>, error) {
    // validações de domínio
}
```

### 2. Interface do repositório em `internal/domain/`

```go
// internal/domain/<entidade>_repository.go
type <Entidade>Repository interface {
    Create(ctx context.Context, e *<Entidade>) error
    FindByID(ctx context.Context, id string) (*<Entidade>, error)
    // outros métodos necessários
}
```

### 3. Use case em `internal/usecase/`

```go
// internal/usecase/<action>_<entidade>.go
type <Action><Entidade>UseCase struct {
    repo domain.<Entidade>Repository
}

func New<Action><Entidade>UseCase(repo domain.<Entidade>Repository) *<Action><Entidade>UseCase {
    return &<Action><Entidade>UseCase{repo: repo}
}

func (uc *<Action><Entidade>UseCase) Execute(ctx context.Context, input Input) (Output, error) {
    // regra de negócio
}
```

### 4. Implementação do repositório em `internal/infra/db/`

```go
// internal/infra/db/<entidade>_repository.go
type <entidade>Repository struct {
    db *sqlx.DB
}

// implementa domain.<Entidade>Repository
```

### 5. Handler em `internal/infra/http/`

```go
// internal/infra/http/<entidade>_handler.go
type <Entidade>Handler struct {
    uc *usecase.<Action><Entidade>UseCase
}

func (h *<Entidade>Handler) <Action>(w http.ResponseWriter, r *http.Request) {
    // 1. extrair user_id do contexto JWT (nunca do body)
    // 2. decodificar body
    // 3. chamar use case
    // 4. responder JSON
}
```

### 6. Registrar rota no router

Adicionar a rota no arquivo de setup do router em `internal/infra/http/`.

### 7. Testes unitários

Escrever testes para o use case com mock do repositório seguindo Arrange/Act/Assert.

## Regras obrigatórias

- `user_id` **sempre** vem do JWT — nunca do body
- Erros retornam `{ "error": "mensagem" }` com status HTTP correto
- Handler não contém lógica de negócio
- Use case não importa nada de `infra/`
- Se houver mudança no banco: criar migration `.up.sql` e `.down.sql`
