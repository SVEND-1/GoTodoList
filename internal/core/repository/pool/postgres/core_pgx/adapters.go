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

func mapError(err error) error {
	const (
		pgxViolatesForeignKeyErrorCode = "23503"
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return postgres.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgxViolatesForeignKeyErrorCode { //При попытки создать задачу и передать userId которого нету в бд
			return fmt.Errorf("%v: %w", err, postgres.ErrViolatesForeignKey)
		}
	}
	return fmt.Errorf("%v: %w", err, pgx.ErrNoRows)
}
