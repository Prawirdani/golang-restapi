package config

import (
	"log"

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

type TokenConfig struct {
	SecretKey string
	Expiry    int // In Hour
}

var Cfg *Config

func LoadConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath("./config") // Respectfully from the root directory

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("fatal config error: %v", err.Error())
	}
	return v
}

func ParseConfig(v *viper.Viper) {
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("fail parse config: %v", err.Error())
	}
	Cfg = &c
}
