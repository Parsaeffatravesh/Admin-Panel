package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

type Role struct {
	ID          uuid.UUID   `json:"id"`
	TenantID    uuid.UUID   `json:"tenant_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	IsSystem    bool        `json:"is_system"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Permissions []string    `json:"permissions,omitempty"`
}

type Permission struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserRole struct {
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuditLog struct {
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	UserID     *uuid.UUID `json:"user_id,omitempty"`
	Action     string     `json:"action"`
	Resource   string     `json:"resource"`
	ResourceID *uuid.UUID `json:"resource_id,omitempty"`
	OldValue   *string    `json:"old_value,omitempty"`
	NewValue   *string    `json:"new_value,omitempty"`
	IPAddress  string     `json:"ip_address"`
	UserAgent  string     `json:"user_agent"`
	CreatedAt  time.Time  `json:"created_at"`
}

type Setting struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	RefreshToken string    `json:"-"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

type ListParams struct {
	Page     int
	PerPage  int
	Sort     string
	Order    string
	Search   string
	Filters  map[string]interface{}
	TenantID uuid.UUID
}
