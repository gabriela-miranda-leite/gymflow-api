# GYM-53 — POST /auth/login + Middleware JWT

## O que foi implementado

Endpoint de login que autentica o usuário e retorna dois tokens: um **access token** (JWT) de curta duração e um **refresh token** de longa duração. Também foi criado o middleware JWT que protege rotas autenticadas.

---

## Fluxo completo

```
POST /auth/login
  → AuthHandler.Login         (lê body, valida input)
  → LoginUseCase.Execute      (regras de negócio)
    → UserRepository.FindByEmail   (busca usuário no banco)
    → bcrypt.CompareHashAndPassword (compara senha com hash)
    → jwt.Generate             (gera access token)
    → crypto/rand              (gera refresh token aleatório)
    → sha256                   (hasha o refresh token)
    → RefreshTokenRepository.Create (salva hash no banco)
  → resposta 200 com os dois tokens
```

---

## Arquivos criados

### `internal/domain/refresh_token.go`
Entidade que representa um refresh token no sistema. Espelha a tabela `refresh_tokens` do banco.

```go
type RefreshToken struct {
    ID        string    `db:"id"`
    UserID    string    `db:"user_id"`
    TokenHash string    `db:"token_hash"`
    ExpiresAt time.Time `db:"expires_at"`
    Revoked   bool      `db:"revoked"`
    CreatedAt time.Time `db:"created_at"`
}
```

### `internal/domain/refresh_token_repository.go`
Interface (contrato) que define como acessar a tabela `refresh_tokens`. O use case depende dessa interface, não da implementação concreta.

```go
type RefreshTokenRepository interface {
    Create(ctx context.Context, token *RefreshToken) error
    FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
    Revoke(ctx context.Context, id string) error
}
```

### `pkg/jwt/jwt.go`
Utilitário de geração e validação de JWT. Fica em `pkg/` porque é usado por dois lugares: o `LoginUseCase` (geração) e o `JWTMiddleware` (validação).

**Generate** — cria um token com:
- `sub`: ID do usuário
- `exp`: expiração em 15 minutos
- Algoritmo: HS256 (HMAC + SHA256)
- Chave: `JWT_SECRET` do ambiente

**Validate** — valida o token e retorna o `userID`:
- Verifica a assinatura com o `JWT_SECRET`
- Verifica que o algoritmo é HMAC (proteção contra ataque do algoritmo `none`)
- Verifica a expiração
- Retorna o `userID` do claim `sub`

### `internal/usecase/login.go`
Regras de negócio do login. Passos do `Execute`:

1. Busca o usuário pelo email
2. Se não encontrar → `ErrInvalidCredentials`
3. Compara a senha com bcrypt
4. Se não bater → `ErrInvalidCredentials` (mesma mensagem — não revela qual campo está errado)
5. Gera o access token JWT
6. Gera 32 bytes aleatórios com `crypto/rand` → converte pra hex (64 chars) → esse é o refresh token
7. Hasha o refresh token com SHA256 → salva o hash no banco (nunca o token em si)
8. Retorna os dois tokens + dados do usuário

**Por que SHA256 no refresh token e bcrypt na senha?**
O bcrypt é propositalmente lento (custo 12 ≈ 250ms por operação) — ótimo para senha porque dificulta brute force. Mas o refresh token seria validado frequentemente, então usar bcrypt seria impraticável. O SHA256 é rápido e suficiente aqui porque o refresh token já é um valor aleatório de 32 bytes (impossível de adivinhar por brute force).

### `internal/infra/db/refresh_token_repository.go`
Implementação concreta do `RefreshTokenRepository` com SQL via sqlx.

- **Create** — INSERT com `user_id`, `token_hash` e `expires_at` (banco gera `id` e `created_at`)
- **FindByHash** — SELECT por `token_hash`, trata `sql.ErrNoRows` como `nil, nil`
- **Revoke** — UPDATE `revoked = true` pelo `id`

### `internal/infra/http/auth_handler.go`
Adicionado método `Login` ao `AuthHandler` existente. O handler:
1. Decodifica o body JSON
2. Chama o `LoginUseCase`
3. Se `ErrInvalidCredentials` → 401
4. Se outro erro → 500 (nunca expõe detalhes internos)
5. Se sucesso → 200 com os tokens e dados do usuário

### `internal/infra/http/middleware/jwt.go`
Middleware que protege rotas autenticadas.

**JWT(next http.Handler) http.Handler**
- Lê o header `Authorization: Bearer <token>`
- Se ausente ou sem prefixo `Bearer` → 401
- Valida o token com `pkg/jwt.Validate`
- Se inválido ou expirado → 401
- Se válido → injeta o `userID` no contexto e passa pro próximo handler

**GetUserIDFromContext(ctx)**
- Helper para recuperar o `userID` do contexto em qualquer handler protegido
- Usa tipo customizado `contextKey` como chave (evita colisão com outros pacotes)

---

## Tabela `refresh_tokens`

```sql
CREATE TABLE refresh_tokens (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id),
    token_hash TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked    BOOLEAN     NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
```

O índice em `user_id` acelera buscas futuras de todos os tokens de um usuário (ex: revogar todos ao fazer logout em todos os dispositivos).

---

## Request / Response

```
POST /auth/login
Content-Type: application/json

Body:
{
    "email": "gabriela@email.com",
    "password": "senha123"
}

200 OK:
{
    "access_token": "eyJhbGci...",
    "refresh_token": "0909a52b...",
    "user": {
        "id": "00b924fb-...",
        "name": "Gabriela",
        "email": "gabriela@email.com"
    }
}

401 Unauthorized:
{ "error": "invalid credentials" }
```

---

## Segurança

| Decisão | Motivo |
|---|---|
| Mesma mensagem para email não encontrado e senha errada | Evita que atacante descubra emails cadastrados |
| Hash do refresh token com SHA256 | Nunca salvar o token real — se o banco vazar, tokens são inúteis |
| Validação do algoritmo JWT | Proteção contra ataque do algoritmo `none` |
| `JWT_SECRET` obrigatório | Sem a variável de ambiente, a aplicação falha explicitamente em vez de assinar com chave vazia |
| Access token expira em 15 min | Limita a janela de uso de um token roubado |

---

## Refresh Token Rotation (GYM-54)

O refresh token expira em 7 dias, mas o app nunca vai deslogar o usuário enquanto ele usar o app regularmente. Isso será implementado no GYM-54:

1. Cliente manda o refresh token para `POST /auth/refresh`
2. Servidor valida o token e verifica que não está revogado
3. Revoga o token antigo (`revoked = true`)
4. Gera novo access token + novo refresh token com prazo renovado
5. Retorna os dois tokens novos

Resultado: enquanto o usuário abrir o app pelo menos a cada 7 dias, nunca é deslogado.

---

## Testes

### `internal/usecase/login_test.go`
- `TestLogin_Success` — email existe, senha correta → retorna tokens e dados do usuário
- `TestLogin_EmailNotFound` — email não cadastrado → `ErrInvalidCredentials`
- `TestLogin_WrongPassword` — email existe mas senha errada → `ErrInvalidCredentials`

### `internal/infra/http/middleware/jwt_test.go`
- `TestJWTMiddleware_ValidToken` — token válido → passa pro handler e injeta `userID` no contexto
- `TestJWTMiddleware_MissingToken` — sem header → 401
- `TestJWTMiddleware_MalformedToken` — token inválido → 401

**Cobertura:**
- `internal/usecase`: 81.1%
- `internal/infra/http/middleware`: 100%
- `internal/domain`: 100%
