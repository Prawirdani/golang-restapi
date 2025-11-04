package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type AppEnv string

const (
	ENV_PRODUCTION  AppEnv = "prod"
	ENV_DEVELOPMENT AppEnv = "dev"
)

type Config struct {
	App         AppConfig
	Postgres    PGConfig
	Cors        CorsConfig
	Token       TokenConfig
	Metrics     MetricsConfig
	SMTP        SMTPConfig
	R2          R2Config
	RabbitMqURL string
}

func (c Config) IsProduction() bool {
	return c.App.Environment == ENV_PRODUCTION
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load() // Load .env in dev

	cfg := &Config{}

	// Parse each struct
	if err := cfg.App.Parse(); err != nil {
		return nil, err
	}
	if err := cfg.Postgres.Parse(); err != nil {
		return nil, err
	}
	if err := cfg.Metrics.Parse(); err != nil {
		return nil, err
	}
	if err := cfg.Cors.Parse(); err != nil {
		return nil, err
	}
	if err := cfg.Token.Parse(); err != nil {
		return nil, err
	}
	if err := cfg.SMTP.Parse(); err != nil {
		return nil, err
	}
	if err := cfg.R2.Parse(); err != nil {
		return nil, err
	}

	cfg.RabbitMqURL = os.Getenv("RABBITMQ_URL")

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.App.Environment != ENV_PRODUCTION && c.App.Environment != ENV_DEVELOPMENT {
		return fmt.Errorf("invalid APP_ENV, expecting %s or %s", ENV_DEVELOPMENT, ENV_PRODUCTION)
	}
	for _, origin := range c.Cors.Origins {
		if _, err := url.ParseRequestURI(origin); err != nil {
			log.Printf("warning: invalid CORS origin: %s\n", origin)
		}
	}
	return nil
}

// =======================
// AppConfig
// =======================

type AppConfig struct {
	Name        string
	Version     string
	Port        int
	Environment AppEnv
}

func (a *AppConfig) Parse() error {
	a.Name = os.Getenv("APP_NAME")
	a.Version = os.Getenv("APP_VERSION")
	a.Environment = AppEnv(strings.ToLower(os.Getenv("APP_ENV")))

	if val := os.Getenv("APP_PORT"); val != "" {
		port, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		a.Port = port
	}
	return nil
}

// =======================
// Postgres
// =======================

type PGConfig struct {
	User            string
	Password        string
	Host            string
	Port            int
	Name            string
	MinConns        int
	MaxConns        int
	MaxConnLifetime time.Duration
}

func (p *PGConfig) Parse() error {
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

// =======================
// Metrics
// =======================

type MetricsConfig struct {
	PrometheusPort int
}

func (m *MetricsConfig) Parse() error {
	if val := os.Getenv("METRICS_PROMETHEUS_PORT"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			m.PrometheusPort = i
		}
	}
	return nil
}

// =======================
// CORS
// =======================

type CorsConfig struct {
	Origins     []string
	Credentials bool
}

func (c *CorsConfig) Parse() error {
	if val := os.Getenv("CORS_ORIGINS"); val != "" {
		c.Origins = strings.Split(val, ",")
	}
	if val := os.Getenv("CORS_CREDENTIALS"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			c.Credentials = b
		}
	}
	return nil
}

// =======================
// Token
// =======================

type TokenConfig struct {
	SecretKey                 string
	AccessTokenExpiry         time.Duration
	RefreshTokenExpiry        time.Duration
	ResetPasswordTokenExpiry  time.Duration
	ResetPasswordFormEndpoint string
}

func (t *TokenConfig) Parse() error {
	t.SecretKey = os.Getenv("TOKEN_SECRETKEY")
	t.ResetPasswordFormEndpoint = os.Getenv("RESET_PASSWORD_FORM_ENDPOINT")

	if val := os.Getenv("TOKEN_ACCESS_TOKEN_EXPIRY"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			t.AccessTokenExpiry = d
		}
	}
	if val := os.Getenv("TOKEN_REFRESH_TOKEN_EXPIRY"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			t.RefreshTokenExpiry = d
		}
	}
	if val := os.Getenv("RESET_PASSWORD_TOKEN_EXPIRY"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			t.ResetPasswordTokenExpiry = d
		}
	}
	return nil
}

// =======================
// SMTP
// =======================

type SMTPConfig struct {
	Host         string
	Port         int
	SenderName   string
	AuthEmail    string
	AuthPassword string
}

func (s *SMTPConfig) Parse() error {
	s.Host = os.Getenv("SMTP_HOST")
	s.SenderName = os.Getenv("SMTP_SENDER_NAME")
	s.AuthEmail = os.Getenv("SMTP_AUTH_EMAIL")
	s.AuthPassword = os.Getenv("SMTP_AUTH_PASSWORD")

	if val := os.Getenv("SMTP_PORT"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			s.Port = i
		}
	}
	return nil
}

// =======================
// R2
// =======================

type R2Config struct {
	PublicBucketURL string
	PublicBucket    string
	PrivateBucket   string
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
}

func (r *R2Config) Parse() error {
	r.PublicBucketURL = os.Getenv("R2_PUBLIC_BUCKET_URL")
	r.PublicBucket = os.Getenv("R2_PUBLIC_BUCKET")
	r.PrivateBucket = os.Getenv("R2_PRIVATE_BUCKET")
	r.AccountID = os.Getenv("R2_ACCOUNT_ID")
	r.AccessKeyID = os.Getenv("R2_ACCESS_KEY_ID")
	r.AccessKeySecret = os.Getenv("R2_ACCESS_KEY_SECRET")
	return nil
}
