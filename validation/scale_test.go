package validation_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	domainerror "github.com/renatofagalde/module-error"
	"github.com/renatofagalde/module-error/validation"
	"github.com/shopspring/decimal"
)

// Espelha os DTOs reais do app-cms, onde os campos sao decimal.Decimal e as
// colunas sao NUMERIC(5,2) para percentual e NUMERIC(15,4) para dinheiro.
type planLineDTO struct {
	RecipientType string          `json:"recipient_type" binding:"required,oneof=seller corretora team_role"`
	Percent       decimal.Decimal `json:"percent"        binding:"required,percent2gt"`
	TaxPercent    decimal.Decimal `json:"tax_percent"    binding:"percent2"`
	PremiumValue  decimal.Decimal `json:"premium_value"  binding:"required,money4gt"`
}

func bindLine(t *testing.T, body string) error {
	t.Helper()
	gin.SetMode(gin.TestMode)
	validation.Register()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/cms/plans", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	var dto planLineDTO
	return validation.BindJSON(c, &dto)
}

func codeFor(t *testing.T, err error, field string) string {
	t.Helper()
	var verr *domainerror.ValidationError
	if !errors.As(err, &verr) {
		t.Fatalf("esperado *ValidationError, recebido %T (%v)", err, err)
	}
	for _, fe := range verr.Errors {
		if fe.Field == field {
			return fe.Code
		}
	}
	return ""
}

func body(percent, tax, premium string) string {
	return `{"recipient_type":"seller","percent":` + percent +
		`,"tax_percent":` + tax + `,"premium_value":` + premium + `}`
}

func TestPercent2gt_Scale(t *testing.T) {
	cases := map[string]bool{
		"15":     true,
		"15.5":   true,
		"15.12":  true,
		"100":    true,
		"0":      false, // CHECK: percent > 0
		"15.123": false, // 3 casas: NUMERIC(5,2) arredondaria em silencio
		"100.01": false, // acima de 100
		"-15":    false,
		"1000":   false,
	}
	for value, wantOK := range cases {
		err := bindLine(t, body(value, "0", "1500.00"))
		if wantOK {
			if err != nil {
				t.Errorf("percent=%s deveria passar: %v", value, err)
			}
			continue
		}
		if err == nil {
			t.Errorf("percent=%s deveria falhar, passou", value)
			continue
		}
		if c := codeFor(t, err, "percent"); c != "invalid_percent" && c != "required" {
			t.Errorf("percent=%s: code inesperado %q", value, c)
		}
	}
}

func TestPercent2_AllowsZero(t *testing.T) {
	if err := bindLine(t, body("10", "0", "1500.00")); err != nil {
		t.Fatalf("tax_percent=0 deveria passar (DEFAULT 0.00): %v", err)
	}
	if err := bindLine(t, body("10", "5.123", "1500.00")); err == nil {
		t.Fatal("tax_percent=5.123 deveria falhar: 3 casas em NUMERIC(5,2)")
	}
}

func TestMoney4gt_RejectsZero(t *testing.T) {
	// CHECK (premium_value > 0): zero nao pode virar violacao de constraint (500).
	err := bindLine(t, body("10", "0", "0"))
	if err == nil {
		t.Fatal("premium_value=0 deveria falhar no binding, nao no Postgres")
	}
	if c := codeFor(t, err, "premium_value"); c != "money_scale_positive" && c != "required" {
		t.Fatalf("code inesperado: %q", c)
	}
}

func TestMoney4gt_Scale(t *testing.T) {
	if err := bindLine(t, body("10", "0", "1500.1234")); err != nil {
		t.Fatalf("4 casas deveria passar: %v", err)
	}
	if err := bindLine(t, body("10", "0", "1500.12345")); err == nil {
		t.Fatal("5 casas deveria falhar")
	}
}

func TestRecipientType_ClosedSet(t *testing.T) {
	// CHECK: recipient_type IN ('seller','corretora','team_role')
	ok := `{"recipient_type":"team_role","percent":25,"tax_percent":0,"premium_value":"1500.00"}`
	if err := bindLine(t, ok); err != nil {
		t.Fatalf("team_role deveria passar: %v", err)
	}
	bad := `{"recipient_type":"vendedor","percent":25,"tax_percent":0,"premium_value":"1500.00"}`
	err := bindLine(t, bad)
	if err == nil {
		t.Fatal("recipient_type invalido deveria falhar")
	}
	if c := codeFor(t, err, "recipient_type"); c != "one_of" {
		t.Fatalf("code inesperado: %q", c)
	}
}
