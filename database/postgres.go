package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
)

// Return PostgreSQL database pooling
func NewPGConnection(cfg *config.Config) (*pgxpool.Pool, error) {
	// DSN Format postgres://username:password@localhost:5432/db_name
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?application_name=%s",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.App.Name,
	)

	pgConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pgConf.MinConns = int32(cfg.DB.MinConns)
	pgConf.MaxConns = int32(cfg.DB.MaxConns)
	pgConf.MaxConnLifetime = time.Minute * time.Duration(cfg.DB.MaxConnLifetime)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConf)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
