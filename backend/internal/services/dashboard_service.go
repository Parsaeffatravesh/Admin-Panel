package services

import (
	"context"
	"time"

	"admin-panel/internal/repository"

	"github.com/google/uuid"
)

type DashboardService struct {
	userRepo  *repository.UserRepository
	roleRepo  *repository.RoleRepository
	auditRepo *repository.AuditLogRepository
}

func NewDashboardService(
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
	auditRepo *repository.AuditLogRepository,
) *DashboardService {
	return &DashboardService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		auditRepo: auditRepo,
	}
}

type DashboardStats struct {
	TotalUsers      int64           `json:"total_users"`
	ActiveUsers     int64           `json:"active_users"`
	TotalRoles      int64           `json:"total_roles"`
	RecentLogins    int64           `json:"recent_logins"`
	UsersByStatus   map[string]int64 `json:"users_by_status"`
	RecentActivity  []ActivityItem  `json:"recent_activity"`
}

type ActivityItem struct {
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	UserEmail string    `json:"user_email"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *DashboardService) GetStats(ctx context.Context, tenantID uuid.UUID) (*DashboardStats, error) {
	totalUsers, err := s.userRepo.Count(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	activeUsers, err := s.userRepo.CountByStatus(ctx, tenantID, "active")
	if err != nil {
		return nil, err
	}

	totalRoles, err := s.roleRepo.Count(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	recentLogins, err := s.auditRepo.CountRecentLogins(ctx, tenantID, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	usersByStatus, err := s.userRepo.CountGroupByStatus(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	recentActivity, err := s.auditRepo.GetRecentActivity(ctx, tenantID, 10)
	if err != nil {
		return nil, err
	}

	activityItems := make([]ActivityItem, len(recentActivity))
	for i, a := range recentActivity {
		activityItems[i] = ActivityItem{
			Action:    a.Action,
			Resource:  a.Resource,
			UserEmail: "",
			CreatedAt: a.CreatedAt,
		}
	}

	return &DashboardStats{
		TotalUsers:     totalUsers,
		ActiveUsers:    activeUsers,
		TotalRoles:     totalRoles,
		RecentLogins:   recentLogins,
		UsersByStatus:  usersByStatus,
		RecentActivity: activityItems,
	}, nil
}
