package core_pgx

import (
	"TodoList/internal/core/repository/pool/postgres"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func (p *PgxPool) Begin(ctx context.Context) (postgres.Tx, error) {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return pgxTx{tx}, nil
}

func (p *PgxPool) BeginTx(ctx context.Context, opts *postgres.TxOptions) (postgres.Tx, error) {
	var txOptions pgx.TxOptions

	if opts != nil {
		switch opts.IsolationLevel {
		case postgres.ReadUncommitted:
			txOptions.IsoLevel = pgx.ReadUncommitted
		case postgres.ReadCommitted:
			txOptions.IsoLevel = pgx.ReadCommitted
		case postgres.RepeatableRead:
			txOptions.IsoLevel = pgx.RepeatableRead
		case postgres.Serializable:
			txOptions.IsoLevel = pgx.Serializable
		}
	}

	tx, err := p.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return pgxTx{tx}, nil
}

func (p *PgxPool) BeginFunc(ctx context.Context, f func(postgres.Tx) error) error {
	tx, err := p.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
	}()

	if err := f(tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("error in transaction: %w (rollback error: %v)", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

func (p *PgxPool) BeginTxFunc(ctx context.Context, opts *postgres.TxOptions, f func(postgres.Tx) error) error {
	const maxRetries = 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		tx, err := p.BeginTx(ctx, opts)
		if err != nil {
			return err
		}

		var txErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					_ = tx.Rollback(ctx)
					panic(r)
				}
			}()

			txErr = f(tx)
		}()

		if txErr != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				return fmt.Errorf("error in transaction: %w (rollback error: %v)", txErr, rbErr)
			}

			if errors.Is(txErr, postgres.ErrTransactionConflict) && attempt < maxRetries-1 {
				time.Sleep(time.Duration(attempt+1) * 10 * time.Millisecond)
				continue
			}
			return txErr
		}

		return tx.Commit(ctx)
	}

	return fmt.Errorf("max retries exceeded for transaction")
}

func (conn *PgxPool) OpTimeout() time.Duration {
	return conn.opTimeout
}

func (p *PgxPool) Close() {
	p.Pool.Close()
}
