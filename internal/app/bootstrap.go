package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type Configuration struct {
	MainRouter *chi.Mux
	Config     *viper.Viper
	DBPool     *pgxpool.Pool
	JSON       *JsonHandler
}

func Bootstrap(d *Configuration) {
}
