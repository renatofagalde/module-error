package domainerror

import "fmt"

type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

// Erros de Validação e Input
var (
	ErrInvalidInput    = New("INVALID_INPUT", "Input inválido")
	ErrInvalidEmail    = New("INVALID_EMAIL", "Email inválido")
	ErrInvalidCPF      = New("INVALID_CPF", "CPF inválido")
	ErrInvalidCNPJ     = New("INVALID_CNPJ", "CNPJ inválido")
	ErrInvalidPhone    = New("INVALID_PHONE", "Telefone inválido")
	ErrInvalidDate     = New("INVALID_DATE", "Data inválida")
	ErrInvalidCurrency = New("INVALID_CURRENCY", "Valor monetário inválido")
	ErrRequiredField   = New("REQUIRED_FIELD", "Campo obrigatório não informado")
)

// Erros de Registro/Recurso
var (
	ErrNotFound       = New("NOT_FOUND", "Registro não encontrado")
	ErrConflict       = New("CONFLICT", "Registro já existente")
	ErrDuplicateEmail = New("DUPLICATE_EMAIL", "Email já cadastrado")
	ErrDuplicateCPF   = New("DUPLICATE_CPF", "CPF já cadastrado")
	ErrDuplicateCNPJ  = New("DUPLICATE_CNPJ", "CNPJ já cadastrado")
	ErrRecordLocked   = New("RECORD_LOCKED", "Registro bloqueado para edição")
	ErrRecordInUse    = New("RECORD_IN_USE", "Registro em uso e não pode ser excluído")
)

// Erros de Autenticação e Autorização
var (
	ErrUnauthorized            = New("UNAUTHORIZED", "Não autorizado")
	ErrForbidden               = New("FORBIDDEN", "Acesso negado")
	ErrInvalidCredentials      = New("INVALID_CREDENTIALS", "Credenciais inválidas")
	ErrSessionExpired          = New("SESSION_EXPIRED", "Sessão expirada")
	ErrTokenInvalid            = New("TOKEN_INVALID", "Token inválido")
	ErrTokenExpired            = New("TOKEN_EXPIRED", "Token expirado")
	ErrInsufficientPermissions = New("INSUFFICIENT_PERMISSIONS", "Permissões insuficientes")
)

// Erros de Negócio - Financeiro
var (
	ErrInsufficientBalance = New("INSUFFICIENT_BALANCE", "Saldo insuficiente")
	ErrPaymentOverdue      = New("PAYMENT_OVERDUE", "Pagamento em atraso")
	ErrPaymentFailed       = New("PAYMENT_FAILED", "Falha no pagamento")
	ErrInvoiceNotPaid      = New("INVOICE_NOT_PAID", "Fatura não paga")
	ErrCreditLimitExceeded = New("CREDIT_LIMIT_EXCEEDED", "Limite de crédito excedido")
)

// Erros de Estado/Status
var (
	ErrInvalidStatus    = New("INVALID_STATUS", "Status inválido para operação")
	ErrStatusConflict   = New("STATUS_CONFLICT", "Conflito de status")
	ErrAccountSuspended = New("ACCOUNT_SUSPENDED", "Conta suspensa")
	ErrAccountInactive  = New("ACCOUNT_INACTIVE", "Conta inativa")
	ErrCompanySuspended = New("COMPANY_SUSPENDED", "Empresa suspensa por inadimplência")
)

// Erros de Idempotência e Concorrência
var (
	ErrDuplicateRequest       = New("DUPLICATE_REQUEST", "Requisição duplicada")
	ErrIdempotencyKeyUsed     = New("IDEMPOTENCY_KEY_USED", "Chave de idempotência já utilizada")
	ErrIdempotencyConflict    = New("IDEMPOTENCY_CONFLICT", "Conflito de idempotência - operação diferente com mesma chave")
	ErrConcurrentModification = New("CONCURRENT_MODIFICATION", "Registro modificado por outro usuário")
	ErrOptimisticLockFailed   = New("OPTIMISTIC_LOCK_FAILED", "Falha no controle de concorrência otimista")
)

// Erros de Limite e Rate Limiting
var (
	ErrRateLimitExceeded   = New("RATE_LIMIT_EXCEEDED", "Limite de requisições excedido")
	ErrQuotaExceeded       = New("QUOTA_EXCEEDED", "Cota excedida")
	ErrMaxAttemptsExceeded = New("MAX_ATTEMPTS_EXCEEDED", "Número máximo de tentativas excedido")
)

