package core_pgx

import (
	"TodoList/internal/core/repository/pool/postgres"
	"context"
)

type txKey struct{}

func WithTx(ctx context.Context, tx postgres.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func TxFromContext(ctx context.Context) (postgres.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(postgres.Tx)
	return tx, ok
}

func (p *PgxPool) ExecutorFromContext(ctx context.Context) postgres.Executor {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return p
}
