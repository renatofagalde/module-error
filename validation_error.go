package domainerror

import "fmt"

// ErrValidationFailed e o erro de topo de qualquer falha de binding/validacao.
var ErrValidationFailed = New("VALIDATION_FAILED", "Um ou mais campos são inválidos")

// FieldError descreve a falha de UM campo do input. Vai no extension member
// "errors" do Problem, permitindo devolver TODAS as falhas de uma vez.
//
// Code e estavel e machine-readable (contrato com o frontend, base de i18n).
// Message e o fallback humano em pt-BR.
// Constraint carrega o limite violado (ex.: "360", "4", "3").
//
// O VALOR RECEBIDO NUNCA APARECE AQUI. Documento, e-mail, telefone e senha
// ecoados em resposta de erro vazam para log de proxy, APM e print de tela.
type FieldError struct {
	Field      string `json:"field"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Constraint string `json:"constraint,omitempty"`
}

// ValidationError agrega as falhas de campo de uma unica requisicao.
type ValidationError struct {
	Errors []FieldError
}

// NewValidationError constroi o erro de validacao a partir das falhas de campo.
func NewValidationError(fields ...FieldError) *ValidationError {
	return &ValidationError{Errors: fields}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s (%d campo(s))",
		ErrValidationFailed.Code, ErrValidationFailed.Message, len(e.Errors))
}

// Unwrap mantem errors.Is/errors.As enxergando um *DomainError, para que o
// mapeamento de status HTTP existente continue funcionando sem alteracao.
func (e *ValidationError) Unwrap() error {
	return ErrValidationFailed
}
