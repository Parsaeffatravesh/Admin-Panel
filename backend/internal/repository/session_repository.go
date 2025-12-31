package repository

import (
	"context"
	"database/sql"
	"time"

	"admin-panel/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, refresh_token, ip_address, user_agent, expires_at, created_at, rotated_at, replaced_by_token, revoked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID, session.UserID, session.RefreshToken,
		session.IPAddress, session.UserAgent, session.ExpiresAt, session.CreatedAt,
		session.RotatedAt, session.ReplacedByToken, session.RevokedAt,
	)
	return err
}

func (r *SessionRepository) GetByRefreshToken(ctx context.Context, token string) (*models.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, ip_address, user_agent, expires_at, created_at, rotated_at, replaced_by_token, revoked_at
		FROM sessions WHERE refresh_token = $1
	`
	session := &models.Session{}
	var rotatedAt sql.NullTime
	var replacedByToken sql.NullString
	var revokedAt sql.NullTime
	err := r.db.QueryRow(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.RefreshToken,
		&session.IPAddress, &session.UserAgent, &session.ExpiresAt, &session.CreatedAt,
		&rotatedAt, &replacedByToken, &revokedAt,
	)
	if err != nil {
		return nil, err
	}
	if rotatedAt.Valid {
		session.RotatedAt = &rotatedAt.Time
	}
	if replacedByToken.Valid {
		session.ReplacedByToken = &replacedByToken.String
	}
	if revokedAt.Valid {
		session.RevokedAt = &revokedAt.Time
	}
	return session, nil
}

func (r *SessionRepository) MarkRotated(ctx context.Context, sessionID uuid.UUID, rotatedAt time.Time, replacedByToken string) error {
	query := `UPDATE sessions SET rotated_at = $2, replaced_by_token = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, sessionID, rotatedAt, replacedByToken)
	return err
}

func (r *SessionRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error {
	query := `UPDATE sessions SET revoked_at = $2 WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, userID, revokedAt)
	return err
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.db.Exec(ctx, query)
	return err
}
