package validation_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	domainerror "github.com/renatofagalde/module-error"
	"github.com/renatofagalde/module-error/httperror"
	"github.com/renatofagalde/module-error/validation"
)

type proposalDTO struct {
	ProductName    string `json:"product_name"    binding:"required,min=3,max=120"`
	Installments   int    `json:"installments"    binding:"required,gte=1,lte=360"`
	PremiumValue   string `json:"premium_value"   binding:"required,money4pos"`
	PersonDocument string `json:"person_document" binding:"required,cpfcnpj"`
	Percent        string `json:"percent"         binding:"required,percent"`
}

const validBody = `{"product_name":"Plano Jedi Premium","installments":12,` +
	`"premium_value":"1500.0000","person_document":"111.444.777-35","percent":"15"}`

func bind(t *testing.T, body string) error {
	t.Helper()
	gin.SetMode(gin.TestMode)
	validation.Register()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/cms/proposals", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	var dto proposalDTO
	return validation.BindJSON(c, &dto)
}

func fieldsOf(t *testing.T, err error) map[string]domainerror.FieldError {
	t.Helper()
	var verr *domainerror.ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("esperado *ValidationError, recebido %T (%v)", err, err)
	}
	out := make(map[string]domainerror.FieldError, len(verr.Errors))
	for _, fe := range verr.Errors {
		out[fe.Field] = fe
	}
	return out
}

func TestBindJSON_Valid(t *testing.T) {
	if err := bind(t, validBody); err != nil {
		t.Fatalf("esperado sucesso, recebido: %v", err)
	}
}

func TestBindJSON_FieldNameIsJSON(t *testing.T) {
	err := bind(t, `{"installments":12,"premium_value":"10.00","person_document":"111.444.777-35","percent":"1"}`)
	f := fieldsOf(t, err)
	if _, ok := f["product_name"]; !ok {
		t.Fatalf("esperado campo 'product_name' em snake_case, recebido: %+v", f)
	}
	if _, ok := f["ProductName"]; ok {
		t.Fatal("nome interno da struct vazou na resposta de erro")
	}
}

func TestBindJSON_ReturnsAllFailuresAtOnce(t *testing.T) {
	err := bind(t, `{"product_name":"a","installments":0,"premium_value":"1e9","person_document":"1","percent":"150"}`)
	f := fieldsOf(t, err)
	for _, want := range []string{"product_name", "installments", "premium_value", "person_document", "percent"} {
		if _, ok := f[want]; !ok {
			t.Errorf("campo ausente na lista de erros: %s", want)
		}
	}
}

func TestBindJSON_MoneyScale(t *testing.T) {
	cases := map[string]bool{
		"1500.00":      true,
		"1500.1234":    true,
		"0":            true,
		"1500.12345":   false, // 5 casas
		"1e9":          false, // notacao cientifica
		"-10.00":       false, // negativo
		"1500,00":      false, // virgula
		"999999999999": false, // 12 digitos inteiros
	}
	for value, wantOK := range cases {
		err := bind(t, `{"product_name":"Plano","installments":1,"premium_value":"`+value+
			`","person_document":"111.444.777-35","percent":"1"}`)
		if wantOK {
			if err != nil {
				t.Errorf("premium_value=%q deveria passar, falhou: %v", value, err)
			}
			continue
		}
		if err == nil {
			t.Errorf("premium_value=%q deveria falhar, passou", value)
			continue
		}
		if fe := fieldsOf(t, err)["premium_value"]; fe.Code != "money_scale_positive" {
			t.Errorf("premium_value=%q: code inesperado %+v", value, fe)
		}
	}
}

func TestBindJSON_InstallmentsRange(t *testing.T) {
	err := bind(t, `{"product_name":"Plano","installments":999999,"premium_value":"10.00","person_document":"111.444.777-35","percent":"1"}`)
	if fe := fieldsOf(t, err)["installments"]; fe.Code != "lte" || fe.Constraint != "360" {
		t.Fatalf("esperado lte/360, recebido %+v", fe)
	}
}

