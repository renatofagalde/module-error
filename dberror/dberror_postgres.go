package dberror

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	domainerror "github.com/renatofagalde/module-error"
	"gorm.io/gorm"
)

type PostgresErrorMapper struct {
	constraintErrors map[string]*domainerror.DomainError
}

func NewPostgresErrorMapper(constraintErrors map[string]*domainerror.DomainError) DBErrorMapper {
	return &PostgresErrorMapper{
		constraintErrors: constraintErrors,
	}
}

func (m *PostgresErrorMapper) Map(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domainerror.ErrConflict
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {

		case "23505":
			if m.constraintErrors != nil {
				if derr, ok := m.constraintErrors[pgErr.ConstraintName]; ok && derr != nil {
					return derr
				}
			}
			// Fallback genérico
			return domainerror.ErrConflict

		case "23503":
			return domainerror.ErrInvalidRelationship

		case "23502":
			return domainerror.ErrRequiredField
		}
	}

	if errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) {
		return domainerror.ErrRequestTimeout
	}

	return domainerror.ErrDatabaseQuery
}
