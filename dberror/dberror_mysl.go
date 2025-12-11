package dberror

import (
	"errors"
	"strings"

	mysql "github.com/go-sql-driver/mysql"
	domainerror "github.com/renatofagalde/module-error"
)

type MySQLErrorMapper struct {
	duplicateIndexErrors map[string]*domainerror.DomainError
}

func NewMySQLErrorMapper(duplicateIndexErrors map[string]*domainerror.DomainError) DBErrorMapper {
	return &MySQLErrorMapper{
		duplicateIndexErrors: duplicateIndexErrors,
	}
}

func (m *MySQLErrorMapper) Map(err error) error {
	if err == nil {
		return nil
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {

		// Duplicate entry
		case 1062:
			msg := mysqlErr.Message

			if m.duplicateIndexErrors != nil {
				for indexName, derr := range m.duplicateIndexErrors {
					if derr != nil && strings.Contains(msg, indexName) {
						return derr
					}
				}
			}

			return domainerror.ErrConflict

		case 1048:
			return domainerror.ErrRequiredField
		}
	}

	return domainerror.ErrDatabaseQuery
}
