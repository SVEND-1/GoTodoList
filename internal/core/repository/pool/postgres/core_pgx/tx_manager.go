package core_pgx

import (
	"TodoList/internal/core/repository/pool/postgres"
	"context"
)

type PgxTxManager struct {
	pool postgres.Pool
}

func NewPgxTxManager(pool postgres.Pool) *PgxTxManager {
	return &PgxTxManager{pool: pool}
}

func (m *PgxTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.pool.BeginFunc(ctx, func(tx postgres.Tx) error {
		txCtx := WithTx(ctx, tx)
		return fn(txCtx)
	})
}
