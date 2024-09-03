package repository

import (
	"context"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/pkg/common"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

type UserRepository struct {
	db     *pgxpool.Pool
	logger logging.Logger
}

func NewUserRepository(pgpool *pgxpool.Pool, logger logging.Logger) *UserRepository {
	return &UserRepository{
		db:     pgpool,
		logger: logger,
	}
}

func (r *UserRepository) InsertUser(ctx context.Context, u entity.User) error {
	query := "INSERT INTO users(id, name, email, password) VALUES($1, $2, $3, $4)"
	_, err := r.db.Exec(ctx, query, u.ID, u.Name, u.Email, u.Password)

	if err != nil {
		r.logger.Error(logging.Postgres, "UserRepository.InsertUser", err.Error())
		if strings.Contains(err.Error(), "users_email_key") {
			return entity.ErrEmailExist
		}
		return err
	}
	return nil
}

func (r *UserRepository) SelectWhere(ctx context.Context, field string, searchVal any) (entity.User, error) {
	var user entity.User
	query := common.ConcatStrings("SELECT id, name, email, password, created_at, updated_at FROM users WHERE ", field, "=$1")

	if err := pgxscan.Get(ctx, r.db, &user, query, searchVal); err != nil {
		if pgxscan.NotFound(err) {
			return entity.User{}, entity.ErrUserNotFound
		}
		r.logger.Error(logging.Postgres, "UserRepository.SelectWhere", err.Error())
		return entity.User{}, err
	}

	return user, nil
}
