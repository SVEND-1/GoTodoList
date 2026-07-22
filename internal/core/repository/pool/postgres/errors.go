package postgres

import "errors"

var (
	ErrNoRows              = errors.New("no rows")
	ErrViolatesForeignKey  = errors.New("violates foreign key")
	ErrTransactionConflict = errors.New("transaction conflict, please retry")
	ErrUnknown             = errors.New("unknown")
)
