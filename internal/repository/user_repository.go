package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/entity"
	"github.com/prawirdani/golang-restapi/pkg/httputil"
)

type UserRepository struct {
	tableName string
}

func NewUserRepository(tableName string) *UserRepository {
	return &UserRepository{
		tableName: tableName,
	}
}

func (r *UserRepository) Create(ctx context.Context, u entity.User, tx pgx.Tx) error {
	query := fmt.Sprintf("INSERT INTO %s(id, name, email, password) VALUES($1, $2, $3, $4)", r.tableName)
	_, err := tx.Exec(ctx, query, u.ID, u.Name, u.Email, u.Password)
	if err != nil {
		// Unique constraint error checker by PG error code.
		if strings.Contains(err.Error(), "23505") {
			return httputil.ErrConflict(fmt.Errorf("Email already exists"))
		}
		return err
	}
	return nil
}

func (r *UserRepository) SelectById(ctx context.Context, userId string, db *pgxpool.Pool) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, name, email, password, created_at, updated_at FROM %s WHERE id=$1", r.tableName)

	if err := db.QueryRow(ctx, query, userId).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) SelectByEmail(ctx context.Context, email string, db *pgxpool.Pool) (*entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, name, email, password, created_at, updated_at FROM %s WHERE email=$1", r.tableName)

	if err := db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}
