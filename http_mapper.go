package domainerror

import "net/http"

// HTTPStatusMapper mapeia DomainError para HTTP status codes
type HTTPStatusMapper struct {
	errorToStatus map[string]int
}

// NewHTTPStatusMapper cria um novo mapper de erros para status HTTP
func NewHTTPStatusMapper() *HTTPStatusMapper {
	mapper := &HTTPStatusMapper{
		errorToStatus: make(map[string]int),
	}
	mapper.initialize()
	return mapper
}

// initialize configura o mapeamento de códigos de erro para status HTTP
func (m *HTTPStatusMapper) initialize() {
	// 400 Bad Request
	m.errorToStatus[ErrInvalidInput.Code] = http.StatusBadRequest
	m.errorToStatus[ErrInvalidEmail.Code] = http.StatusBadRequest
	m.errorToStatus[ErrInvalidCPF.Code] = http.StatusBadRequest
	m.errorToStatus[ErrInvalidCNPJ.Code] = http.StatusBadRequest
	m.errorToStatus[ErrInvalidPhone.Code] = http.StatusBadRequest
	m.errorToStatus[ErrInvalidDate.Code] = http.StatusBadRequest
	m.errorToStatus[ErrInvalidCurrency.Code] = http.StatusBadRequest
	m.errorToStatus[ErrRequiredField.Code] = http.StatusBadRequest
	m.errorToStatus[ErrInvalidFileType.Code] = http.StatusBadRequest

	// 401 Unauthorized
	m.errorToStatus[ErrUnauthorized.Code] = http.StatusUnauthorized
	m.errorToStatus[ErrInvalidCredentials.Code] = http.StatusUnauthorized
	m.errorToStatus[ErrTokenInvalid.Code] = http.StatusUnauthorized
	m.errorToStatus[ErrTokenExpired.Code] = http.StatusUnauthorized
	m.errorToStatus[ErrSessionExpired.Code] = http.StatusUnauthorized

	// 403 Forbidden
	m.errorToStatus[ErrForbidden.Code] = http.StatusForbidden
	m.errorToStatus[ErrInsufficientPermissions.Code] = http.StatusForbidden
	m.errorToStatus[ErrAccountSuspended.Code] = http.StatusForbidden
	m.errorToStatus[ErrAccountInactive.Code] = http.StatusForbidden
	m.errorToStatus[ErrCompanySuspended.Code] = http.StatusForbidden
	m.errorToStatus[ErrModuleNotContracted.Code] = http.StatusForbidden

	// 404 Not Found
	m.errorToStatus[ErrNotFound.Code] = http.StatusNotFound
	m.errorToStatus[ErrFileNotFound.Code] = http.StatusNotFound

	// 405 Method Not Allowed
	m.errorToStatus[ErrMethodNotAllowed.Code] = http.StatusMethodNotAllowed

	// 406 Not Acceptable
	m.errorToStatus[ErrNotAcceptable.Code] = http.StatusNotAcceptable

	// 408 Request Timeout
	m.errorToStatus[ErrRequestTimeout.Code] = http.StatusRequestTimeout

	// 409 Conflict
	m.errorToStatus[ErrConflict.Code] = http.StatusConflict
	m.errorToStatus[ErrDuplicateEmail.Code] = http.StatusConflict
	m.errorToStatus[ErrDuplicateCPF.Code] = http.StatusConflict
	m.errorToStatus[ErrDuplicateCNPJ.Code] = http.StatusConflict
	m.errorToStatus[ErrRecordLocked.Code] = http.StatusConflict
	m.errorToStatus[ErrStatusConflict.Code] = http.StatusConflict
	m.errorToStatus[ErrIdempotencyConflict.Code] = http.StatusConflict
	m.errorToStatus[ErrConcurrentModification.Code] = http.StatusConflict
	m.errorToStatus[ErrCircularReference.Code] = http.StatusConflict
	m.errorToStatus[ErrDuplicateRequest.Code] = http.StatusConflict
	m.errorToStatus[ErrIdempotencyKeyUsed.Code] = http.StatusConflict
	m.errorToStatus[ErrDuplicateLead.Code] = http.StatusConflict

	// 410 Gone
	m.errorToStatus[ErrResourceGone.Code] = http.StatusGone
	m.errorToStatus[ErrResourceArchived.Code] = http.StatusGone

	// 412 Precondition Failed
	m.errorToStatus[ErrPreconditionFailed.Code] = http.StatusPreconditionFailed
	m.errorToStatus[ErrETagMismatch.Code] = http.StatusPreconditionFailed
	m.errorToStatus[ErrOptimisticLockFailed.Code] = http.StatusPreconditionFailed

	// 413 Payload Too Large
	m.errorToStatus[ErrFileTooLarge.Code] = http.StatusRequestEntityTooLarge

	// 415 Unsupported Media Type
	m.errorToStatus[ErrUnsupportedMediaType.Code] = http.StatusUnsupportedMediaType

	// 417 Expectation Failed
	m.errorToStatus[ErrExpectationFailed.Code] = http.StatusExpectationFailed

	// 422 Unprocessable Entity
	m.errorToStatus[ErrInvalidStatus.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrLeadAlreadyConverted.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrInvalidLeadStatus.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrInvalidRelationship.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrInsufficientBalance.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrPaymentOverdue.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrPaymentFailed.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrInvoiceNotPaid.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrCreditLimitExceeded.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrCustomerNotActive.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrContractExpired.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrContractNotActive.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrOrphanRecord.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrDependencyExists.Code] = http.StatusUnprocessableEntity
	m.errorToStatus[ErrRecordInUse.Code] = http.StatusUnprocessableEntity

	// 423 Locked
	m.errorToStatus[ErrRecordLocked.Code] = http.StatusLocked

	// 424 Failed Dependency
	m.errorToStatus[ErrFailedDependency.Code] = http.StatusFailedDependency

	// 429 Too Many Requests
	m.errorToStatus[ErrRateLimitExceeded.Code] = http.StatusTooManyRequests
	m.errorToStatus[ErrQuotaExceeded.Code] = http.StatusTooManyRequests
	m.errorToStatus[ErrMaxAttemptsExceeded.Code] = http.StatusTooManyRequests

	// 451 Unavailable For Legal Reasons
	m.errorToStatus[ErrUnavailableForLegalReasons.Code] = http.StatusUnavailableForLegalReasons

	// 500 Internal Server Error
	m.errorToStatus[ErrInternalServer.Code] = http.StatusInternalServerError
	m.errorToStatus[ErrDatabaseQuery.Code] = http.StatusInternalServerError
	m.errorToStatus[ErrFileUploadFailed.Code] = http.StatusInternalServerError

	// 502 Bad Gateway
	m.errorToStatus[ErrThirdPartyAPIError.Code] = http.StatusBadGateway
	m.errorToStatus[ErrExternalServiceUnavailable.Code] = http.StatusBadGateway

	// 503 Service Unavailable
	m.errorToStatus[ErrServiceUnavailable.Code] = http.StatusServiceUnavailable
	m.errorToStatus[ErrDatabaseConnection.Code] = http.StatusServiceUnavailable

	// 504 Gateway Timeout
	m.errorToStatus[ErrExternalServiceTimeout.Code] = http.StatusGatewayTimeout
}

// GetHTTPStatus retorna o status HTTP correspondente ao erro de domínio
func (m *HTTPStatusMapper) GetHTTPStatus(err error) int {
	if domainErr, ok := err.(*DomainError); ok {
		if status, exists := m.errorToStatus[domainErr.Code]; exists {
			return status
		}
	}
	// Default para erro genérico
	return http.StatusInternalServerError
}

// GetHTTPStatusByCode retorna o status HTTP correspondente ao código de erro
func (m *HTTPStatusMapper) GetHTTPStatusByCode(code string) int {
	if status, exists := m.errorToStatus[code]; exists {
		return status
	}
	return http.StatusInternalServerError
}
