package httperror

import (
	"errors"
	"net/http"

	domainerror "github.com/renatofagalde/module-error"
)

type HTTPStatusMapper interface {
	Status(err error) int
}

type DefaultHTTPStatusMapper struct {
	statusByCode map[string]int
}

var httpErrorMapper = NewDefaultHTTPStatusMapper()

func WriteError(c *gin.Context, err error) {
	status := httpErrorMapper.Status(err)

	if derr, ok := err.(*domainerror.DomainError); ok {
		c.JSON(status, gin.H{
			"code":    derr.Code,
			"message": derr.Message,
		})
		return
	}

	c.JSON(status, gin.H{
		"code":    domainerror.ErrInternalServer.Code,
		"message": "Erro interno do servidor",
	})
}

func NewDefaultHTTPStatusMapper() *DefaultHTTPStatusMapper {
	m := &DefaultHTTPStatusMapper{
		statusByCode: make(map[string]int),
	}

	// ---------------------------------------------------------
	// 400 – Bad Request (erros de validação / input inválido)
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrInvalidInput.Code] = http.StatusBadRequest
	m.statusByCode[domainerror.ErrInvalidEmail.Code] = http.StatusBadRequest
	m.statusByCode[domainerror.ErrInvalidCPF.Code] = http.StatusBadRequest
	m.statusByCode[domainerror.ErrInvalidCNPJ.Code] = http.StatusBadRequest
	m.statusByCode[domainerror.ErrInvalidPhone.Code] = http.StatusBadRequest
	m.statusByCode[domainerror.ErrInvalidDate.Code] = http.StatusBadRequest
	m.statusByCode[domainerror.ErrInvalidCurrency.Code] = http.StatusBadRequest
	m.statusByCode[domainerror.ErrRequiredField.Code] = http.StatusBadRequest

	// ---------------------------------------------------------
	// 401 – Unauthorized / 403 – Forbidden
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrUnauthorized.Code] = http.StatusUnauthorized
	m.statusByCode[domainerror.ErrInvalidCredentials.Code] = http.StatusUnauthorized
	m.statusByCode[domainerror.ErrTokenInvalid.Code] = http.StatusUnauthorized
	m.statusByCode[domainerror.ErrTokenExpired.Code] = http.StatusUnauthorized
	m.statusByCode[domainerror.ErrSessionExpired.Code] = http.StatusUnauthorized

	m.statusByCode[domainerror.ErrForbidden.Code] = http.StatusForbidden
	m.statusByCode[domainerror.ErrInsufficientPermissions.Code] = http.StatusForbidden

	// ---------------------------------------------------------
	// 402 – Payment Required / 422 – Unprocessable Entity
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrInsufficientBalance.Code] = http.StatusPaymentRequired
	m.statusByCode[domainerror.ErrPaymentOverdue.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrPaymentFailed.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrInvoiceNotPaid.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrCreditLimitExceeded.Code] = http.StatusPaymentRequired

	// ---------------------------------------------------------
	// 404 – Not Found
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrNotFound.Code] = http.StatusNotFound
	m.statusByCode[domainerror.ErrFileNotFound.Code] = http.StatusNotFound
	m.statusByCode[domainerror.ErrCustomerNotActive.Code] = http.StatusNotFound

	// ---------------------------------------------------------
	// 409 – Conflict (duplicidade, estado inválido, relacionamento em uso)
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrConflict.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrDuplicateEmail.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrDuplicateCPF.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrDuplicateCNPJ.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrDuplicateLead.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrStatusConflict.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrRecordLocked.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrRecordInUse.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrDuplicateRequest.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrIdempotencyKeyUsed.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrIdempotencyConflict.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrOptimisticLockFailed.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrDependencyExists.Code] = http.StatusConflict
	m.statusByCode[domainerror.ErrLeadAlreadyConverted.Code] = http.StatusConflict

	// ---------------------------------------------------------
	// 410 – Gone / 422 – Unprocessable Entity (negócio/estado)
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrResourceGone.Code] = http.StatusGone
	m.statusByCode[domainerror.ErrResourceArchived.Code] = http.StatusGone

	m.statusByCode[domainerror.ErrInvalidStatus.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrOrphanRecord.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrCircularReference.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrInvalidRelationship.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrInvalidLeadStatus.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrContractExpired.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrContractNotActive.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrModuleNotContracted.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrAccountSuspended.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrAccountInactive.Code] = http.StatusUnprocessableEntity
	m.statusByCode[domainerror.ErrCompanySuspended.Code] = http.StatusUnprocessableEntity

	// ---------------------------------------------------------
	// 413 / 415 – Arquivo / media type
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrFileTooLarge.Code] = http.StatusRequestEntityTooLarge
	m.statusByCode[domainerror.ErrInvalidFileType.Code] = http.StatusUnsupportedMediaType
	m.statusByCode[domainerror.ErrFileUploadFailed.Code] = http.StatusInternalServerError

	// ---------------------------------------------------------
	// Protocolos HTTP específicos
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrMethodNotAllowed.Code] = http.StatusMethodNotAllowed
	m.statusByCode[domainerror.ErrNotAcceptable.Code] = http.StatusNotAcceptable
	m.statusByCode[domainerror.ErrRequestTimeout.Code] = http.StatusRequestTimeout
	m.statusByCode[domainerror.ErrUnsupportedMediaType.Code] = http.StatusUnsupportedMediaType
	m.statusByCode[domainerror.ErrExpectationFailed.Code] = http.StatusExpectationFailed

	// ---------------------------------------------------------
	// 412 – Precondition Failed / ETag
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrPreconditionFailed.Code] = http.StatusPreconditionFailed
	m.statusByCode[domainerror.ErrETagMismatch.Code] = http.StatusPreconditionFailed

	// ---------------------------------------------------------
	// 424 – Failed Dependency
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrFailedDependency.Code] = http.StatusFailedDependency

	// ---------------------------------------------------------
	// 429 – Rate Limit / Quota
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrRateLimitExceeded.Code] = http.StatusTooManyRequests
	m.statusByCode[domainerror.ErrQuotaExceeded.Code] = http.StatusTooManyRequests
	m.statusByCode[domainerror.ErrMaxAttemptsExceeded.Code] = http.StatusTooManyRequests

	// ---------------------------------------------------------
	// 451 – Legal reasons
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrUnavailableForLegalReasons.Code] = http.StatusUnavailableForLegalReasons

	// ---------------------------------------------------------
	// 500 – Erros internos / banco / APIs externas genéricas
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrInternalServer.Code] = http.StatusInternalServerError
	m.statusByCode[domainerror.ErrDatabaseConnection.Code] = http.StatusInternalServerError
	m.statusByCode[domainerror.ErrDatabaseQuery.Code] = http.StatusInternalServerError
	m.statusByCode[domainerror.ErrThirdPartyAPIError.Code] = http.StatusBadGateway
	m.statusByCode[domainerror.ErrOptimisticLockFailed.Code] = http.StatusConflict // (já mapeado, mas ok)

	// ---------------------------------------------------------
	// 502 / 503 / 504 – serviços externos / indisponibilidade
	// ---------------------------------------------------------
	m.statusByCode[domainerror.ErrExternalServiceUnavailable.Code] = http.StatusServiceUnavailable
	m.statusByCode[domainerror.ErrExternalServiceTimeout.Code] = http.StatusGatewayTimeout
	m.statusByCode[domainerror.ErrServiceUnavailable.Code] = http.StatusServiceUnavailable

	return m
}

func (m *DefaultHTTPStatusMapper) Status(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var derr *domainerror.DomainError
	if !errors.As(err, &derr) {
		return http.StatusInternalServerError
	}

	if status, ok := m.statusByCode[derr.Code]; ok {
		return status
	}
	return http.StatusInternalServerError
}
