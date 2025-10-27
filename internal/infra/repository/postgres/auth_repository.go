package postgres

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/pkg/log"
	strs "github.com/prawirdani/golang-restapi/pkg/strings"
)

type authRepository struct {
	db *db
}

func NewAuthRepository(pool *pgxpool.Pool) *authRepository {
	return &authRepository{
		db: &db{pool: pool},
	}
}

// Implements auth.Repository
func (r *authRepository) InsertSession(ctx context.Context, sess auth.Session) error {
	query := "INSERT INTO sessions(user_id, refresh_token, user_agent, expires_at, accessed_at) VALUES($1, $2, $3, $4, $5)"

	conn := r.db.GetConn(ctx)
	if _, err := conn.Exec(ctx, query,
		sess.UserID,
		sess.RefreshToken,
		sess.UserAgent,
		sess.ExpiresAt,
		sess.AccessedAt,
	); err != nil {
		log.Error("failed to insert session", "error", err.Error())
		return err
	}

	return nil
}

// Implements auth.Repository
func (r *authRepository) GetUserSessionBy(
	ctx context.Context,
	field string,
	searchVal any,
) (auth.Session, error) {
	query := strs.Concatenate(
		"UPDATE sessions SET accessed_at=NOW() WHERE ",
		field,
		"=$1 RETURNING *",
	)

	conn := r.db.GetConn(ctx)
	if r.db.IsTxConn(conn) {
		query += "\nFOR UPDATE"
	}

	var sess auth.Session
	if err := pgxscan.Get(ctx, conn, &sess, query, searchVal); err != nil {
		if pgxscan.NotFound(err) {
			return auth.Session{}, auth.ErrSessionNotFound
		}
		log.Error("failed to get user session", "error", err.Error())
		return auth.Session{}, err
	}

	return sess, nil
}

// Implements auth.Repository
func (r *authRepository) DeleteSession(ctx context.Context, field string, val any) error {
	query := strs.Concatenate("DELETE FROM sessions WHERE ", field, "=$1")
	conn := r.db.GetConn(ctx)
	_, err := conn.Exec(ctx, query, val)
	if err != nil {
		log.Error("failed to delete session", "error", err.Error())
		return err
	}
	return nil
}

// Implements auth.Repository
func (r *authRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := "DELETE FROM sessions WHERE expires_at < NOW()"
	conn := r.db.GetConn(ctx)

	_, err := conn.Exec(ctx, query)
	if err != nil {
		log.Error("failed to delete expired sessions", "error", err.Error())
		return err
	}

	return nil
}

// GetResetPasswordTokenObj implements auth.Repository.
func (r *authRepository) GetResetPasswordTokenObj(
	ctx context.Context,
	tokenValue string,
) (auth.ResetPasswordToken, error) {
	query := "SELECT user_id, value, expires_at, used_at FROM reset_password_tokens WHERE value=$1"

	conn := r.db.GetConn(ctx)
	if r.db.IsTxConn(conn) {
		query += "\nFOR UPDATE"
	}

	var tokenObj auth.ResetPasswordToken
	if err := pgxscan.Get(ctx, conn, &tokenObj, query, tokenValue); err != nil {
		if pgxscan.NotFound(err) {
			return tokenObj, auth.ErrResetPasswordTokenNotFound
		}
		log.Error("failed to get reset password token", "error", err.Error())
		return tokenObj, err
	}

	return tokenObj, nil
}

// InsertResetPasswordToken implements auth.Repository.
func (r *authRepository) InsertResetPasswordToken(
	ctx context.Context,
	token auth.ResetPasswordToken,
) error {
	query := "INSERT INTO reset_password_tokens(user_id, value, expires_at) VALUES($1, $2, $3)"
	conn := r.db.GetConn(ctx)

	if _, err := conn.Exec(ctx, query, token.UserId, token.Value, token.ExpiresAt); err != nil {
		log.Error("failed to insert reset password token", "error", err.Error())
		return err
	}

	return nil
}

// UseResetPasswordToken implements auth.Repository.
func (r *authRepository) UseResetPasswordToken(
	ctx context.Context,
	tokenObj auth.ResetPasswordToken,
) error {
	query := "UPDATE reset_password_tokens SET used_at=$1 WHERE value=$2"
	conn := r.db.GetConn(ctx)

	now := time.Now()

	if _, err := conn.Exec(ctx, query, now, tokenObj.Value); err != nil {
		log.Error("failed to set used_at reset password token", "error", err.Error())
		return err
	}

	return nil
}
