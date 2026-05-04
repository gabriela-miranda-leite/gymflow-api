# gymflow-api

API backend do GymFlow construída em Go com Clean Architecture.

## Pré-requisitos

- [Go 1.22+](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/)
- [golangci-lint](https://golangci-lint.run/usage/install/) (para lint)
- [golang-migrate](https://github.com/golang-migrate/migrate) (para migrations)

## Setup local

1. Clone o repositório:
   ```bash
   git clone https://github.com/gabriela-miranda-leite/gymflow-api.git
   cd gymflow-api
   ```

2. Copie o arquivo de variáveis de ambiente:
   ```bash
   cp .env.example .env
   ```
   Edite `.env` com suas configurações locais.

3. Inicie o servidor:
   ```bash
   make run
   ```

## Comandos disponíveis

| Comando           | Descrição                                 |
|-------------------|-------------------------------------------|
| `make run`        | Inicia o servidor em modo desenvolvimento |
| `make build`      | Compila o binário em `bin/api`            |
| `make test`       | Executa os testes                         |
| `make lint`       | Executa o linter                          |
| `make migrate-up` | Aplica as migrations pendentes            |
| `make seed`       | Popula o banco com dados iniciais         |

## Estrutura de pastas

```
cmd/api/          ← entry point (main.go)
internal/
  domain/         ← entidades e interfaces (zero deps externos)
  usecase/        ← regras de negócio
  infra/
    http/         ← handlers, router, middleware
    db/           ← queries SQL, repositórios
pkg/              ← utilitários reutilizáveis (jwt, hash, validator)
migrations/       ← arquivos .sql versionados
```
