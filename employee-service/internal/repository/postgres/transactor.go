package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transactor struct {
	db *pgxpool.Pool
}

func NewTransactor(db *pgxpool.Pool) *Transactor {
	return &Transactor{db: db}
}

type txKey struct{}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}

	return nil
}

func (t *Transactor) WithTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			return
		}
	}(tx, ctx)

	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		return fmt.Errorf("could not execute transaction: %w", err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}
