package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
)

// Return PostgreSQL database pooling
func NewPGConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	// DSN Format postgres://username:password@localhost:5432/db_name
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?application_name=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Name,
		cfg.App.Name,
	)

	pgConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pgConf.MinConns = int32(cfg.Postgres.MinConns)
	pgConf.MaxConns = int32(cfg.Postgres.MaxConns)
	pgConf.MaxConnLifetime = cfg.Postgres.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConf)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
