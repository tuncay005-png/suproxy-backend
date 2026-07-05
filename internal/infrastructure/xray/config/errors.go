package config

import "errors"

var (
	ErrInvalidConfig     = errors.New("invalid config")
	ErrConfigNotFound    = errors.New("config file not found")
	ErrConfigWriteFailed = errors.New("failed to write config")
	ErrConfigReadFailed  = errors.New("failed to read config")
	ErrBackupFailed      = errors.New("failed to create backup")
	ErrRestoreFailed     = errors.New("failed to restore from backup")
	ErrBackupNotFound    = errors.New("backup not found")
	ErrInvalidBackupTime = errors.New("invalid backup time")
)
