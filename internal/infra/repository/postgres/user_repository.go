package postgres

import (
	"context"
	"strings"
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
func (r *userRepository) Insert(ctx context.Context, u user.User) error {
	// log.DebugCtx(ctx, "insert user", log.JSON("args", u))

	query := "INSERT INTO users(id, name, email, phone, password, profile_image) VALUES($1, $2, $3, $4, $5, $6)"
	conn := r.db.GetConn(ctx)

	_, err := conn.Exec(ctx, query, u.ID, u.Name, u.Email, u.Phone, u.Password, u.ProfileImage)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return user.ErrEmailExist
		}
		log.ErrorCtx(ctx, "failed to insert user", "error", err.Error())
		return err
	}
	return nil
}

// Implements user.Repository
func (r *userRepository) GetUserBy(
	ctx context.Context,
	field string,
	value any,
) (user.User, error) {
	log.DebugCtx(ctx, "get user data", "search_field", field, "search_arg", value)

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

	if err := pgxscan.Get(ctx, conn, &u, query, value); err != nil {
		if pgxscan.NotFound(err) {
			return user.User{}, user.ErrUserNotFound
		}
		log.ErrorCtx(ctx, "failed to get user", "field", field, "error", err.Error())
		return user.User{}, err
	}

	return u, nil
}

// Implements user.Repository.
func (r *userRepository) UpdateUser(ctx context.Context, u user.User) error {
	// log.DebugCtx(ctx, "update user", log.JSON("args", u))

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
		if strings.Contains(err.Error(), "users_email_key") {
			return user.ErrEmailExist
		}
		log.ErrorCtx(ctx, "failed to update user", "error", err.Error())
		return err
	}

	return nil
}

// Implements user.Repository
func (r *userRepository) DeleteUser(ctx context.Context, userId string) error {
	log.WarnCtx(ctx, "deleting user data", "user_id", userId)

	query := "UPDATE users SET deleted_at=$1 WHERE id=$2"
	conn := r.db.GetConn(ctx)

	deleteTime := time.Now()

	_, err := conn.Exec(ctx, query, deleteTime, userId)
	if err != nil {
		log.ErrorCtx(ctx, "failed to delete user", "error", err.Error())
		return err
	}

	return nil
}
