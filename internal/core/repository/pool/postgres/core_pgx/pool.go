package core_pgx

import (
	"TodoList/internal/core/repository/pool/postgres"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewPool(
	ctx context.Context,
	config Config,
) (*PgxPool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping connection pool: %w", err)
	}
	return &PgxPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *PgxPool) Query(ctx context.Context, sql string, args ...any) (postgres.Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxRows{rows}, nil
}
func (p *PgxPool) QueryRow(ctx context.Context, sql string, args ...any) postgres.Row {
	row := p.Pool.QueryRow(ctx, sql, args...)
	return pgxRow{row}
}
func (p *PgxPool) Exec(ctx context.Context, sql string, arguments ...any) (postgres.CommandTag, error) {
	commandTag, err := p.Pool.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxCommandTag{commandTag}, nil
}

func (conn *PgxPool) OpTimeout() time.Duration {
	return conn.opTimeout
}
