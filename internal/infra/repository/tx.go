package repository

import (
	"context"
)

// Transactor defines an interface for executing functions within a database transaction.
// It abstracts transaction handling to ensure consistency across repositories.
// Implementator Usage:
//   - Begins a new transaction and injects it into the provided context.
//   - Commits the transaction if the function returns nil.
//   - Rolls back the transaction if the function returns an error.
//
// IMPORTANT:
//   - Use the provided `ctx` when calling repository methods to ensure the transactional connection is used.
//   - Using the original context may result in executing outside the transaction.
type Transactor interface {
	Transact(ctx context.Context, fn func(ctx context.Context) error) error
}
