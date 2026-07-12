package validation

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// messageFor traduz uma falha do validator em (code, message, constraint).
//
// code    -> contrato estavel, minusculo, base de i18n no frontend
// message -> fallback humano em pt-BR
// unidade -> "caracteres" para string, "itens" para colecao, valor para numero
func messageFor(fe validator.FieldError) (code, message, constraint string) {
	param := fe.Param()
	unit := unitFor(fe.Kind())

	switch fe.Tag() {
	case "required":
		return "required", "Campo obrigatório", ""

	case "min":
		if unit == "" {
			return "min", fmt.Sprintf("Deve ser no mínimo %s", param), param
		}
		return "min", fmt.Sprintf("Deve ter no mínimo %s %s", param, unit), param

	case "max":
		if unit == "" {
			return "max", fmt.Sprintf("Deve ser no máximo %s", param), param
		}
		return "max", fmt.Sprintf("Deve ter no máximo %s %s", param, unit), param

	case "len":
		if unit == "" {
			return "len", fmt.Sprintf("Deve ser exatamente %s", param), param
		}
		return "len", fmt.Sprintf("Deve ter exatamente %s %s", param, unit), param

	case "gte":
		return "gte", fmt.Sprintf("Deve ser maior ou igual a %s", param), param
	case "lte":
		return "lte", fmt.Sprintf("Deve ser menor ou igual a %s", param), param
	case "gt":
		return "gt", fmt.Sprintf("Deve ser maior que %s", param), param
	case "lt":
		return "lt", fmt.Sprintf("Deve ser menor que %s", param), param

	case "oneof":
		return "one_of", fmt.Sprintf("Valor não permitido. Valores aceitos: %s", param), param

	case "email":
		return "invalid_email", "E-mail inválido", ""

	case "uuid", "uuid4", "uuid7":
		return "invalid_uuid", "Identificador inválido", ""

	case "url", "fqdn", "hostname":
		return "invalid_url", "Endereço inválido", ""

	case "numeric":
		return "not_numeric", "Deve conter apenas números", ""

	case "money4":
		return "money_scale", "Valor monetário deve ter no máximo 4 casas decimais", "4"

	case "money4pos":
		return "money_scale_positive", "Valor monetário deve ser positivo e ter no máximo 4 casas decimais", "4"

	case "percent":
		return "invalid_percent", "Percentual deve estar entre 0 e 100", "100"

	case "percent2":
		return "invalid_percent", "Percentual deve estar entre 0 e 100, com no máximo 2 casas decimais", "100"

	case "percent2gt":
		return "invalid_percent", "Percentual deve ser maior que 0 e no máximo 100, com até 2 casas decimais", "100"

	case "money4gt":
		return "money_scale_positive", "Valor monetário deve ser maior que zero e ter no máximo 4 casas decimais", "4"

	case "cpfcnpj":
		return "invalid_document", "CPF ou CNPJ inválido", ""

	case "phonebr":
		return "invalid_phone", "Telefone inválido", ""

	case "eqfield", "nefield", "eqcsfield":
		return "field_mismatch", fmt.Sprintf("Não confere com o campo %s", param), param
	}

	// Fallback: nunca devolvemos a mensagem crua do validator, que expoe a
	// struct interna (ex.: "Key: 'CreateProposalDTO.PremiumValue' ...").
	return "invalid_field", "Valor inválido", param
}

func unitFor(k reflect.Kind) string {
	switch k {
	case reflect.String:
		return "caracteres"
	case reflect.Slice, reflect.Array, reflect.Map:
		return "itens"
	}
	return ""
}
