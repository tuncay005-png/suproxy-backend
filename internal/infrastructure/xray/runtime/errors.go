package runtime

import "errors"

var (
	ErrProcessNotFound          = errors.New("process not found in registry")
	ErrProcessAlreadyRegistered = errors.New("process already registered")
	ErrInvalidProcessInfo       = errors.New("invalid process info")
	ErrProcessNotRunning        = errors.New("process is not running")
	ErrProcessAlreadyRunning    = errors.New("process is already running")
	ErrProcessStartFailed       = errors.New("failed to start process")
	ErrProcessStopFailed        = errors.New("failed to stop process")
	ErrProcessKillFailed        = errors.New("failed to kill process")
	ErrConfigFileNotFound       = errors.New("config file not found")
	ErrBinaryNotFound           = errors.New("xray binary not found")
	ErrInvalidPID               = errors.New("invalid process ID")
)
