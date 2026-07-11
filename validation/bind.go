package validation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	domainerror "github.com/renatofagalde/module-error"
)

// BindJSON substitui c.ShouldBindJSON nos handlers. Em caso de falha devolve
// sempre um erro do dominio, pronto para httperror.WriteError.
//
//	if err := validation.BindJSON(c, &dto); err != nil {
//	    httperror.WriteError(c, err)
//	    return
//	}
func BindJSON(c *gin.Context, obj any) error {
	return Translate(c.ShouldBindJSON(obj))
}

// BindURI valida parametros de rota (tags uri).
func BindURI(c *gin.Context, obj any) error {
	return Translate(c.ShouldBindUri(obj))
}

// BindQuery valida query string (tags form).
func BindQuery(c *gin.Context, obj any) error {
	return Translate(c.ShouldBindQuery(obj))
}

// Translate converte o erro cru do gin/validator no erro do dominio.
// Todas as falhas de campo voltam de uma vez, nunca uma por requisicao.
func Translate(err error) error {
	if err == nil {
		return nil
	}

	// Falhas de regra: uma entrada por campo.
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make([]domainerror.FieldError, 0, len(ve))
		for _, fe := range ve {
			code, message, constraint := messageFor(fe)
			fields = append(fields, domainerror.FieldError{
				Field:      fieldPath(fe),
				Code:       code,
				Message:    message,
				Constraint: constraint,
			})
		}
		return domainerror.NewValidationError(fields...)
	}

	// Tipo errado no JSON: "installments": "abc". Nao chega no validator.
	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		field := ute.Field
		if field == "" {
			field = "body"
		}
		return domainerror.NewValidationError(domainerror.FieldError{
			Field:      field,
			Code:       "invalid_type",
			Message:    "Tipo de dado inválido para este campo",
			Constraint: ute.Type.String(),
		})
	}

	// JSON malformado.
	var se *json.SyntaxError
	if errors.As(err, &se) {
		return domainerror.NewValidationError(domainerror.FieldError{
			Field:   "body",
			Code:    "malformed_json",
			Message: "Corpo da requisição não é um JSON válido",
		})
	}

	// Body vazio.
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return domainerror.NewValidationError(domainerror.FieldError{
			Field:   "body",
			Code:    "required",
			Message: "Corpo da requisição é obrigatório",
		})
	}

	// Body acima do limite (http.MaxBytesReader no middleware).
	var mbe *http.MaxBytesError
	if errors.As(err, &mbe) {
		return domainerror.ErrFileTooLarge
	}

	// Struct mal configurada: bug de programacao, nao input do cliente.
	var ive *validator.InvalidValidationError
	if errors.As(err, &ive) {
		return domainerror.ErrInternalServer
	}

	return domainerror.ErrInvalidInput
}

// fieldPath devolve o caminho do campo em notacao JSON, incluindo indice de
// colecao quando ha "dive": items[0].percent
func fieldPath(fe validator.FieldError) string {
	ns := fe.Namespace()
	if i := strings.Index(ns, "."); i >= 0 {
		ns = ns[i+1:] // remove o nome da struct raiz
	}
	if ns == "" {
		return fe.Field()
	}
	return ns
}
