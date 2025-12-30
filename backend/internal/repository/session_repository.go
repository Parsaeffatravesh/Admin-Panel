package repository

import (
	"context"

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
		INSERT INTO sessions (id, user_id, refresh_token, ip_address, user_agent, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		session.ID, session.UserID, session.RefreshToken,
		session.IPAddress, session.UserAgent, session.ExpiresAt, session.CreatedAt,
	)
	return err
}

func (r *SessionRepository) GetByRefreshToken(ctx context.Context, token string) (*models.Session, error) {
	query := `
		SELECT id, user_id, refresh_token, ip_address, user_agent, expires_at, created_at
		FROM sessions WHERE refresh_token = $1
	`
	session := &models.Session{}
	err := r.db.QueryRow(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.RefreshToken,
		&session.IPAddress, &session.UserAgent, &session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) Update(ctx context.Context, session *models.Session) error {
	query := `UPDATE sessions SET refresh_token = $2, expires_at = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, session.ID, session.RefreshToken, session.ExpiresAt)
	return err
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.db.Exec(ctx, query)
	return err
}
