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

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	query := `
		INSERT INTO roles (id, tenant_id, name, description, is_system, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	_, err := r.db.Exec(ctx, query,
		role.ID, role.TenantID, role.Name, role.Description, role.IsSystem, role.CreatedAt, role.UpdatedAt,
	)
	return err
}

func (r *RoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	query := `
		SELECT id, tenant_id, name, description, is_system, created_at, updated_at
		FROM roles WHERE id = $1
	`
	role := &models.Role{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID, &role.TenantID, &role.Name, &role.Description,
		&role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*models.Role, error) {
	query := `
		SELECT id, tenant_id, name, description, is_system, created_at, updated_at
		FROM roles WHERE tenant_id = $1 AND name = $2
	`
	role := &models.Role{}
	err := r.db.QueryRow(ctx, query, tenantID, name).Scan(
		&role.ID, &role.TenantID, &role.Name, &role.Description,
		&role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) List(ctx context.Context, params *models.ListParams) ([]*models.Role, int64, error) {
	var conditions []string
	var args []interface{}
	argCount := 1

	conditions = append(conditions, fmt.Sprintf("tenant_id = $%d", argCount))
	args = append(args, params.TenantID)
	argCount++

	if params.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argCount, argCount))
		args = append(args, "%"+params.Search+"%")
		argCount++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM roles %s", whereClause)
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PerPage

	query := fmt.Sprintf(`
		SELECT id, tenant_id, name, description, is_system, created_at, updated_at
		FROM roles %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	args = append(args, params.PerPage, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(
			&role.ID, &role.TenantID, &role.Name, &role.Description,
			&role.IsSystem, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		roles = append(roles, role)
	}

	return roles, total, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *models.Role) error {
	query := `UPDATE roles SET name = $2, description = $3, updated_at = $4 WHERE id = $1`
	role.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query, role.ID, role.Name, role.Description, role.UpdatedAt)
	return err
}

func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM role_permissions WHERE role_id = $1", id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, "DELETE FROM user_roles WHERE role_id = $1", id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, "DELETE FROM roles WHERE id = $1", id)
	return err
}

func (r *RoleRepository) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM roles WHERE tenant_id = $1", tenantID).Scan(&count)
	return count, err
}

func (r *RoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `INSERT INTO user_roles (user_id, role_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, userID, roleID, time.Now())
	return err
}

func (r *RoleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := r.db.Exec(ctx, query, userID, roleID)
	return err
}

func (r *RoleRepository) RemoveAllRolesFromUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *RoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*models.Role, error) {
	query := `
		SELECT r.id, r.tenant_id, r.name, r.description, r.is_system, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		role := &models.Role{}
		err := rows.Scan(&role.ID, &role.TenantID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `INSERT INTO role_permissions (role_id, permission_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, roleID, permissionID, time.Now())
	return err
}

func (r *RoleRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
	_, err := r.db.Exec(ctx, query, roleID, permissionID)
	return err
}

func (r *RoleRepository) RemoveAllPermissionsFromRole(ctx context.Context, roleID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1`
	_, err := r.db.Exec(ctx, query, roleID)
	return err
}

func (r *RoleRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*models.Permission, error) {
	query := `
		SELECT p.id, p.name, p.resource, p.action, p.description, p.created_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`
	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*models.Permission
	for rows.Next() {
		perm := &models.Permission{}
		err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description, &perm.CreatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

func (r *RoleRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]*models.Permission, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.resource, p.action, p.description, p.created_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*models.Permission
	for rows.Next() {
		perm := &models.Permission{}
		err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description, &perm.CreatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

func (r *RoleRepository) GetAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	query := `SELECT id, name, resource, action, description, created_at FROM permissions ORDER BY resource, action`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*models.Permission
	for rows.Next() {
		perm := &models.Permission{}
		err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description, &perm.CreatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}
	return permissions, nil
}
