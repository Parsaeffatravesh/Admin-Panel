package repository

import (
	"context"
	"time"

	"admin-panel/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeatureFlagRepository struct {
	pool *pgxpool.Pool
}

func NewFeatureFlagRepository(pool *pgxpool.Pool) *FeatureFlagRepository {
	return &FeatureFlagRepository{pool: pool}
}

func (r *FeatureFlagRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*models.FeatureFlag, error) {
	query := `
		SELECT id, tenant_id, key, name, description, enabled, metadata, created_at, updated_at
		FROM feature_flags
		WHERE tenant_id = $1
		ORDER BY key ASC
	`

	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flags []*models.FeatureFlag
	for rows.Next() {
		var flag models.FeatureFlag
		err := rows.Scan(
			&flag.ID,
			&flag.TenantID,
			&flag.Key,
			&flag.Name,
			&flag.Description,
			&flag.Enabled,
			&flag.Metadata,
			&flag.CreatedAt,
			&flag.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		flags = append(flags, &flag)
	}

	return flags, nil
}

func (r *FeatureFlagRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.FeatureFlag, error) {
	query := `
		SELECT id, tenant_id, key, name, description, enabled, metadata, created_at, updated_at
		FROM feature_flags
		WHERE id = $1
	`

	var flag models.FeatureFlag
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&flag.ID,
		&flag.TenantID,
		&flag.Key,
		&flag.Name,
		&flag.Description,
		&flag.Enabled,
		&flag.Metadata,
		&flag.CreatedAt,
		&flag.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &flag, nil
}

func (r *FeatureFlagRepository) GetByKey(ctx context.Context, tenantID uuid.UUID, key string) (*models.FeatureFlag, error) {
	query := `
		SELECT id, tenant_id, key, name, description, enabled, metadata, created_at, updated_at
		FROM feature_flags
		WHERE tenant_id = $1 AND key = $2
	`

	var flag models.FeatureFlag
	err := r.pool.QueryRow(ctx, query, tenantID, key).Scan(
		&flag.ID,
		&flag.TenantID,
		&flag.Key,
		&flag.Name,
		&flag.Description,
		&flag.Enabled,
		&flag.Metadata,
		&flag.CreatedAt,
		&flag.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &flag, nil
}

func (r *FeatureFlagRepository) Create(ctx context.Context, flag *models.FeatureFlag) error {
	now := time.Now()
	flag.ID = uuid.New()
	flag.CreatedAt = now
	flag.UpdatedAt = now

	query := `
		INSERT INTO feature_flags (id, tenant_id, key, name, description, enabled, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.pool.Exec(ctx, query,
		flag.ID,
		flag.TenantID,
		flag.Key,
		flag.Name,
		flag.Description,
		flag.Enabled,
		flag.Metadata,
		flag.CreatedAt,
		flag.UpdatedAt,
	)
	return err
}

func (r *FeatureFlagRepository) Update(ctx context.Context, flag *models.FeatureFlag) error {
	flag.UpdatedAt = time.Now()

	query := `
		UPDATE feature_flags
		SET name = $2, description = $3, enabled = $4, metadata = $5, updated_at = $6
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		flag.ID,
		flag.Name,
		flag.Description,
		flag.Enabled,
		flag.Metadata,
		flag.UpdatedAt,
	)
	return err
}

func (r *FeatureFlagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM feature_flags WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *FeatureFlagRepository) IsEnabled(ctx context.Context, tenantID uuid.UUID, key string) bool {
	query := `SELECT enabled FROM feature_flags WHERE tenant_id = $1 AND key = $2`
	var enabled bool
	err := r.pool.QueryRow(ctx, query, tenantID, key).Scan(&enabled)
	if err != nil {
		return false
	}
	return enabled
}
