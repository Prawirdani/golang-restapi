package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func uniqueViolationErr(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == constraintName {
		return true
	}

	return false
}

func noRowsErr(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