func TestBindJSON_InvalidDocument(t *testing.T) {
	// 11 digitos, tamanho correto, digito verificador errado.
	err := bind(t, `{"product_name":"Plano","installments":1,"premium_value":"10.00","person_document":"111.111.111-11","percent":"1"}`)
	if fe := fieldsOf(t, err)["person_document"]; fe.Code != "invalid_document" {
		t.Fatalf("esperado invalid_document, recebido %+v", fe)
	}
}

func TestBindJSON_PercentCap(t *testing.T) {
	err := bind(t, `{"product_name":"Plano","installments":1,"premium_value":"10.00","person_document":"111.444.777-35","percent":"150"}`)
	if fe := fieldsOf(t, err)["percent"]; fe.Code != "invalid_percent" {
		t.Fatalf("esperado invalid_percent, recebido %+v", fe)
	}
}

func TestBindJSON_InvalidType(t *testing.T) {
	err := bind(t, `{"product_name":"Plano","installments":"doze","premium_value":"10.00","person_document":"111.444.777-35","percent":"1"}`)
	if fe := fieldsOf(t, err)["installments"]; fe.Code != "invalid_type" {
		t.Fatalf("esperado invalid_type, recebido %+v", fe)
	}
}

func TestBindJSON_EmptyBody(t *testing.T) {
	if fe := fieldsOf(t, bind(t, ``))["body"]; fe.Code != "required" {
		t.Fatalf("esperado body/required, recebido %+v", fe)
	}
}

func TestWriteError_ProblemDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validation.Register()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/cms/proposals",
		bytes.NewBufferString(`{"product_name":"a","installments":0,"premium_value":"x","person_document":"1","percent":"9"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	var dto proposalDTO
	httperror.WriteError(c, validation.BindJSON(c, &dto))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperado 400, recebido %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != httperror.ProblemContentType {
		t.Fatalf("esperado %q, recebido %q", httperror.ProblemContentType, ct)
	}

	var p domainerror.Problem
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		t.Fatalf("corpo nao e JSON: %v", err)
	}
	if p.Type != "https://api.gestao.one/errors/validation-failed" {
		t.Errorf("type inesperado: %s", p.Type)
	}
	if p.Status != http.StatusBadRequest || p.Code != domainerror.ErrValidationFailed.Code {
		t.Errorf("status/code inesperados: %d / %s", p.Status, p.Code)
	}
	if p.Instance != "/cms/proposals" {
		t.Errorf("instance inesperado: %s", p.Instance)
	}
	if len(p.Errors) < 4 {
		t.Errorf("esperado varios erros de campo, recebido %d", len(p.Errors))
	}
}

func TestWriteError_UnknownErrorDoesNotLeak(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/cms/apolices/1", nil)

	httperror.WriteError(c, errors.New("pq: relation \"apolice\" does not exist"))

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("esperado 500, recebido %d", rec.Code)
	}
	if bytes.Contains(rec.Body.Bytes(), []byte("relation")) {
		t.Fatal("mensagem interna vazou para o cliente")
	}
}

func TestTypeURI(t *testing.T) {
	got := domainerror.TypeURI(domainerror.ErrInsufficientBalance.Code)
	if got != "https://api.gestao.one/errors/insufficient-balance" {
		t.Fatalf("type URI inesperado: %s", got)
	}
}

func TestValidationError_UnwrapsToDomainError(t *testing.T) {
	err := error(domainerror.NewValidationError())
	var derr *domainerror.DomainError
	if !errors.As(err, &derr) {
		t.Fatal("ValidationError deveria fazer Unwrap para *DomainError")
	}
	if derr.Code != domainerror.ErrValidationFailed.Code {
		t.Fatalf("code inesperado: %s", derr.Code)
	}
}
