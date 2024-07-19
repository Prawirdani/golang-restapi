package config

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

type AppEnv string

const (
	ENV_PRODUCTION  AppEnv = "PROD"
	ENV_DEVELOPMENT AppEnv = "DEV"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Cors  CorsConfig
	Token TokenConfig
}

func (c Config) IsProduction() bool {
	return c.App.Environment == ENV_PRODUCTION
}

type AppConfig struct {
	Name        string
	Version     string
	Port        int
	Environment AppEnv
}

type DBConfig struct {
	Username        string
	Password        string
	Host            string
	Port            int
	Name            string
	MinConns        int // PG Pool minimum connections
	MaxConns        int // PG Pool maximum connections
	MaxConnLifetime int // PG Pool maximun connection lifetime, In Minute
}

type CorsConfig struct {
	Origins     []string
	Credentials bool
}

type TokenConfig struct {
	SecretKey          string
	AccessTokenExpiry  int
	RefreshTokenExpiry int
}

// Load and Parse Config, pass the path of the config file relatively from the root dir
func LoadConfig(path string) (*Config, error) {
	var c Config
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fail load config: %v", err.Error())
	}

	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("fail parse config: %v", err.Error())
	}

	if c.App.Environment != ENV_PRODUCTION && c.App.Environment != ENV_DEVELOPMENT {
		return nil, errors.New("Invalid app.Environtment value, expecting DEV or PROD")
	}

	// Validate origins URL
	for _, origin := range c.Cors.Origins {
		if _, err := url.ParseRequestURI(origin); err != nil {
			return nil, fmt.Errorf("Invalid cors.Origins URL: %s", origin)
		}
	}

	return &c, nil
}
