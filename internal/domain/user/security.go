package user

import (
	"time"
)

const (
	MaxFailedLoginAttempts = 5
	AccountLockDuration    = 15 * time.Minute
)

func (u *User) RecordSuccessfulLogin(ipAddress string) {
	now := time.Now().UTC()
	u.LastLoginAt = &now
	u.LastLoginIP = ipAddress
	u.FailedLoginCount = 0
	u.LockedUntil = nil
	u.UpdatedAt = now
}

func (u *User) RecordFailedLogin() {
	u.FailedLoginCount++
	u.UpdatedAt = time.Now().UTC()

	if u.FailedLoginCount >= MaxFailedLoginAttempts {
		lockUntil := time.Now().UTC().Add(AccountLockDuration)
		u.LockedUntil = &lockUntil
	}
}

func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().UTC().Before(*u.LockedUntil)
}

func (u *User) Unlock() {
	u.LockedUntil = nil
	u.FailedLoginCount = 0
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) RecordPasswordChange() {
	now := time.Now().UTC()
	u.PasswordChangedAt = &now
	u.UpdatedAt = now
}
