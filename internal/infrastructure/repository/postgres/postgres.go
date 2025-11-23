package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
)

// NewPool Return PostgreSQL database pooling
func NewPool(cfg config.Postgres) (*pgxpool.Pool, error) {
	// DSN Format postgres://username:password@localhost:5432/db_name
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	pgConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pgConf.MinConns = int32(cfg.MinConns)
	pgConf.MaxConns = int32(cfg.MaxConns)
	pgConf.MaxConnLifetime = cfg.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConf)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
