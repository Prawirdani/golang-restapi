package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/domain/user"
	"github.com/prawirdani/golang-restapi/pkg/log"
	strs "github.com/prawirdani/golang-restapi/pkg/strings"
)

type userRepository struct {
	db *db
}

func NewUserRepository(connPool *pgxpool.Pool) *userRepository {
	return &userRepository{
		db: &db{pool: connPool},
	}
}

// Store implements [user.Repository].
func (r *userRepository) Store(ctx context.Context, u *user.User) error {
	if u == nil {
		log.WarnCtx(ctx, "Store called with nil user")
		return errors.New("user is nil")
	}

	query := "INSERT INTO users(id, name, email, phone, password, profile_image) VALUES($1, $2, $3, $4, $5, $6)"
	conn := r.db.GetConn(ctx)

	_, err := conn.Exec(ctx, query, u.ID, u.Name, u.Email, u.Phone, u.Password, u.ProfileImage)
	if err != nil {
		if uniqueViolationErr(err, "users_email_key") {
			return user.ErrEmailExists
		}

		log.ErrorCtx(ctx, "Failed to store user", err)
		return err
	}
	return nil
}

// GetByEmail implements [user.Repository].
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return r.getUserBy(ctx, "email", email)
}

// GetByID implements [user.Repository].
func (r *userRepository) GetByID(ctx context.Context, userID string) (*user.User, error) {
	return r.getUserBy(ctx, "id", userID)
}

// Update implements [user.Repository].
func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	if u == nil {
		log.WarnCtx(ctx, "Update called with nil user")
		return errors.New("user is nil")
	}

	query := "UPDATE users SET name=$1, email=$2, phone=$3, password=$4, profile_image=$5, updated_at=$6 WHERE id=$7"
	updatedAt := time.Now()

	conn := r.db.GetConn(ctx)
	_, err := conn.Exec(
		ctx,
		query,
		u.Name,
		u.Email,
		u.Phone,
		u.Password,
		u.ProfileImage,
		updatedAt,
		u.ID,
	)
	if err != nil {
		if uniqueViolationErr(err, "users_email_key") {
			return user.ErrEmailExists
		}

		log.ErrorCtx(ctx, "Failed to update user", err)
		return err
	}

	return nil
}

// Delete implements [user.Repository].
func (r *userRepository) Delete(ctx context.Context, u *user.User) error {
	if u == nil {
		log.WarnCtx(ctx, "Delete user called with nil user")
		return errors.New("user is nil")
	}

	query := "UPDATE users SET deleted_at=$1 WHERE id=$2"
	conn := r.db.GetConn(ctx)

	deleteTime := time.Now()
	_, err := conn.Exec(ctx, query, deleteTime, u.ID)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to delete user", err)
		return err
	}

	return nil
}

func (r *userRepository) getUserBy(
	ctx context.Context,
	field string,
	value any,
) (*user.User, error) {
	query := strs.Concatenate(
		"SELECT id, name, email, phone, password, profile_image, created_at, updated_at FROM users WHERE ",
		field,
		"=$1",
	)
	conn := r.db.GetConn(ctx)
	if r.db.IsTxConn(conn) {
		query += "\nFOR UPDATE"
	}

	var u user.User
	if err := pgxscan.Get(ctx, conn, &u, query, value); err != nil {
		if noRowsErr(err) {
			return nil, user.ErrNotFound
		}
		log.ErrorCtx(ctx, "Failed to get user", err, "field", field)
		return nil, err
	}

	return &u, nil
}
