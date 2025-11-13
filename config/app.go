package config

import (
	"os"
	"strconv"
	"strings"
)

type App struct {
	Name        string
	Version     string
	Port        int
	Environment AppEnv
}

func (a *App) Parse() error {
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
