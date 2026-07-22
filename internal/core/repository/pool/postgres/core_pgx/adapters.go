package core_pgx

import (
	"TodoList/internal/core/repository/pool/postgres"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...interface{}) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		return mapError(err)
	}
	return err
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

type pgxTx struct {
	pgx.Tx
}

func mapError(err error) error {
	const (
		pgxViolatesForeignKeyErrorCode = "23503"
		pgxSerializationFailureCode    = "40001"
		pgxDeadlockDetectedCode        = "40P01"
		pgxUniqueViolationCode         = "23505"
		pgxNotNullViolationCode        = "23502"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return postgres.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgxViolatesForeignKeyErrorCode:
			return fmt.Errorf("%v: %w", err, postgres.ErrViolatesForeignKey)
		case pgxSerializationFailureCode, pgxDeadlockDetectedCode:
			return fmt.Errorf("%v: %w", err, postgres.ErrTransactionConflict)
		case pgxUniqueViolationCode:
			return fmt.Errorf("%v: %w", err, errors.New("unique constraint violation"))
		case pgxNotNullViolationCode:
			return fmt.Errorf("%v: %w", err, errors.New("not null constraint violation"))
		default:
			return fmt.Errorf("database error (code %s): %v", pgErr.Code, err)
		}
	}

	return err
}
