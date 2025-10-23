package config

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

type AppEnv string

const (
	ENV_PRODUCTION  AppEnv = "PROD"
	ENV_DEVELOPMENT AppEnv = "DEV"
)

type Config struct {
	App      AppConfig     `mapstructure:",squash"`
	Postgres PGConfig      `mapstructure:",squash"`
	Cors     CorsConfig    `mapstructure:",squash"`
	Token    TokenConfig   `mapstructure:",squash"`
	Metrics  MetricsConfig `mapstructure:",squash"`
	SMTP     SMTPConfig    `mapstructure:",squash"`
	R2       R2Config      `mapstructure:",squash"`
}

func (c Config) IsProduction() bool {
	return c.App.Environment == ENV_PRODUCTION
}

type AppConfig struct {
	Name        string `mapstructure:"APP_NAME"`
	Version     string `mapstructure:"APP_VERSION"`
	Port        int    `mapstructure:"APP_PORT"`
	Environment AppEnv `mapstructure:"APP_ENV"`
}

type PGConfig struct {
	User            string        `mapstructure:"DB_USER"`
	Password        string        `mapstructure:"DB_PASSWORD"`
	Host            string        `mapstructure:"DB_HOST"`
	Port            int           `mapstructure:"DB_PORT"`
	Name            string        `mapstructure:"DB_NAME"`
	MinConns        int           `mapstructure:"DB_MINCONNS"`         // PG Pool minimum connections
	MaxConns        int           `mapstructure:"DB_MAXCONNS"`         // PG Pool maximum connections
	MaxConnLifetime time.Duration `mapstructure:"DB_MAXCONN_LIFETIME"` // PG Pool maximun connection lifetime
}

type MetricsConfig struct {
	Enable         bool `mapstructure:"METRICS_ENABLE"`
	PrometheusPort int  `mapstructure:"METRICS_PROMETHEUS_PORT"`
}

type CorsConfig struct {
	Origins     []string `mapstructure:"CORS_ORIGINS"`
	Credentials bool     `mapstructure:"CORS_CREDENTIALS"`
}

type TokenConfig struct {
	SecretKey                 string        `mapstructure:"TOKEN_SECRETKEY"`
	AccessTokenExpiry         time.Duration `mapstructure:"TOKEN_ACCESS_TOKEN_EXPIRY"`
	RefreshTokenExpiry        time.Duration `mapstructure:"TOKEN_REFRESH_TOKEN_EXPIRY"`
	ResetPasswordTokenExpiry  time.Duration `mapstructure:"RESET_PASSWORD_TOKEN_EXPIRY"`
	ResetPasswordFormEndpoint string        `mapstructure:"RESET_PASSWORD_FORM_ENDPOINT"`
}

type SMTPConfig struct {
	Host         string `mapstructure:"SMTP_HOST"`
	Port         int    `mapstructure:"SMTP_PORT"`
	SenderName   string `mapstructure:"SMTP_SENDER_NAME"`
	AuthEmail    string `mapstructure:"SMTP_AUTH_EMAIL"`
	AuthPassword string `mapstructure:"SMTP_AUTH_PASSWORD"`
}

type R2Config struct {
	PublicBucketURL string `mapstructure:"R2_PUBLIC_BUCKET_URL"`
	PublicBucket    string `mapstructure:"R2_PUBLIC_BUCKET"`
	PrivateBucket   string `mapstructure:"R2_PRIVATE_BUCKET"`
	AccountID       string `mapstructure:"R2_ACCOUNT_ID"`
	AccessKeyID     string `mapstructure:"R2_ACCESS_KEY_ID"`
	AccessKeySecret string `mapstructure:"R2_ACCESS_KEY_SECRET"`
}

// Load and Parse Config, pass the path of the config file relatively from the root dir
func LoadConfig(filepath string) (*Config, error) {
	// Set the file name and path for the .env file
	viper.SetConfigFile(filepath)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("fail parse config: %v", err.Error())
	}

	if c.App.Environment != ENV_PRODUCTION && c.App.Environment != ENV_DEVELOPMENT {
		return nil, errors.New("invalid app.Environtment value, expecting 'DEV' or 'PROD'")
	}

	// Validate origins URL
	for _, origin := range c.Cors.Origins {
		if _, err := url.ParseRequestURI(origin); err != nil {
			return nil, fmt.Errorf("invalid cors origins url: %s", origin)
		}
	}

	return &c, nil
}
