package domainerror

import "strings"

// ProblemBaseURI e o prefixo dos URIs de tipo de erro. Sobrescreva no bootstrap
// se a documentacao publica mudar de endereco.
var ProblemBaseURI = "https://api.gestao.one/errors/"

// Problem e o envelope unico de erro da API, conforme RFC 7807
// (atualizada pela RFC 9457 - Problem Details for HTTP APIs).
//
// Content-Type: application/problem+json
//
//	{
//	  "type":     "https://api.gestao.one/errors/validation-failed",
//	  "title":    "Um ou mais campos são inválidos",
//	  "status":   400,
//	  "instance": "/cms/proposals",
//	  "code":     "VALIDATION_FAILED",
//	  "trace_id": "1-6890abc...",
//	  "errors":   [ ... ]
//	}
//
// code, trace_id e errors sao extension members (permitidos pela spec).
// code e o contrato estavel com o cliente; errors so aparece em validacao.
type Problem struct {
	Type     string       `json:"type"`
	Title    string       `json:"title"`
	Status   int          `json:"status"`
	Detail   string       `json:"detail,omitempty"`
	Instance string       `json:"instance,omitempty"`
	Code     string       `json:"code"`
	TraceID  string       `json:"trace_id,omitempty"`
	Errors   []FieldError `json:"errors,omitempty"`
}

// TypeURI deriva o URI do tipo a partir do codigo do erro de dominio.
// INSUFFICIENT_BALANCE -> https://api.gestao.one/errors/insufficient-balance
//
// Derivar em vez de manter um catalogo paralelo de URIs: o codigo continua
// sendo a unica fonte de verdade e nao ha como as duas listas divergirem.
func TypeURI(code string) string {
	return ProblemBaseURI + strings.ToLower(strings.ReplaceAll(code, "_", "-"))
}

// NewProblem monta o Problem a partir de um erro do catalogo.
// O title vem do catalogo (pt-BR): uma unica fonte de texto para o erro.
func NewProblem(d *DomainError, status int) *Problem {
	return &Problem{
		Type:   TypeURI(d.Code),
		Title:  d.Message,
		Status: status,
		Code:   d.Code,
	}
}
