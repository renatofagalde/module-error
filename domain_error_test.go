package domainerror

import (
	"net/http"
	"testing"
)

func TestDomainError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *DomainError
		expected string
	}{
		{
			name:     "Invalid input error",
			err:      ErrInvalidInput,
			expected: "INVALID_INPUT: Input inválido",
		},
		{
			name:     "Not found error",
			err:      ErrNotFound,
			expected: "NOT_FOUND: Registro não encontrado",
		},
		{
			name:     "Unauthorized error",
			err:      ErrUnauthorized,
			expected: "UNAUTHORIZED: Não autorizado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	code := "TEST_ERROR"
	message := "Mensagem de teste"

	err := New(code, message)

	if err.Code != code {
		t.Errorf("Code = %v, want %v", err.Code, code)
	}

	if err.Message != message {
		t.Errorf("Message = %v, want %v", err.Message, message)
	}
}

func TestHTTPStatusMapper_GetHTTPStatus(t *testing.T) {
	mapper := NewHTTPStatusMapper()

	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "Invalid input returns 400",
			err:            ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthorized returns 401",
			err:            ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Forbidden returns 403",
			err:            ErrForbidden,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Not found returns 404",
			err:            ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Conflict returns 409",
			err:            ErrConflict,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Rate limit returns 429",
			err:            ErrRateLimitExceeded,
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "Internal server error returns 500",
			err:            ErrInternalServer,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := mapper.GetHTTPStatus(tt.err)
			if status != tt.expectedStatus {
				t.Errorf("GetHTTPStatus() = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestHTTPStatusMapper_GetHTTPStatusByCode(t *testing.T) {
	mapper := NewHTTPStatusMapper()

	tests := []struct {
		name           string
		code           string
		expectedStatus int
	}{
		{
			name:           "INVALID_INPUT returns 400",
			code:           "INVALID_INPUT",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "NOT_FOUND returns 404",
			code:           "NOT_FOUND",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unknown code returns 500",
			code:           "UNKNOWN_ERROR",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := mapper.GetHTTPStatusByCode(tt.code)
			if status != tt.expectedStatus {
				t.Errorf("GetHTTPStatusByCode() = %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}
