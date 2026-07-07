package dto

import (
	"time"

	"github.com/google/uuid"
)

// Admin User List Request with filters
type AdminUserListRequest struct {
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit" binding:"max=100"`
	Role      string `form:"role" binding:"omitempty,oneof=user admin"`
	Status    string `form:"status" binding:"omitempty,oneof=active inactive suspended"`
	Email     string `form:"email"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=created_at email status"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// Admin User Response with additional fields
type AdminUserResponse struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	FirstName         string     `json:"first_name"`
	LastName          string     `json:"last_name"`
	Phone             string     `json:"phone"`
	Avatar            string     `json:"avatar"`
	Status            string     `json:"status"`
	Role              string     `json:"role"`
	LastLoginAt       *time.Time `json:"last_login_at"`
	LastLoginIP       string     `json:"last_login_ip"`
	FailedLoginCount  int        `json:"failed_login_count"`
	LockedUntil       *time.Time `json:"locked_until"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// Admin User List Response
type AdminUserListResponse struct {
	Users  []*AdminUserResponse `json:"users"`
	Total  int64                `json:"total"`
	Offset int                  `json:"offset"`
	Limit  int                  `json:"limit"`
}

// Update User Status Request
type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive suspended"`
}

// Update User Role Request
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user admin"`
}
