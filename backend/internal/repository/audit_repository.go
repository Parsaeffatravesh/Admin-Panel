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

type AuditLogRepository struct {
	db *pgxpool.Pool
}

func NewAuditLogRepository(db *pgxpool.Pool) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Log(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, tenant_id, user_id, action, resource, resource_id, old_value, new_value, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		log.ID, log.TenantID, log.UserID, log.Action, log.Resource,
		log.ResourceID, log.OldValue, log.NewValue, log.IPAddress, log.UserAgent, log.CreatedAt,
	)
	return err
}

func (r *AuditLogRepository) List(ctx context.Context, params *models.ListParams) ([]*models.AuditLog, int64, error) {
	var conditions []string
	var args []interface{}
	argCount := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argCount))
	args = append(args, params.TenantID)
	argCount++

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(action ILIKE $%d OR resource ILIKE $%d)", argCount, argCount))
		args = append(args, "%"+params.Search+"%")
		argCount++
	}

	if action, ok := params.Filters["action"].(string); ok && action != "" {
		conditions = append(conditions, fmt.Sprintf("action = $%d", argCount))
		args = append(args, action)
		argCount++
	}

	if resource, ok := params.Filters["resource"].(string); ok && resource != "" {
		conditions = append(conditions, fmt.Sprintf("resource = $%d", argCount))
		args = append(args, resource)
		argCount++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_logs %s", whereClause)
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PerPage

	query := fmt.Sprintf(`
		SELECT id, tenant_id, user_id, action, resource, resource_id, old_value, new_value, ip_address, user_agent, created_at
		FROM audit_logs %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	args = append(args, params.PerPage, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		err := rows.Scan(
			&log.ID, &log.TenantID, &log.UserID, &log.Action, &log.Resource,
			&log.ResourceID, &log.OldValue, &log.NewValue, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

func (r *AuditLogRepository) CountRecentLogins(ctx context.Context, tenantID uuid.UUID, duration time.Duration) (int64, error) {
	query := `SELECT COUNT(*) FROM audit_logs WHERE tenant_id = $1 AND action = 'login' AND created_at > $2`
	var count int64
	err := r.db.QueryRow(ctx, query, tenantID, time.Now().Add(-duration)).Scan(&count)
	return count, err
}

func (r *AuditLogRepository) GetRecentActivity(ctx context.Context, tenantID uuid.UUID, limit int) ([]*models.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, resource, resource_id, old_value, new_value, ip_address, user_agent, created_at
		FROM audit_logs WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(ctx, query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		err := rows.Scan(
			&log.ID, &log.TenantID, &log.UserID, &log.Action, &log.Resource,
			&log.ResourceID, &log.OldValue, &log.NewValue, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
