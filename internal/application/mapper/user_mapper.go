package mapper

import (
	"github.com/suproxy/backend/internal/application/dto"
	"github.com/suproxy/backend/internal/domain/user"
)

func ToUserResponse(u *user.User) *dto.UserResponse {
	if u == nil {
		return nil
	}

	return &dto.UserResponse{
		ID:        u.ID,
		Email:     u.Email.String(),
		FirstName: u.Profile.FirstName,
		LastName:  u.Profile.LastName,
		Phone:     u.Profile.Phone,
		Avatar:    u.Profile.Avatar,
		Status:    string(u.Status),
		Role:      string(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToUserListResponse(users []*user.User, total int64, offset, limit int) *dto.UserListResponse {
	responses := make([]*dto.UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, ToUserResponse(u))
	}

	return &dto.UserListResponse{
		Users:  responses,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}
}
