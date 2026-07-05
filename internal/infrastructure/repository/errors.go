package repository

import "errors"

var (
	ErrNotFound       = errors.New("entity not found")
	ErrDuplicate      = errors.New("entity already exists")
	ErrInvalidInput   = errors.New("invalid input")
	ErrDatabase       = errors.New("database error")
	ErrTransaction    = errors.New("transaction error")
	ErrNoRowsAffected = errors.New("no rows affected")
)
