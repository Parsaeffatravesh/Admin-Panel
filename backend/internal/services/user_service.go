package services

import (
	"context"
	"errors"

	"admin-panel/internal/models"
	"admin-panel/internal/repository"
	"admin-panel/internal/utils"

	"github.com/google/uuid"
)

var (
	ErrEmailExists = errors.New("email already exists")
)

type UserService struct {
	userRepo  *repository.UserRepository
	roleRepo  *repository.RoleRepository
	auditRepo *repository.AuditLogRepository
}

func NewUserService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	auditRepo *repository.AuditLogRepository,
) *UserService {
	return &UserService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		auditRepo: auditRepo,
	}
}

type CreateUserRequest struct {
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8"`
	FirstName string    `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string    `json:"last_name" validate:"required,min=1,max=100"`
	RoleIDs   []string  `json:"role_ids" validate:"dive,uuid"`
}

type UpdateUserRequest struct {
	Email     *string  `json:"email,omitempty" validate:"omitempty,email"`
	FirstName *string  `json:"first_name,omitempty" validate:"omitempty,min=1,max=100"`
	LastName  *string  `json:"last_name,omitempty" validate:"omitempty,min=1,max=100"`
	Status    *string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended"`
	RoleIDs   []string `json:"role_ids,omitempty" validate:"omitempty,dive,uuid"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (s *UserService) Create(ctx context.Context, req *CreateUserRequest, tenantID uuid.UUID) (*models.User, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailExists
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Email:        req.Email,
		PasswordHash: passwordHash,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Status:       "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	for _, roleIDStr := range req.RoleIDs {
		roleID, err := uuid.Parse(roleIDStr)
		if err != nil {
			continue
		}
		s.roleRepo.AssignRoleToUser(ctx, user.ID, roleID)
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context, params *models.ListParams) (*models.PaginatedResponse, error) {
	users, total, err := s.userRepo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / params.PerPage
	if int(total)%params.PerPage > 0 {
		totalPages++
	}

	return &models.PaginatedResponse{
		Data:       users,
		Total:      total,
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Email != nil && *req.Email != user.Email {
		existing, _ := s.userRepo.GetByEmail(ctx, *req.Email)
		if existing != nil {
			return nil, ErrEmailExists
		}
		user.Email = *req.Email
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	if req.RoleIDs != nil {
		s.roleRepo.RemoveAllRolesFromUser(ctx, id)
		for _, roleIDStr := range req.RoleIDs {
			roleID, err := uuid.Parse(roleIDStr)
			if err != nil {
				continue
			}
			s.roleRepo.AssignRoleToUser(ctx, id, roleID)
		}
	}

	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ResetPassword(ctx context.Context, id uuid.UUID, req *ResetPasswordRequest) error {
	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, id, passwordHash)
}

func (s *UserService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*models.Role, error) {
	return s.roleRepo.GetUserRoles(ctx, userID)
}
