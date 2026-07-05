package runtime

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// ProcessInfo represents information about a running process
type ProcessInfo struct {
	InstanceID uuid.UUID
	ProcessID  int
	StartedAt  time.Time
	ConfigPath string
	LogPath    string
	ErrorPath  string
	Command    string
	Args       []string
}

// Registry manages running Xray processes in memory
type Registry interface {
	// Register registers a new process
	Register(info *ProcessInfo) error

	// Remove removes a process from registry
	Remove(instanceID uuid.UUID) error

	// Find finds a process by instance ID
	Find(instanceID uuid.UUID) (*ProcessInfo, bool)

	// List lists all registered processes
	List() []*ProcessInfo

	// RunningCount returns the number of running processes
	RunningCount() int

	// IsRegistered checks if an instance is registered
	IsRegistered(instanceID uuid.UUID) bool

	// Update updates process information
	Update(instanceID uuid.UUID, updater func(*ProcessInfo)) error

	// Clear clears all processes from registry
	Clear()
}

// processRegistry implements Registry with thread-safe operations
type processRegistry struct {
	mu        sync.RWMutex
	processes map[uuid.UUID]*ProcessInfo
}

// NewRegistry creates a new process registry
func NewRegistry() Registry {
	return &processRegistry{
		processes: make(map[uuid.UUID]*ProcessInfo),
	}
}

func (r *processRegistry) Register(info *ProcessInfo) error {
	if info == nil {
		return ErrInvalidProcessInfo
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if already registered
	if _, exists := r.processes[info.InstanceID]; exists {
		return ErrProcessAlreadyRegistered
	}

	r.processes[info.InstanceID] = info
	return nil
}

func (r *processRegistry) Remove(instanceID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.processes[instanceID]; !exists {
		return ErrProcessNotFound
	}

	delete(r.processes, instanceID)
	return nil
}

func (r *processRegistry) Find(instanceID uuid.UUID) (*ProcessInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, exists := r.processes[instanceID]
	if !exists {
		return nil, false
	}

	// Return a copy to prevent external modification
	infoCopy := *info
	return &infoCopy, true
}

func (r *processRegistry) List() []*ProcessInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*ProcessInfo, 0, len(r.processes))
	for _, info := range r.processes {
		// Create copies to prevent external modification
		infoCopy := *info
		result = append(result, &infoCopy)
	}
	return result
}

func (r *processRegistry) RunningCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.processes)
}

func (r *processRegistry) IsRegistered(instanceID uuid.UUID) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.processes[instanceID]
	return exists
}

func (r *processRegistry) Update(instanceID uuid.UUID, updater func(*ProcessInfo)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, exists := r.processes[instanceID]
	if !exists {
		return ErrProcessNotFound
	}

	updater(info)
	return nil
}

func (r *processRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.processes = make(map[uuid.UUID]*ProcessInfo)
}
