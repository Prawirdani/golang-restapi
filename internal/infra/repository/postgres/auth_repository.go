package postgres

import (
	"context"
	"errors"

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

// StoreSession implements [auth.Repository]
func (r *authRepository) StoreSession(ctx context.Context, session *auth.Session) error {
	if session == nil {
		log.WarnCtx(ctx, "StoreSession called with nil session data")
		return errors.New("session is nil")
	}

	query := "INSERT INTO sessions(id, user_id, user_agent, expires_at, accessed_at) VALUES($1, $2, $3, $4, $5)"
	conn := r.db.GetConn(ctx)
	if _, err := conn.Exec(ctx, query, session.ID, session.UserID, session.UserAgent, session.ExpiresAt, session.AccessedAt); err != nil {
		log.ErrorCtx(ctx, "Failed to store session", err)
		return err
	}

	return nil
}

// GetSession implements [auth.Repository]
func (r *authRepository) GetSession(
	ctx context.Context,
	sessID string,
) (*auth.Session, error) {
	query := strs.Concatenate(
		"UPDATE sessions SET accessed_at=NOW() WHERE id",
		"=$1 RETURNING *",
	)

	conn := r.db.GetConn(ctx)
	if r.db.IsTxConn(conn) {
		query += "\nFOR UPDATE"
	}

	var sess auth.Session
	if err := pgxscan.Get(ctx, conn, &sess, query, sessID); err != nil {
		if noRowsErr(err) {
			return nil, auth.ErrSessionNotFound
		}
		log.ErrorCtx(ctx, "Failed to get session data", err)
		return nil, err
	}

	return &sess, nil
}

// UpdateSession implements [auth.Repository]
func (r *authRepository) UpdateSession(ctx context.Context, session *auth.Session) error {
	if session == nil {
		log.WarnCtx(ctx, "UpdateSession called with nil session struct ptr")
		return errors.New("session is nil")
	}

	query := "UPDATE sessions SET expires_at=$1 WHERE id=$2"
	conn := r.db.GetConn(ctx)

	if _, err := conn.Exec(ctx, query, session.ExpiresAt, session.ID); err != nil {
		log.ErrorCtx(ctx, "Failed to updated session", err)
		return err
	}

	return nil
}

// GetResetPasswordToken implements [auth.Repository]
func (r *authRepository) GetResetPasswordToken(
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

// StoreResetPasswordToken implements [auth.Repository]
func (r *authRepository) StoreResetPasswordToken(
	ctx context.Context,
	token *auth.ResetPasswordToken,
) error {
	if token == nil {
		log.WarnCtx(ctx, "StoreResetPasswordToken called with nil token ptr")
		return errors.New("reset password token is nil")
	}

	query := "INSERT INTO reset_password_tokens(user_id, value, expires_at) VALUES($1, $2, $3)"
	conn := r.db.GetConn(ctx)

	if _, err := conn.Exec(ctx, query, token.UserID, token.Value, token.ExpiresAt); err != nil {
		log.ErrorCtx(ctx, "Failed to store reset password token", err)
		return err
	}

	return nil
}

// UpdateResetPasswordToken implements [auth.Repository]
func (r *authRepository) UpdateResetPasswordToken(
	ctx context.Context,
	token *auth.ResetPasswordToken,
) error {
	if token == nil {
		log.WarnCtx(ctx, "UpdateResetPasswordToken called with nil token object")
		return errors.New("reset password token is nil")
	}

	query := "UPDATE reset_password_tokens SET used_at=$1 WHERE value=$2"
	conn := r.db.GetConn(ctx)

	if _, err := conn.Exec(ctx, query, token.UsedAt, token.Value); err != nil {
		log.ErrorCtx(ctx, "Failed to update reset password token", err)
		return err
	}

	return nil
}
