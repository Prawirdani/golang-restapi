package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type Configuration struct {
	MainRouter *chi.Mux
	Config     *viper.Viper
	DBPool     *pgxpool.Pool
	Validator  *validator.Validate
}

func Bootstrap(d *Configuration) {
}
