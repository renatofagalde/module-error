// Package validation instala os guardrails de input do projeto sobre o
// validator do gin: validadores customizados, nome de campo em JSON e
// traducao das falhas para o envelope Problem Details.
package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	registerOnce sync.Once

	// money4: ate 11 digitos inteiros e no maximo 4 casas decimais.
	// Rejeita "1e9", "1.23456789", "NaN", "+1", " 10", "10,50".
	reMoney4 = regexp.MustCompile(`^-?\d{1,11}(\.\d{1,4})?$`)

	reNonDigit = regexp.MustCompile(`\D`)
)

// Register e idempotente e deve ser chamado uma vez no bootstrap
// (cmd/api/main.go e cmd/lambda/main.go), antes de servir requisicoes.
func Register() {
	registerOnce.Do(func() {
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			return
		}

		// O erro devolve o nome do campo em JSON (premium_value), nunca o nome
		// Go (PremiumValue). Evita vazar a estrutura interna e entrega ao
		// frontend uma chave que ele consegue casar com o formulario.
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			for _, tag := range []string{"json", "uri", "form"} {
				name := strings.SplitN(fld.Tag.Get(tag), ",", 2)[0]
				if name == "-" {
					return ""
				}
				if name != "" {
					return name
				}
			}
			return fld.Name
		})

		_ = v.RegisterValidation("money4", isMoney4)
		_ = v.RegisterValidation("money4pos", isMoney4Positive)
		_ = v.RegisterValidation("percent", isPercent)
		_ = v.RegisterValidation("cpfcnpj", isCPFCNPJ)
		_ = v.RegisterValidation("phonebr", isPhoneBR)
	})
}

// fieldString extrai o valor textual do campo. Cobre string e qualquer tipo
// que implemente fmt.Stringer (decimal.Decimal, por exemplo), sem que este
// pacote precise depender de shopspring/decimal.
func fieldString(fl validator.FieldLevel) (string, bool) {
	f := fl.Field()
	if f.Kind() == reflect.String {
		return f.String(), true
	}
	if s, ok := f.Interface().(fmt.Stringer); ok {
		return s.String(), true
	}
	return "", false
}

func isMoney4(fl validator.FieldLevel) bool {
	s, ok := fieldString(fl)
	if !ok {
		return false
	}
	return reMoney4.MatchString(s)
}

func isMoney4Positive(fl validator.FieldLevel) bool {
	s, ok := fieldString(fl)
	if !ok || !reMoney4.MatchString(s) {
		return false
	}
	return !strings.HasPrefix(s, "-")
}

// isPercent aceita 0..100, com no maximo 4 casas. Serve para string,
// decimal.Decimal e tipos numericos.
func isPercent(fl validator.FieldLevel) bool {
	f := fl.Field()

	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := f.Int()
		return n >= 0 && n <= 100
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return f.Uint() <= 100
	case reflect.Float32, reflect.Float64:
		n := f.Float()
		return n >= 0 && n <= 100
	}

	s, ok := fieldString(fl)
	if !ok || !reMoney4.MatchString(s) {
		return false
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false
	}
	return n >= 0 && n <= 100
}

func isCPFCNPJ(fl validator.FieldLevel) bool {
	s, ok := fieldString(fl)
	if !ok {
		return false
	}
	digits := reNonDigit.ReplaceAllString(s, "")
	switch len(digits) {
	case 11:
		return validCPF(digits)
	case 14:
		return validCNPJ(digits)
	}
	return false
}

func isPhoneBR(fl validator.FieldLevel) bool {
	s, ok := fieldString(fl)
	if !ok {
		return false
	}
	n := len(reNonDigit.ReplaceAllString(s, ""))
	return n == 10 || n == 11
}

func allSameDigit(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

// validCPF confere os dois digitos verificadores. Validar so o tamanho deixa
// passar 111.111.111-11, que e a primeira coisa que um auditor testa.
func validCPF(d string) bool {
	if allSameDigit(d) {
		return false
	}
	for _, dv := range []int{9, 10} {
		sum := 0
		weight := dv + 1
		for i := 0; i < dv; i++ {
			sum += int(d[i]-'0') * weight
			weight--
		}
		check := (sum * 10) % 11
		if check == 10 {
			check = 0
		}
		if check != int(d[dv]-'0') {
			return false
		}
	}
	return true
}

func validCNPJ(d string) bool {
	if allSameDigit(d) {
		return false
	}
	weights := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	for _, dv := range []int{12, 13} {
		w := weights[len(weights)-dv:]
		sum := 0
		for i := 0; i < dv; i++ {
			sum += int(d[i]-'0') * w[i]
		}
		check := sum % 11
		if check < 2 {
			check = 0
		} else {
			check = 11 - check
		}
		if check != int(d[dv]-'0') {
			return false
		}
	}
	return true
}
