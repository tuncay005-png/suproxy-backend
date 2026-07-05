package user

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserAlreadyActive     = errors.New("user already active")
	ErrUserAlreadyInactive   = errors.New("user already inactive")
	ErrUserAlreadySuspended  = errors.New("user already suspended")
	ErrInvalidEmail          = errors.New("invalid email")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserNotActive         = errors.New("user not active")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrInsufficientPermission = errors.New("insufficient permission")
)
