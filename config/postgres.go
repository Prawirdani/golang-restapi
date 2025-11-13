package config

import (
	"os"
	"strconv"
	"time"
)

type Postgres struct {
	User            string
	Password        string
	Host            string
	Port            int
	Name            string
	MinConns        int
	MaxConns        int
	MaxConnLifetime time.Duration
}

func (p *Postgres) Parse() error {
	p.User = os.Getenv("DB_USER")
	p.Password = os.Getenv("DB_PASSWORD")
	p.Host = os.Getenv("DB_HOST")
	p.Name = os.Getenv("DB_NAME")

	if val := os.Getenv("DB_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			p.Port = port
		}
	}
	if val := os.Getenv("DB_MINCONNS"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			p.MinConns = i
		}
	}
	if val := os.Getenv("DB_MAXCONNS"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			p.MaxConns = i
		}
	}
	if val := os.Getenv("DB_MAXCONN_LIFETIME"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			p.MaxConnLifetime = d
		}
	}
	return nil
}
