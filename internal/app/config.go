package app

import (
	"log"

	"github.com/spf13/viper"
)

func NewConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("fatal config error: %v", err.Error())
	}

	return v
}
