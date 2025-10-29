package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/pkg/log"
)

type PGQuerier interface {
	Exec(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(
		ctx context.Context,
		tableName pgx.Identifier,
		columnNames []string,
		rowSrc pgx.CopyFromSource,
	) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type ctxKey struct{}

var txCtxKey ctxKey

type db struct {
	pool *pgxpool.Pool
}

// GetConn returns the TXed connection if exists, otherwise it will return the regular pool connection.
func (r *db) GetConn(ctx context.Context) PGQuerier {
	tx, ok := ctx.Value(txCtxKey).(pgx.Tx)
	if !ok {
		return r.pool
	}
	return tx
}

// GetConnTx returns the transactioned connection. returns error if the connection is not transactioned.
func (r *db) GetConnTx(ctx context.Context) (PGQuerier, error) {
	tx, ok := ctx.Value(txCtxKey).(pgx.Tx)
	if !ok {
		return nil, errors.New("required transaction connection")
	}
	return tx, nil
}

func (r *db) IsTxConn(conn any) bool {
	switch conn.(type) {
	case pgx.Tx:
		return true
	default:
		return false
	}
}

func (r *db) Transact(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	log.DebugCtx(ctx, "Transaction Begin")

	txCtx := context.WithValue(ctx, txCtxKey, tx)
	err = fn(txCtx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx failed: %v, rollback failed: %w", err, rbErr)
		}
		log.DebugCtx(ctx, "Transaction Rollback")
		return fmt.Errorf("tx failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	log.DebugCtx(ctx, "Transaction Commited")

	return nil
}

// Return *db that implement repository.Transcator interface
func NewTransactor(pgpool *pgxpool.Pool) *db {
	return &db{pool: pgpool}
}

// Implements repository.Transactor
// func (r *db) Transact(
// 	ctx context.Context,
// 	fn func(ctx context.Context) error,
// ) error {
// 	conn, err := r.pool.Acquire(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to acquire connection: %w", err)
// 	}
// 	defer conn.Release()
//
// 	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
// 	if err != nil {
// 		return fmt.Errorf("failed to begin transaction: %w", err)
// 	}
//
// 	// Create a done channel to signal when fn is complete
// 	done := make(chan struct{})
//
// 	// Create a channel for the error result
// 	errChan := make(chan error, 1)
//
// 	// Run the transaction function in a goroutine
// 	go func() {
// 		defer close(done)
// 		txCtx := context.WithValue(ctx, txKey{}, tx)
// 		if err := fn(txCtx); err != nil {
// 			errChan <- err
// 			return
// 		}
// 		errChan <- nil
// 	}()
//
// 	// Wait for either the transaction to complete or the context to be canceled
// 	select {
// 	case <-done:
// 		err := <-errChan
// 		if err != nil {
// 			// log.Println("Rolling back transaction due to error:", err)
// 			if rbErr := tx.Rollback(ctx); rbErr != nil {
// 				return fmt.Errorf(
// 					"transaction failed: %v, rollback failed: %w",
// 					err,
// 					rbErr,
// 				)
// 			}
// 			return fmt.Errorf("transaction failed: %w", err)
// 		}
// 		// log.Println("Committing transaction")
// 		if err := tx.Commit(ctx); err != nil {
// 			return fmt.Errorf("failed to commit transaction: %w", err)
// 		}
// 		return nil
// 	case <-ctx.Done():
// 		// log.Println("Rolling back transaction due to context cancellation")
// 		if rbErr := tx.Rollback(ctx); rbErr != nil {
// 			return fmt.Errorf("context canceled, rollback failed: %w", rbErr)
// 		}
// 		return ctx.Err()
// 	}
// }
