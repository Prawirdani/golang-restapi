package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/config"
)

// Return PostgreSQL database pooling
func NewPGPool(cfg config.DBConfig) *pgxpool.Pool {
	// DSN Format postgres://username:password@localhost:5432/db_name
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	pgConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("Error parsing postgres dns address", err)
	}

	pgConf.MinConns = int32(cfg.MinConns)
	pgConf.MaxConns = int32(cfg.MaxConns)
	pgConf.MaxConnLifetime = time.Minute * time.Duration(cfg.MaxConnLifetime)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConf)
	if err != nil {
		slog.Error("PGSQL Init Failed", "cause", err)
		os.Exit(1)
	}

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("PostgreSQL Ping error", "cause", err)
		os.Exit(1)
	}

	slog.Info("PostgreSQL DB Connection Established")
	return pool
}
