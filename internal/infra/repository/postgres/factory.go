package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

type RepositoryFactory struct {
	pool   *pgxpool.Pool
	logger logging.Logger
}

func NewRepositoryFactory(pool *pgxpool.Pool, logger logging.Logger) *RepositoryFactory {
	return &RepositoryFactory{pool, logger}
}

func (f *RepositoryFactory) User() *userRepository {
	return NewUserRepository(f.pool, f.logger)
}

func (f *RepositoryFactory) Auth() *authRepository {
	return NewAuthRepository(f.pool, f.logger)
}
