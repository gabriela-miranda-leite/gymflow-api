# gymflow-api — Contexto do Projeto

## Stack

| Camada | Tecnologia |
|---|---|
| Linguagem | Go 1.23 |
| Router HTTP | chi |
| Banco de dados | PostgreSQL |
| Query builder | sqlx + pgx |
| Autenticação | golang-jwt |
| Testes | testify |
| Hot reload | air |
| Migrations | golang-migrate |
| Lint | golangci-lint v2 |

## Arquitetura: Clean Architecture

```
domain → usecase → infra
```

- **domain**: entidades e interfaces de repositório. Zero dependências externas. Nunca importa de `infra` ou `usecase`.
- **usecase**: regras de negócio. Depende apenas de interfaces definidas em `domain`. Nunca importa de `infra/http` ou `infra/db` diretamente.
- **infra/http**: handlers HTTP, router, middleware. Chama use cases.
- **infra/db**: implementações concretas dos repositórios. Usa sqlx/pgx para acessar o banco.
- **pkg**: utilitários reutilizáveis (jwt, hash, validator). Sem regra de negócio.

## Estrutura de pastas

```
cmd/api/            ← entry point (main.go)
internal/
  domain/           ← entidades e interfaces (zero deps externos)
  usecase/          ← regras de negócio
  infra/
    http/           ← handlers, router, middleware
    db/             ← queries SQL, repositórios
pkg/                ← utilitários reutilizáveis (jwt, hash, validator)
migrations/         ← arquivos .sql versionados (.up.sql e .down.sql)
```

## Fluxo de dados

```
Request HTTP
  → Handler (infra/http)   — valida input, extrai user_id do JWT
  → UseCase (usecase)      — executa a regra de negócio
  → Repository interface (domain) — chamada ao banco via infra/db
  → Response JSON
```

## Regras de camada

- `domain` nunca importa de `infra`, `usecase` ou pacotes externos além da stdlib
- `usecase` nunca importa de `infra/http` ou `infra/db` — usa apenas interfaces
- `infra/http` nunca contém lógica de negócio — só validação de input e serialização
- **`user_id` nunca vem do body da requisição** — sempre extraído do JWT no middleware

## Naming conventions

| Tipo | Padrão | Exemplo |
|---|---|---|
| Use case | `<Action><Entity>UseCase` | `CreateUserUseCase` |
| Repository interface | `<Entity>Repository` | `UserRepository` |
| Repository impl | `<Entity>Repository` (em `infra/db`) | `userRepository` |
| Handler | `<Entity>Handler` | `UserHandler` |
| Arquivo de teste | `<arquivo>_test.go` | `user_test.go` |

## Comandos do Makefile

```bash
make run          # inicia o servidor (hot reload com air em dev)
make build        # compila o binário em bin/api
make test         # roda testes e gera coverage.html
make lint         # roda golangci-lint
make migrate-up   # aplica migrations pendentes
make seed         # popula o banco com dados iniciais
```

## Variáveis de ambiente

```env
DATABASE_URL=postgres://user:password@localhost:5432/gymflow?sslmode=disable
JWT_SECRET=your-secret-key-here
PORT=8080
```

## Padrão de erro

Todo erro deve retornar JSON com o status HTTP correto:

```go
// handler
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{"error": "mensagem legível"})
```

Nunca expor erros internos (ex: mensagens do banco) para o cliente.

## Padrão de teste: Arrange / Act / Assert

```go
func TestNomeDoComportamento(t *testing.T) {
    // Arrange — prepara dados e mocks
    // Act     — executa a função sendo testada
    // Assert  — verifica resultado com testify
    assert.Equal(t, expected, got)
}
```

Use `require.NoError(t, err)` antes de acessar o resultado para evitar panic em caso de erro.
