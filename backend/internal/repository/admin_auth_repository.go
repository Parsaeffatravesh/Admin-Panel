package repository

import (
	"context"
	"time"

	"admin-panel/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminAuthRepository struct {
	pool *pgxpool.Pool
}

func NewAdminAuthRepository(pool *pgxpool.Pool) *AdminAuthRepository {
	return &AdminAuthRepository{pool: pool}
}

func (r *AdminAuthRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.AdminAuth, error) {
	query := `
		SELECT user_id, admin_password_hash, is_admin, enabled_at, created_at, updated_at
		FROM admin_auth
		WHERE user_id = $1
	`

	var auth models.AdminAuth
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&auth.UserID,
		&auth.AdminPasswordHash,
		&auth.IsAdmin,
		&auth.EnabledAt,
		&auth.CreatedAt,
		&auth.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *AdminAuthRepository) SetAdmin(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	now := time.Now()
	query := `
		INSERT INTO admin_auth (user_id, admin_password_hash, is_admin, enabled_at, created_at, updated_at)
		VALUES ($1, $2, TRUE, $3, $3, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			admin_password_hash = EXCLUDED.admin_password_hash,
			is_admin = TRUE,
			enabled_at = COALESCE(admin_auth.enabled_at, $3),
			updated_at = $3
	`

	_, err := r.pool.Exec(ctx, query, userID, passwordHash, now)
	return err
}

func (r *AdminAuthRepository) UnsetAdmin(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE admin_auth
		SET is_admin = FALSE, updated_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

func (r *AdminAuthRepository) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	query := `
		SELECT is_admin FROM admin_auth WHERE user_id = $1
	`

	var isAdmin bool
	err := r.pool.QueryRow(ctx, query, userID).Scan(&isAdmin)
	if err != nil {
		return false, nil
	}

	return isAdmin, nil
}

func (r *AdminAuthRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM admin_auth WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}
