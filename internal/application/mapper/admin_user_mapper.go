package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/user"
)

// ToAdminUserResponse converts domain user to admin user response with all fields
func ToAdminUserResponse(u *user.User) *dto.AdminUserResponse {
	if u == nil {
		return nil
	}

	return &dto.AdminUserResponse{
		ID:                u.ID,
		Email:             u.Email.String(),
		FirstName:         u.Profile.FirstName,
		LastName:          u.Profile.LastName,
		Phone:             u.Profile.Phone,
		Avatar:            u.Profile.Avatar,
		Status:            string(u.Status),
		Role:              string(u.Role),
		LastLoginAt:       u.LastLoginAt,
		LastLoginIP:       u.LastLoginIP,
		FailedLoginCount:  u.FailedLoginCount,
		LockedUntil:       u.LockedUntil,
		PasswordChangedAt: u.PasswordChangedAt,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}

// ToAdminUserListResponse converts domain users to admin user list response
func ToAdminUserListResponse(users []*user.User, total int64, offset, limit int) *dto.AdminUserListResponse {
	responses := make([]*dto.AdminUserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, ToAdminUserResponse(u))
	}

	return &dto.AdminUserListResponse{
		Users:  responses,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}
}
