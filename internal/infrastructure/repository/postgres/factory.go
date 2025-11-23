package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryFactory struct {
	pool *pgxpool.Pool
}

func NewRepositoryFactory(pool *pgxpool.Pool) *RepositoryFactory {
	return &RepositoryFactory{pool}
}

func (f *RepositoryFactory) User() *userRepository {
	return NewUserRepository(f.pool)
}

func (f *RepositoryFactory) Auth() *authRepository {
	return NewAuthRepository(f.pool)
}
