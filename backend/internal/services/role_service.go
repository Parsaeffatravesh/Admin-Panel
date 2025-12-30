package services

import (
	"context"
	"errors"

	"admin-panel/internal/models"
	"admin-panel/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrRoleNameExists = errors.New("role name already exists")
	ErrSystemRole     = errors.New("cannot modify system role")
)

type RoleService struct {
	roleRepo  *repository.RoleRepository
	auditRepo *repository.AuditLogRepository
}

func NewRoleService(
	roleRepo *repository.RoleRepository,
	auditRepo *repository.AuditLogRepository,
) *RoleService {
	return &RoleService{
		roleRepo:  roleRepo,
		auditRepo: auditRepo,
	}
}

type CreateRoleRequest struct {
	Name          string   `json:"name" validate:"required,min=1,max=100"`
	Description   string   `json:"description" validate:"max=500"`
	PermissionIDs []string `json:"permission_ids" validate:"dive,uuid"`
}

type UpdateRoleRequest struct {
	Name          *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description   *string  `json:"description,omitempty" validate:"omitempty,max=500"`
	PermissionIDs []string `json:"permission_ids,omitempty" validate:"omitempty,dive,uuid"`
}

func (s *RoleService) Create(ctx context.Context, req *CreateRoleRequest, tenantID uuid.UUID) (*models.Role, error) {
	existing, _ := s.roleRepo.GetByName(ctx, tenantID, req.Name)
	if existing != nil {
		return nil, ErrRoleNameExists
	}

	role := &models.Role{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	for _, permIDStr := range req.PermissionIDs {
		permID, err := uuid.Parse(permIDStr)
		if err != nil {
			continue
		}
		s.roleRepo.AssignPermissionToRole(ctx, role.ID, permID)
	}

	return role, nil
}

func (s *RoleService) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}

func (s *RoleService) List(ctx context.Context, params *models.ListParams) (*models.PaginatedResponse, error) {
	roles, total, err := s.roleRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / params.PerPage
	if int(total)%params.PerPage > 0 {
		totalPages++
	}

	return &models.PaginatedResponse{
		Data:       roles,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (s *RoleService) Update(ctx context.Context, id uuid.UUID, req *UpdateRoleRequest) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if role.IsSystem {
		return nil, ErrSystemRole
	}

	if req.Name != nil && *req.Name != role.Name {
		existing, _ := s.roleRepo.GetByName(ctx, role.TenantID, *req.Name)
		if existing != nil {
			return nil, ErrRoleNameExists
		}
		role.Name = *req.Name
	}

	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	if req.PermissionIDs != nil {
		s.roleRepo.RemoveAllPermissionsFromRole(ctx, id)
		for _, permIDStr := range req.PermissionIDs {
			permID, err := uuid.Parse(permIDStr)
			if err != nil {
				continue
			}
			s.roleRepo.AssignPermissionToRole(ctx, id, permID)
		}
	}

	return role, nil
}

func (s *RoleService) Delete(ctx context.Context, id uuid.UUID) error {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if role.IsSystem {
		return ErrSystemRole
	}

	return s.roleRepo.Delete(ctx, id)
}

func (s *RoleService) GetAllPermissions(ctx context.Context) ([]*models.Permission, error) {
	return s.roleRepo.GetAllPermissions(ctx)
}

func (s *RoleService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*models.Permission, error) {
	return s.roleRepo.GetRolePermissions(ctx, roleID)
}
