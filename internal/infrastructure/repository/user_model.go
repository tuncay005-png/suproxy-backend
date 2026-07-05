package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/user"
)

type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	FirstName string    `gorm:"type:varchar(100)"`
	LastName  string    `gorm:"type:varchar(100)"`
	Phone     string    `gorm:"type:varchar(20)"`
	Avatar    string    `gorm:"type:varchar(500)"`
	Status    string    `gorm:"type:varchar(20);not null;default:'active'"`
	Role      string    `gorm:"type:varchar(20);not null;default:'user'"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (UserModel) TableName() string {
	return "users"
}

func toUserModel(u *user.User) *UserModel {
	return &UserModel{
		ID:        u.ID,
		Email:     u.Email.String(),
		Password:  u.Password.Hash(),
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

func toDomainUser(m *UserModel) (*user.User, error) {
	email, err := user.NewEmail(m.Email)
	if err != nil {
		return nil, err
	}

	password, err := user.NewPassword(m.Password)
	if err != nil {
		return nil, err
	}

	profile := user.NewProfile(m.FirstName, m.LastName, m.Phone, m.Avatar)

	return &user.User{
		ID:        m.ID,
		Email:     email,
		Password:  password,
		Profile:   profile,
		Status:    user.Status(m.Status),
		Role:      user.Role(m.Role),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
