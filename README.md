# gymflow-api

[![CI](https://github.com/gabriela-miranda-leite/gymflow-api/actions/workflows/ci.yml/badge.svg)](https://github.com/gabriela-miranda-leite/gymflow-api/actions/workflows/ci.yml)

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

## Testes

Execute os testes e gere o relatório de cobertura:

```bash
make test
```

Isso gera dois arquivos:
- `coverage.out` — dados brutos de cobertura
- `coverage.html` — relatório visual (abra no browser para ver quais linhas foram testadas)

### Padrão de teste (Arrange / Act / Assert)

```go
func TestNomeDoComportamento(t *testing.T) {
    // Arrange — prepara os dados
    // Act     — executa o que está sendo testado
    // Assert  — verifica o resultado com testify
    assert.Equal(t, expected, got)
}
```

## Pre-commit hook

Instale o hook para que o lint rode automaticamente antes de cada commit:

```bash
pre-commit install
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
