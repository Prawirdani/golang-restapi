package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

// Return PostgreSQL database pooling
func NewPGPool(config *viper.Viper) *pgxpool.Pool {
	// DSN Format postgres://username:password@localhost:5432/db_name
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s",
		config.GetString("db.username"),
		config.GetString("db.password"),
		config.GetString("db.host"),
		config.GetInt("db.port"),
		config.GetString("db.name"),
	)

	pgConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("Error parsing postgres dns address", err)
	}

	pgConf.MinConns = config.GetInt32("db.pool.min")
	pgConf.MaxConns = config.GetInt32("db.pool.max")
	pgConf.MaxConnLifetime = time.Minute * time.Duration(config.GetInt("db.pool.lifetime"))

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
