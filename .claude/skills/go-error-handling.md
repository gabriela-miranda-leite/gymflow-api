---
name: go-error-handling
description: Padrões de tratamento de erro do gymflow-api
---

## Padrão de resposta de erro HTTP

Todo handler deve retornar JSON com a chave `error`:

```go
func writeError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
```

## Mapeamento de erros de domínio → status HTTP

| Situação | Status |
|---|---|
| Input inválido | 400 Bad Request |
| Não autenticado | 401 Unauthorized |
| Sem permissão | 403 Forbidden |
| Entidade não encontrada | 404 Not Found |
| Conflito (ex: email já existe) | 409 Conflict |
| Erro interno / banco | 500 Internal Server Error |

## Erros de domínio tipados

Defina erros sentinela em `domain/` para que o usecase e o handler possam identificá-los:

```go
// internal/domain/errors.go
var (
    ErrNotFound      = errors.New("not found")
    ErrAlreadyExists = errors.New("already exists")
    ErrUnauthorized  = errors.New("unauthorized")
)
```

No handler, use `errors.Is` para mapear ao status correto:

```go
if errors.Is(err, domain.ErrNotFound) {
    writeError(w, http.StatusNotFound, "recurso não encontrado")
    return
}
if errors.Is(err, domain.ErrAlreadyExists) {
    writeError(w, http.StatusConflict, "recurso já existe")
    return
}
// fallback
writeError(w, http.StatusInternalServerError, "erro interno")
```

## Nunca expor erros internos

```go
// errado — vaza detalhes do banco
writeError(w, 500, err.Error())

// correto — mensagem genérica para o cliente, log interno
log.Printf("db error: %v", err)
writeError(w, 500, "erro interno")
```
