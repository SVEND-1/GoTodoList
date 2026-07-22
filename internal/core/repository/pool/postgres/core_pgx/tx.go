package core_pgx

import (
	"TodoList/internal/core/repository/pool/postgres"
	"context"
)

func (t pgxTx) Query(ctx context.Context, sql string, args ...any) (postgres.Rows, error) {
	rows, err := t.Tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxRows{rows}, nil
}

func (t pgxTx) QueryRow(ctx context.Context, sql string, args ...any) postgres.Row {
	row := t.Tx.QueryRow(ctx, sql, args...)
	return pgxRow{row}
}

func (t pgxTx) Exec(ctx context.Context, sql string, arguments ...any) (postgres.CommandTag, error) {
	commandTag, err := t.Tx.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxCommandTag{commandTag}, nil
}

func (t pgxTx) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func (t pgxTx) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}
