package postgres

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/golang-restapi/internal/auth"
	"github.com/prawirdani/golang-restapi/pkg/logging"
	strs "github.com/prawirdani/golang-restapi/pkg/strings"
)

type authRepository struct {
	db     *db
	logger logging.Logger
}

func NewAuthRepository(
	pool *pgxpool.Pool,
	logger logging.Logger,
) *authRepository {
	return &authRepository{
		db:     &db{pool: pool},
		logger: logger,
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
		r.logger.Error(logging.Postgres, "AuthRepository.InsertSession", err.Error())
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
		r.logger.Error(logging.Postgres, "AuthRepository.SelectSession", err.Error())
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
		r.logger.Error(logging.Postgres, "AuthRepository.DeleteSession", err.Error())
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
		r.logger.Error(
			logging.Postgres,
			"AuthRepository.DeleteExpiredSessions",
			err.Error(),
		)
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
		r.logger.Error(logging.Postgres, "AuthRepository.GetResetPasswordTokenObj", err.Error())
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
		r.logger.Error(logging.Postgres, "AuthRepository.InsertResetPasswordToken", err.Error())
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
		r.logger.Error(logging.Postgres, "AuthRepository.UseResetPasswordToken", err.Error())
		return err
	}

	return nil
}
