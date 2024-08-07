package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/pkg/errors"
	"github.com/prawirdani/golang-restapi/pkg/logging"
)

var (
	ErrorUserNotFound = errors.NotFound("User not found")
)

type UserField string

const (
	UserID    UserField = "id"
	UserEmail UserField = "email"
	UserName  UserField = "name"
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
		return err
	}
	return nil
}

func (r *UserRepository) SelectWhere(ctx context.Context, field UserField, searchVal any) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, name, email, password, created_at, updated_at FROM users WHERE %s=$1", field)

	err := r.db.QueryRow(ctx, query, searchVal).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrorUserNotFound
		}
		r.logger.Error(logging.Postgres, "UserRepository.SelectWhere", err.Error())
		return nil, err
	}

	return &user, nil
}
