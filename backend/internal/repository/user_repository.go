package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"admin-panel/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.db.Exec(ctx, query,
		user.ID, user.TenantID, user.Email, user.PasswordHash,
		user.FirstName, user.LastName, user.Status, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, status, created_at, updated_at, last_login_at
		FROM users WHERE id = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Status,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, status, created_at, updated_at, last_login_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Status,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context, params *models.ListParams) ([]*models.User, int64, error) {
	var conditions []string
	var args []interface{}
	argCount := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argCount))
	args = append(args, params.TenantID)
	argCount++

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(email ILIKE $%d OR first_name ILIKE $%d OR last_name ILIKE $%d)",
			argCount, argCount, argCount,
		))
		args = append(args, "%"+params.Search+"%")
		argCount++
	}

	if status, ok := params.Filters["status"].(string); ok && status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, status)
		argCount++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	sortColumn := "created_at"
	if params.Sort != "" {
		allowedSorts := map[string]bool{"email": true, "first_name": true, "last_name": true, "status": true, "created_at": true}
		if allowedSorts[params.Sort] {
			sortColumn = params.Sort
		}
	}

	sortOrder := "DESC"
	if params.Order == "asc" {
		sortOrder = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage

	query := fmt.Sprintf(`
		SELECT id, tenant_id, email, password_hash, first_name, last_name, status, created_at, updated_at, last_login_at
		FROM users %s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sortColumn, sortOrder, argCount, argCount+1)

	args = append(args, params.PerPage, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.TenantID, &user.Email, &user.PasswordHash,
			&user.FirstName, &user.LastName, &user.Status,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET email = $2, first_name = $3, last_name = $4, status = $5, updated_at = $6
		WHERE id = $1
	`
	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.FirstName, user.LastName, user.Status, user.UpdatedAt)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, passwordHash, time.Now())
	return err
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, time.Now())
	return err
}

func (r *UserRepository) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE tenant_id = $1", tenantID).Scan(&count)
	return count, err
}

func (r *UserRepository) CountByStatus(ctx context.Context, tenantID uuid.UUID, status string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE tenant_id = $1 AND status = $2", tenantID, status).Scan(&count)
	return count, err
}

func (r *UserRepository) CountGroupByStatus(ctx context.Context, tenantID uuid.UUID) (map[string]int64, error) {
	query := `SELECT status, COUNT(*) FROM users WHERE tenant_id = $1 GROUP BY status`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		result[status] = count
	}
	return result, nil
}