// Erros de Integração Externa
var (
	ErrExternalServiceUnavailable = New("EXTERNAL_SERVICE_UNAVAILABLE", "Serviço externo indisponível")
	ErrExternalServiceTimeout     = New("EXTERNAL_SERVICE_TIMEOUT", "Timeout em serviço externo")
	ErrThirdPartyAPIError         = New("THIRD_PARTY_API_ERROR", "Erro em API de terceiros")
)

// Erros de Relacionamento/Dependência
var (
	ErrOrphanRecord        = New("ORPHAN_RECORD", "Registro órfão - relacionamento obrigatório ausente")
	ErrCircularReference   = New("CIRCULAR_REFERENCE", "Referência circular detectada")
	ErrInvalidRelationship = New("INVALID_RELATIONSHIP", "Relacionamento inválido")
	ErrDependencyExists    = New("DEPENDENCY_EXISTS", "Não é possível excluir - existem dependências")
)

// Erros de CRM Específicos
var (
	ErrLeadAlreadyConverted = New("LEAD_ALREADY_CONVERTED", "Lead já convertido em cliente")
	ErrInvalidLeadStatus    = New("INVALID_LEAD_STATUS", "Status do lead não permite esta operação")
	ErrDuplicateLead        = New("DUPLICATE_LEAD", "Lead duplicado")
	ErrCustomerNotActive    = New("CUSTOMER_NOT_ACTIVE", "Cliente não está ativo")
	ErrContractExpired      = New("CONTRACT_EXPIRED", "Contrato expirado")
	ErrContractNotActive    = New("CONTRACT_NOT_ACTIVE", "Contrato não está ativo")
	ErrModuleNotContracted  = New("MODULE_NOT_CONTRACTED", "Módulo não contratado pela empresa")
)

// Erros de Arquivo/Upload
var (
	ErrFileTooLarge     = New("FILE_TOO_LARGE", "Arquivo muito grande")
	ErrInvalidFileType  = New("INVALID_FILE_TYPE", "Tipo de arquivo inválido")
	ErrFileUploadFailed = New("FILE_UPLOAD_FAILED", "Falha no upload do arquivo")
	ErrFileNotFound     = New("FILE_NOT_FOUND", "Arquivo não encontrado")
)

// Erros de Protocolo HTTP
var (
	ErrMethodNotAllowed     = New("METHOD_NOT_ALLOWED", "Método HTTP não permitido")
	ErrNotAcceptable        = New("NOT_ACCEPTABLE", "Formato de resposta não suportado")
	ErrRequestTimeout       = New("REQUEST_TIMEOUT", "Tempo de requisição excedido")
	ErrUnsupportedMediaType = New("UNSUPPORTED_MEDIA_TYPE", "Tipo de mídia não suportado")
	ErrExpectationFailed    = New("EXPECTATION_FAILED", "Expectativa não atendida")
)

// Erros de Precondição e Versionamento
var (
	ErrPreconditionFailed = New("PRECONDITION_FAILED", "Pré-condição falhou")
	ErrETagMismatch       = New("ETAG_MISMATCH", "ETag não corresponde - recurso modificado")
)

// Erros de Remoção e Arquivamento
var (
	ErrResourceGone     = New("RESOURCE_GONE", "Recurso foi permanentemente removido")
	ErrResourceArchived = New("RESOURCE_ARCHIVED", "Recurso foi arquivado")
)

// Erros de Dependência e Compliance
var (
	ErrFailedDependency           = New("FAILED_DEPENDENCY", "Falha em dependência necessária")
	ErrUnavailableForLegalReasons = New("UNAVAILABLE_FOR_LEGAL_REASONS", "Indisponível por razões legais")
)

// Erros de Sistema
var (
	ErrInternalServer     = New("INTERNAL_SERVER_ERROR", "Erro interno do servidor")
	ErrDatabaseConnection = New("DATABASE_CONNECTION_ERROR", "Erro de conexão com banco de dados")
	ErrDatabaseQuery      = New("DATABASE_QUERY_ERROR", "Erro na execução da query")
	ErrServiceUnavailable = New("SERVICE_UNAVAILABLE", "Serviço temporariamente indisponível")
)
