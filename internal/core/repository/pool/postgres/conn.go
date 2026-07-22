package postgres

import (
	"context"
	"time"
)

type Pool interface {
	Executor
	Close()

	Begin(ctx context.Context) (Tx, error)
	BeginTx(ctx context.Context, opts *TxOptions) (Tx, error)
	BeginFunc(ctx context.Context, f func(Tx) error) error
	BeginTxFunc(ctx context.Context, opts *TxOptions, f func(Tx) error) error

	ExecutorFromContext(ctx context.Context) Executor

	OpTimeout() time.Duration
}

type Rows interface {
	Scan(dest ...any) error
	Next() bool
	Err() error
	Close()
}

type Row interface {
	Scan(dest ...any) error
}

type CommandTag interface {
	RowsAffected() int64
}

type Tx interface {
	Executor
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type TxOptions struct {
	IsolationLevel IsolationLevel
}

type IsolationLevel int

const (
	ReadUncommitted IsolationLevel = iota
	ReadCommitted
	RepeatableRead
	Serializable
)
