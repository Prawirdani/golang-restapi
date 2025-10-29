package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/entity/user"
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

// Implements user.Repository
func (r *userRepository) Insert(ctx context.Context, u *user.User) error {
	if u == nil {
		log.WarnCtx(ctx, "Insert user called with nil user")
		return errors.New("user is nil")
	}

	log.DebugCtx(ctx, "Insert user into db", "args", u)

	query := "INSERT INTO users(id, name, email, phone, password, profile_image) VALUES($1, $2, $3, $4, $5, $6)"
	conn := r.db.GetConn(ctx)

	_, err := conn.Exec(ctx, query, u.ID, u.Name, u.Email, u.Phone, u.Password, u.ProfileImage)
	if err != nil {
		if uniqueViolationErr(err, "users_email_key") {
			return user.ErrEmailExist
		}

		log.ErrorCtx(ctx, "Failed to insert user", "error", err.Error())
		return err
	}
	return nil
}

// Implements user.Repository
func (r *userRepository) GetUserBy(
	ctx context.Context,
	field string,
	value any,
) (*user.User, error) {
	var u user.User
	query := strs.Concatenate(
		"SELECT id, name, email, phone, password, profile_image, created_at, updated_at, deleted_at FROM users WHERE ",
		field,
		"=$1",
	)
	conn := r.db.GetConn(ctx)
	if r.db.IsTxConn(conn) {
		query += "\nFOR UPDATE"
	}

	log.DebugCtx(ctx, "Get user data", "search_field", field, "search_arg", value)
	if err := pgxscan.Get(ctx, conn, &u, query, value); err != nil {
		if noRowsErr(err) {
			return nil, user.ErrUserNotFound
		}
		log.ErrorCtx(ctx, "Failed to get user", "field", field, "error", err.Error())
		return nil, err
	}

	return &u, nil
}

// Implements user.Repository.
func (r *userRepository) UpdateUser(ctx context.Context, u *user.User) error {
	if u == nil {
		log.WarnCtx(ctx, "Update user called with nil user")
		return errors.New("user is nil")
	}

	log.DebugCtx(ctx, "Updating user data", "args", u)

	query := "UPDATE users SET name=$1, email=$2, phone=$3, password=$4, profile_image=$5, updated_at=$6 WHERE id=$7"
	updateTime := time.Now()

	conn := r.db.GetConn(ctx)
	_, err := conn.Exec(
		ctx,
		query,
		u.Name,
		u.Email,
		u.Phone,
		u.Password,
		u.ProfileImage,
		updateTime,
		u.ID,
	)
	if err != nil {
		if uniqueViolationErr(err, "users_email_key") {
			return user.ErrEmailExist
		}

		log.ErrorCtx(ctx, "Failed to update user", "error", err.Error())
		return err
	}

	return nil
}

// Implements user.Repository
func (r *userRepository) DeleteUser(ctx context.Context, u *user.User) error {
	if u == nil {
		log.WarnCtx(ctx, "Delete user called with nil user")
		return errors.New("user is nil")
	}

	log.WarnCtx(ctx, "Deleting user data", "id", u.ID.String(), "name", u.Name)

	query := "UPDATE users SET deleted_at=$1 WHERE id=$2"
	conn := r.db.GetConn(ctx)

	deleteTime := time.Now()

	_, err := conn.Exec(ctx, query, deleteTime, u.ID)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to delete user", "error", err.Error())
		return err
	}

	return nil
}
