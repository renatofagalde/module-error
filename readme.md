# 🚨 module-error

Módulo Go para gerenciamento de erros de domínio em APIs REST, com mapeamento automático para HTTP status codes.

## 🎯 Características

- ✅ Erros tipados e reutilizáveis
- ✅ Mapeamento automático para HTTP status codes
- ✅ Cobertura completa de cenários de CRM/SaaS multi-tenant
- ✅ Suporte a idempotência e concorrência
- ✅ Conformidade com padrões REST puristas
- ✅ Testes unitários incluídos
- ✅ ~85% de cobertura dos status HTTP mais comuns

## 📦 Instalação
```bash
go get github.com/renatofagalde/module-error
```

## 🚀 Uso Básico
```go
import "github.com/renatofagalde/module-error"

// Retornar erro de domínio
func GetUser(id string) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, domain_error.ErrNotFound
    }
    return user, nil
}

// Mapear para HTTP status code
func HandleError(w http.ResponseWriter, err error) {
    mapper := domain_error.NewHTTPStatusMapper()
    statusCode := mapper.GetHTTPStatus(err)
    
    if domainErr, ok := err.(*domain_error.DomainError); ok {
        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(domainErr)
    }
}
```

## 📋 Categorias de Erros

### 🔴 Validação e Input (400)
```go
domain_error.ErrInvalidInput
domain_error.ErrInvalidEmail
domain_error.ErrInvalidCPF
domain_error.ErrInvalidCNPJ
domain_error.ErrRequiredField
```

### 🔐 Autenticação (401)
```go
domain_error.ErrUnauthorized
domain_error.ErrInvalidCredentials
domain_error.ErrTokenInvalid
domain_error.ErrTokenExpired
domain_error.ErrSessionExpired
```

### 🚫 Autorização (403)
```go
domain_error.ErrForbidden
domain_error.ErrInsufficientPermissions
domain_error.ErrAccountSuspended
domain_error.ErrCompanySuspended
domain_error.ErrModuleNotContracted
```

### 🔍 Não Encontrado (404)
```go
domain_error.ErrNotFound
domain_error.ErrFileNotFound
```

### ⚠️ Conflito (409)
```go
domain_error.ErrConflict
domain_error.ErrDuplicateEmail
domain_error.ErrRecordLocked
domain_error.ErrConcurrentModification
domain_error.ErrIdempotencyConflict
```

### 💰 Regras de Negócio (422)
```go
domain_error.ErrInsufficientBalance
domain_error.ErrPaymentOverdue
domain_error.ErrLeadAlreadyConverted
domain_error.ErrContractExpired
domain_error.ErrInvalidStatus
```

### ⏱️ Rate Limiting (429)
```go
domain_error.ErrRateLimitExceeded
domain_error.ErrQuotaExceeded
domain_error.ErrMaxAttemptsExceeded
```

### 💥 Erros de Sistema (500+)
```go
domain_error.ErrInternalServer
domain_error.ErrDatabaseQuery
domain_error.ErrExternalServiceUnavailable
domain_error.ErrExternalServiceTimeout
```

## 🔧 Exemplo AWS Lambda Handler
```go
package main

import (
    "context"
    "encoding/json"
    "github.com/aws/aws-lambda-go/events"
    "github.com/renatofagalde/module-error"
)

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    mapper := domain_error.NewHTTPStatusMapper()
    
    // Valida empresa suspensa
    if isCompanySuspended(ctx) {
        err := domain_error.ErrCompanySuspended
        return errorResponse(mapper, err), nil
    }
    
    // Valida módulo contratado
    if !hasModuleAccess(ctx, "CRM") {
        err := domain_error.ErrModuleNotContracted
        return errorResponse(mapper, err), nil
    }
    
    // Lógica de negócio...
    user, err := createUser(req)
    if err != nil {
        return errorResponse(mapper, err), nil
    }
    
    return successResponse(201, user), nil
}

func errorResponse(mapper *domain_error.HTTPStatusMapper, err error) events.APIGatewayProxyResponse {
    statusCode := mapper.GetHTTPStatus(err)
    body, _ := json.Marshal(err)
    
    return events.APIGatewayProxyResponse{
        StatusCode: statusCode,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: string(body),
    }
}
```

## 🧪 Testes
```bash
# Executar testes
go test -v

# Cobertura de testes
go test -cover

# Relatório de cobertura
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📊 Mapeamento HTTP Status Codes

| Status | Categoria | Exemplo |
|--------|-----------|---------|
| 400 | Bad Request | `ErrInvalidInput` |
| 401 | Unauthorized | `ErrTokenExpired` |
| 403 | Forbidden | `ErrCompanySuspended` |
| 404 | Not Found | `ErrNotFound` |
| 409 | Conflict | `ErrDuplicateEmail` |
| 422 | Unprocessable | `ErrInsufficientBalance` |
| 429 | Too Many Requests | `ErrRateLimitExceeded` |
| 500 | Internal Error | `ErrInternalServer` |
| 502 | Bad Gateway | `ErrThirdPartyAPIError` |
| 503 | Unavailable | `ErrServiceUnavailable` |

## 🏗️ Arquitetura
```
module-error/
├── domain_error.go          # Definições de erros
├── http_mapper.go           # Mapeamento HTTP
├── domain_error_test.go     # Testes unitários
├── go.mod                   # Módulo Go
├── .gitignore              # Git ignore
└── readme.md               # Documentação
```

## 🤝 Contribuindo

Para adicionar novos erros:

1. Adicione a variável de erro em `domain_error.go`
2. Adicione o mapeamento HTTP em `http_mapper.go`
3. Adicione testes em `domain_error_test.go`
4. Execute `go test -v` para validar

## 📝 Licença

Este módulo faz parte do ecossistema de microserviços da API REST multi-tenant.

## 👤 Autor

**Renato Fagalde**
- GitHub: [@renatofagalde](https://github.com/renatofagalde)

---

⭐️ **Padrão REST Purista** | **Multi-tenant Ready** | **Production Ready**
