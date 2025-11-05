package postgres

import (
	"context"
	"errors"
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
func (r *authRepository) InsertSession(ctx context.Context, sess *auth.Session) error {
	if sess == nil {
		log.WarnCtx(ctx, "Insert session called with nil session data")
		return errors.New("sess is nil")
	}

	log.DebugCtx(ctx, "Inserting session data", "args", sess)

	query := "INSERT INTO sessions(user_id, refresh_token, user_agent, expires_at, accessed_at) VALUES($1, $2, $3, $4, $5)"
	conn := r.db.GetConn(ctx)
	if _, err := conn.Exec(ctx, query, sess.UserID, sess.RefreshToken, sess.UserAgent, sess.ExpiresAt, sess.AccessedAt); err != nil {
		log.ErrorCtx(ctx, "Failed to insert session", err)
		return err
	}

	return nil
}

// Implements auth.Repository
func (r *authRepository) GetUserSessionBy(
	ctx context.Context,
	field string,
	value any,
) (*auth.Session, error) {
	log.DebugCtx(ctx, "Getting session data", "search_field", field, "search_arg", value)

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
	if err := pgxscan.Get(ctx, conn, &sess, query, value); err != nil {
		if noRowsErr(err) {
			return nil, auth.ErrSessionNotFound
		}
		log.ErrorCtx(ctx, "Failed to get user session", err)
		return nil, err
	}

	return &sess, nil
}

// Implements auth.Repository
func (r *authRepository) DeleteSession(ctx context.Context, field string, value any) error {
	log.DebugCtx(ctx, "Deleting session data", "search_field", field, "search_arg", value)

	query := strs.Concatenate("DELETE FROM sessions WHERE ", field, "=$1")
	conn := r.db.GetConn(ctx)
	_, err := conn.Exec(ctx, query, value)
	if err != nil {
		log.ErrorCtx(ctx, "Failed to delete session", err)
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
		log.ErrorCtx(ctx, "Failed to delete expired sessions", err)
		return err
	}

	return nil
}

// GetResetPasswordTokenObj implements auth.Repository.
func (r *authRepository) GetResetPasswordTokenObj(
	ctx context.Context,
	tokenValue string,
) (*auth.ResetPasswordToken, error) {
	query := "SELECT user_id, value, expires_at, used_at FROM reset_password_tokens WHERE value=$1"

	conn := r.db.GetConn(ctx)
	if r.db.IsTxConn(conn) {
		query += "\nFOR UPDATE"
	}

	var tokenObj auth.ResetPasswordToken
	if err := pgxscan.Get(ctx, conn, &tokenObj, query, tokenValue); err != nil {
		if noRowsErr(err) {
			return nil, auth.ErrResetPasswordTokenNotFound
		}
		log.ErrorCtx(ctx, "Failed to get reset password token", err)
		return nil, err
	}

	return &tokenObj, nil
}

// InsertResetPasswordToken implements auth.Repository.
func (r *authRepository) InsertResetPasswordToken(
	ctx context.Context,
	token *auth.ResetPasswordToken,
) error {
	if token == nil {
		log.WarnCtx(ctx, "Insert reset password token called with nil token object")
		return errors.New("reset password token is nil")
	}

	log.DebugCtx(ctx, "Inserting reset password token", "args", token)

	query := "INSERT INTO reset_password_tokens(user_id, value, expires_at) VALUES($1, $2, $3)"
	conn := r.db.GetConn(ctx)

	if _, err := conn.Exec(ctx, query, token.UserId, token.Value, token.ExpiresAt); err != nil {
		log.ErrorCtx(ctx, "Failed to insert reset password token", err)
		return err
	}

	return nil
}

// InvalidateResetPasswordToken implements auth.Repository.
func (r *authRepository) InvalidateResetPasswordToken(
	ctx context.Context,
	token *auth.ResetPasswordToken,
) error {
	if token == nil {
		log.WarnCtx(ctx, "Invalidate reset password token called with nil token object")
		return errors.New("reset password token is nil")
	}

	query := "UPDATE reset_password_tokens SET used_at=$1 WHERE value=$2"
	conn := r.db.GetConn(ctx)

	now := time.Now()

	if _, err := conn.Exec(ctx, query, now, token.Value); err != nil {
		log.ErrorCtx(ctx, "Failed to invalidate reset password token", err)
		return err
	}

	return nil
}
