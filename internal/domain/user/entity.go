package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID
	Email             Email
	Password          Password
	Profile           Profile
	Status            Status
	Role              Role
	LastLoginAt       *time.Time
	LastLoginIP       string
	FailedLoginCount  int
	LockedUntil       *time.Time
	PasswordChangedAt *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewUser(email Email, password Password, profile Profile) (*User, error) {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		Profile:   profile,
		Status:    StatusActive,
		Role:      RoleUser,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (u *User) Activate() error {
	if u.Status == StatusActive {
		return ErrUserAlreadyActive
	}
	u.Status = StatusActive
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) Deactivate() error {
	if u.Status == StatusInactive {
		return ErrUserAlreadyInactive
	}
	u.Status = StatusInactive
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) Suspend() error {
	if u.Status == StatusSuspended {
		return ErrUserAlreadySuspended
	}
	u.Status = StatusSuspended
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func (u *User) UpdateProfile(profile Profile) {
	u.Profile = profile
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) ChangePassword(newPassword Password) {
	u.Password = newPassword
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
