package config

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Cors  CorsConfig
	Token TokenConfig
}

type AppConfig struct {
	Version     string
	Port        int
	Environment string
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
	AllowedOrigins string
	AllowedMethods string
	Credentials    bool
}

// Convert AllowedOrigins into Array of string
func (cc CorsConfig) OriginsToArray() []string {
	origins := strings.Split(cc.AllowedOrigins, ",")
	// Validate Origins URL
	for _, origin := range origins {
		_, err := url.ParseRequestURI(origin)
		if err != nil {
			log.Fatal(err)
		}
	}
	return origins
}

// Convert AllowedMethods into Array of string
func (cc CorsConfig) MethodsToArray() []string {
	return strings.Split(cc.AllowedMethods, ",")
}

type TokenConfig struct {
	SecretKey string
	Expiry    int // In Hour
}

// Load and Parse Config
func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(path) // Respectfully from the root directory

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fail load config: %v", err.Error())
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("fail parse config: %v", err.Error())
	}
	return &c, nil
}
