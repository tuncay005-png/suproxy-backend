package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/suproxy/backend/internal/domain/user"
	"github.com/suproxy/backend/internal/domain/xray"
	"github.com/suproxy/backend/internal/infrastructure/logger"
	xrayConfig "github.com/suproxy/backend/internal/infrastructure/xray/config"
)

// Mock repositories and dependencies
type MockXrayInstanceRepo struct {
	mock.Mock
}

func (m *MockXrayInstanceRepo) Create(ctx context.Context, instance *xray.XrayInstance) error {
	args := m.Called(ctx, instance)
	return args.Error(0)
}

func (m *MockXrayInstanceRepo) FindByID(ctx context.Context, id uuid.UUID) (*xray.XrayInstance, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*xray.XrayInstance), args.Error(1)
}

func (m *MockXrayInstanceRepo) FindByNodeID(ctx context.Context, nodeID uuid.UUID) (*xray.XrayInstance, error) {
	args := m.Called(ctx, nodeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*xray.XrayInstance), args.Error(1)
}

func (m *MockXrayInstanceRepo) FindRunning(ctx context.Context) ([]*xray.XrayInstance, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*xray.XrayInstance), args.Error(1)
}

func (m *MockXrayInstanceRepo) Update(ctx context.Context, instance *xray.XrayInstance) error {
	args := m.Called(ctx, instance)
	return args.Error(0)
}

func (m *MockXrayInstanceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockXrayInstanceRepo) List(ctx context.Context, offset, limit int) ([]*xray.XrayInstance, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*xray.XrayInstance), args.Error(1)
}

func (m *MockXrayInstanceRepo) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type MockInboundRepo struct {
	mock.Mock
}

func (m *MockInboundRepo) Create(ctx context.Context, inbound *xray.Inbound) error {
	args := m.Called(ctx, inbound)
	return args.Error(0)
}

func (m *MockInboundRepo) FindByID(ctx context.Context, id uuid.UUID) (*xray.Inbound, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*xray.Inbound), args.Error(1)
}

func (m *MockInboundRepo) FindByInstanceID(ctx context.Context, instanceID uuid.UUID) ([]*xray.Inbound, error) {
	args := m.Called(ctx, instanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*xray.Inbound), args.Error(1)
}

func (m *MockInboundRepo) FindEnabledByInstanceID(ctx context.Context, instanceID uuid.UUID) ([]*xray.Inbound, error) {
	args := m.Called(ctx, instanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*xray.Inbound), args.Error(1)
}

func (m *MockInboundRepo) Update(ctx context.Context, inbound *xray.Inbound) error {
	args := m.Called(ctx, inbound)
	return args.Error(0)
}

func (m *MockInboundRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInboundRepo) List(ctx context.Context, offset, limit int) ([]*xray.Inbound, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*xray.Inbound), args.Error(1)
}

func (m *MockInboundRepo) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type MockClientRepo struct {
	mock.Mock
}

func (m *MockClientRepo) Create(ctx context.Context, client *xray.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientRepo) FindByID(ctx context.Context, id uuid.UUID) (*xray.Client, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*xray.Client), args.Error(1)
}

func (m *MockClientRepo) FindByInboundID(ctx context.Context, inboundID uuid.UUID) ([]*xray.Client, error) {
	args := m.Called(ctx, inboundID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*xray.Client), args.Error(1)
}

func (m *MockClientRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*xray.Client, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*xray.Client), args.Error(1)
}

func (m *MockClientRepo) FindByUUID(ctx context.Context, clientUUID string) (*xray.Client, error) {
	args := m.Called(ctx, clientUUID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*xray.Client), args.Error(1)
}

func (m *MockClientRepo) FindEnabledByInboundID(ctx context.Context, inboundID uuid.UUID) ([]*xray.Client, error) {
	args := m.Called(ctx, inboundID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*xray.Client), args.Error(1)
}

func (m *MockClientRepo) Update(ctx context.Context, client *xray.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientRepo) List(ctx context.Context, offset, limit int) ([]*xray.Client, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*xray.Client), args.Error(1)
}

func (m *MockClientRepo) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

type MockRealityConfigRepo struct {
	mock.Mock
}

func (m *MockRealityConfigRepo) Create(ctx context.Context, config *xray.RealityConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockRealityConfigRepo) FindByID(ctx context.Context, id uuid.UUID) (*xray.RealityConfig, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*xray.RealityConfig), args.Error(1)
}

func (m *MockRealityConfigRepo) FindByInboundID(ctx context.Context, inboundID uuid.UUID) (*xray.RealityConfig, error) {
	args := m.Called(ctx, inboundID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*xray.RealityConfig), args.Error(1)
}

func (m *MockRealityConfigRepo) Update(ctx context.Context, config *xray.RealityConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockRealityConfigRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockAuditRepo struct {
	mock.Mock
}

func (m *MockAuditRepo) Create(ctx context.Context, log interface{}) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

type MockConfigGenerator struct {
	mock.Mock
}

func (m *MockConfigGenerator) GenerateJSON(ctx context.Context, instanceID uuid.UUID) ([]byte, error) {
	args := m.Called(ctx, instanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

type MockConfigWriter struct {
	mock.Mock
}

func (m *MockConfigWriter) Write(ctx context.Context, instanceID uuid.UUID, configJSON []byte) error {
	args := m.Called(ctx, instanceID, configJSON)
	return args.Error(0)
}

func (m *MockConfigWriter) Read(ctx context.Context, instanceID uuid.UUID) ([]byte, error) {
	args := m.Called(ctx, instanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockConfigWriter) Backup(ctx context.Context, instanceID uuid.UUID) error {
	args := m.Called(ctx, instanceID)
	return args.Error(0)
}

func (m *MockConfigWriter) Restore(ctx context.Context, instanceID uuid.UUID, backupTime time.Time) error {
	args := m.Called(ctx, instanceID, backupTime)
	return args.Error(0)
}

func (m *MockConfigWriter) Delete(ctx context.Context, instanceID uuid.UUID) error {
	args := m.Called(ctx, instanceID)
	return args.Error(0)
}

func (m *MockConfigWriter) DeleteBackup(ctx context.Context, instanceID uuid.UUID, timestamp int64) error {
	args := m.Called(ctx, instanceID, timestamp)
	return args.Error(0)
}

func (m *MockConfigWriter) GetPath(instanceID uuid.UUID) string {
	args := m.Called(instanceID)
	return args.String(0)
}

func (m *MockConfigWriter) ListBackups(ctx context.Context, instanceID uuid.UUID) ([]xrayConfig.BackupInfo, error) {
	args := m.Called(ctx, instanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]xrayConfig.BackupInfo), args.Error(1)
}

type MockProcessManager struct {
	mock.Mock
}

func (m *MockProcessManager) Start(ctx context.Context, instanceID uuid.UUID) error {
	args := m.Called(ctx, instanceID)
	return args.Error(0)
}

func (m *MockProcessManager) Stop(ctx context.Context, instanceID uuid.UUID) error {
	args := m.Called(ctx, instanceID)
	return args.Error(0)
}

func (m *MockProcessManager) Restart(ctx context.Context, instanceID uuid.UUID) error {
	args := m.Called(ctx, instanceID)
	return args.Error(0)
}

func (m *MockProcessManager) Reload(ctx context.Context, instanceID uuid.UUID) error {
	args := m.Called(ctx, instanceID)
	return args.Error(0)
}

func (m *MockProcessManager) Status(ctx context.Context, instanceID uuid.UUID) (interface{}, error) {
	args := m.Called(ctx, instanceID)
	return args.Get(0), args.Error(1)
}

func (m *MockProcessManager) IsRunning(ctx context.Context, instanceID uuid.UUID) (bool, error) {
	args := m.Called(ctx, instanceID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProcessManager) GetProcessID(ctx context.Context, instanceID uuid.UUID) (int, error) {
	args := m.Called(ctx, instanceID)
	return args.Int(0), args.Error(1)
}

// Helper function to create test user
func createTestUser() *user.User {
	email, _ := user.NewEmail("test@example.com")
	password, _ := user.NewPassword("hashed_password")
	profile := user.NewProfile("Test", "User", "", "")
	testUser, _ := user.NewUser(email, password, profile)
	return testUser
}

// Helper function to create test instance
func createTestInstance() *xray.XrayInstance {
	nodeID := uuid.New()
	instance, _ := xray.NewXrayInstance(nodeID, "test-instance", xray.ProtocolVLESS, 1080)
	return instance
}

// Helper function to create test inbound
func createTestInbound(instanceID uuid.UUID) *xray.Inbound {
	inbound, _ := xray.NewInbound(instanceID, xray.ProtocolVLESS, 443, "tcp", "0.0.0.0", xray.SecurityREALITY)
	inbound.Enable()
	return inbound
}

// Test: Successful provisioning
func TestProvisionUserToXray_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	testUser := createTestUser()
	testInstance := createTestInstance()
	testInbound := createTestInbound(testInstance.ID)
	
	mockXrayInstanceRepo := new(MockXrayInstanceRepo)
	mockInboundRepo := new(MockInboundRepo)
	mockClientRepo := new(MockClientRepo)
	mockRealityRepo := new(MockRealityConfigRepo)
	mockAuditRepo := new(MockAuditRepo)
	mockConfigGenerator := new(MockConfigGenerator)
	mockConfigWriter := new(MockConfigWriter)
	mockProcessManager := new(MockProcessManager)
	log := logger.New("info", "json")

	service := NewXrayProvisioningService(
		mockXrayInstanceRepo,
		mockInboundRepo,
		mockClientRepo,
		mockRealityRepo,
		mockAuditRepo,
		mockConfigGenerator,
		mockConfigWriter,
		mockProcessManager,
		log,
	)

	// Expectations
	mockClientRepo.On("FindByUserID", ctx, testUser.ID).Return([]*xray.Client{}, nil)
	mockXrayInstanceRepo.On("FindRunning", ctx).Return([]*xray.XrayInstance{testInstance}, nil)
	mockInboundRepo.On("FindEnabledByInstanceID", ctx, testInstance.ID).Return([]*xray.Inbound{testInbound}, nil)
	mockClientRepo.On("Create", ctx, mock.AnythingOfType("*xray.Client")).Return(nil)
	mockAuditRepo.On("Create", ctx, mock.Anything).Return(nil)
	mockConfigGenerator.On("GenerateJSON", ctx, testInstance.ID).Return([]byte("{}"), nil)
	mockConfigWriter.On("Backup", ctx, testInstance.ID).Return(nil)
	mockConfigWriter.On("Write", ctx, testInstance.ID, []byte("{}")).Return(nil)
	mockProcessManager.On("Reload", ctx, testInstance.ID).Return(nil)

	// Execute
	err := service.ProvisionUserToXray(ctx, testUser, "127.0.0.1", "test-agent")

	// Assert
	assert.NoError(t, err)
	mockClientRepo.AssertExpectations(t)
	mockXrayInstanceRepo.AssertExpectations(t)
	mockInboundRepo.AssertExpectations(t)
	mockConfigGenerator.AssertExpectations(t)
	mockConfigWriter.AssertExpectations(t)
	mockProcessManager.AssertExpectations(t)
}

// Test: Idempotency - existing client
func TestProvisionUserToXray_ExistingClient(t *testing.T) {
	// Setup
	ctx := context.Background()
	testUser := createTestUser()
	testInstance := createTestInstance()
	testInbound := createTestInbound(testInstance.ID)
	
	existingClient, _ := xray.NewClient(testInbound.ID, testUser.ID, uuid.New().String(), "xtls-rprx-vision", testUser.Email.String())
	
	mockClientRepo := new(MockClientRepo)
	mockAuditRepo := new(MockAuditRepo)
	log := logger.New("info", "json")

	service := NewXrayProvisioningService(
		nil,
		nil,
		mockClientRepo,
		nil,
		mockAuditRepo,
		nil,
		nil,
		nil,
		log,
	)

	// Expectations
	mockClientRepo.On("FindByUserID", ctx, testUser.ID).Return([]*xray.Client{existingClient}, nil)
	mockAuditRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Execute
	err := service.ProvisionUserToXray(ctx, testUser, "127.0.0.1", "test-agent")

	// Assert
	assert.NoError(t, err)
	mockClientRepo.AssertExpectations(t)
	mockClientRepo.AssertNotCalled(t, "Create")
}

// Test: Config generation error - client should be rolled back
func TestProvisionUserToXray_ConfigGenerationError_ClientRollback(t *testing.T) {
	// Setup
	ctx := context.Background()
	testUser := createTestUser()
	testInstance := createTestInstance()
	testInbound := createTestInbound(testInstance.ID)
	
	mockXrayInstanceRepo := new(MockXrayInstanceRepo)
	mockInboundRepo := new(MockInboundRepo)
	mockClientRepo := new(MockClientRepo)
	mockRealityRepo := new(MockRealityConfigRepo)
	mockAuditRepo := new(MockAuditRepo)
	mockConfigGenerator := new(MockConfigGenerator)
	mockConfigWriter := new(MockConfigWriter)
	log := logger.New("info", "json")

	service := NewXrayProvisioningService(
		mockXrayInstanceRepo,
		mockInboundRepo,
		mockClientRepo,
		mockRealityRepo,
		mockAuditRepo,
		mockConfigGenerator,
		mockConfigWriter,
		nil,
		log,
	)

	// Expectations
	mockClientRepo.On("FindByUserID", ctx, testUser.ID).Return([]*xray.Client{}, nil)
	mockXrayInstanceRepo.On("FindRunning", ctx).Return([]*xray.XrayInstance{testInstance}, nil)
	mockInboundRepo.On("FindEnabledByInstanceID", ctx, testInstance.ID).Return([]*xray.Inbound{testInbound}, nil)
	mockClientRepo.On("Create", ctx, mock.AnythingOfType("*xray.Client")).Return(nil)
	mockAuditRepo.On("Create", ctx, mock.Anything).Return(nil)
	mockConfigGenerator.On("GenerateJSON", ctx, testInstance.ID).Return(nil, errors.New("generation failed"))
	mockConfigWriter.On("Backup", ctx, testInstance.ID).Return(nil)
	mockClientRepo.On("Delete", ctx, mock.AnythingOfType("uuid.UUID")).Return(nil)

	// Execute
	err := service.ProvisionUserToXray(ctx, testUser, "127.0.0.1", "test-agent")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config reload failed, client rolled back")
	mockClientRepo.AssertCalled(t, "Delete", ctx, mock.AnythingOfType("uuid.UUID"))
}

// Test: Reload error with successful config rollback
func TestRegenerateAndReload_ReloadError_ConfigRollbackSuccess(t *testing.T) {
	// Setup
	ctx := context.Background()
	instanceID := uuid.New()
	userID := uuid.New()
	
	mockConfigGenerator := new(MockConfigGenerator)
	mockConfigWriter := new(MockConfigWriter)
	mockProcessManager := new(MockProcessManager)
	mockAuditRepo := new(MockAuditRepo)
	log := logger.New("info", "json")

	service := NewXrayProvisioningService(
		nil,
		nil,
		nil,
		nil,
		mockAuditRepo,
		mockConfigGenerator,
		mockConfigWriter,
		mockProcessManager,
		log,
	)

	// Expectations
	mockConfigGenerator.On("GenerateJSON", ctx, instanceID).Return([]byte("{}"), nil)
	mockConfigWriter.On("Backup", ctx, instanceID).Return(nil)
	mockConfigWriter.On("Write", ctx, instanceID, []byte("{}")).Return(nil)
	mockProcessManager.On("Reload", ctx, instanceID).Return(errors.New("reload failed"))
	mockConfigWriter.On("ListBackups", ctx, instanceID).Return([]xrayConfig.BackupInfo{{
		InstanceID: instanceID,
		Timestamp:  time.Unix(123456, 0).UTC(),
		Path:       "/backup/test.json",
		Size:       100,
	}}, nil)
	mockConfigWriter.On("Restore", ctx, instanceID, mock.MatchedBy(func(t time.Time) bool {
		return t.Unix() == 123456
	})).Return(nil)
	mockProcessManager.On("Reload", ctx, instanceID).Return(nil).Once()
	mockAuditRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Execute
	err := service.RegenerateAndReload(ctx, instanceID, userID, "127.0.0.1", "test-agent")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reload failed, rolled back to previous config")
	mockConfigWriter.AssertCalled(t, "Restore", ctx, instanceID, int64(123456))
}

// Test: Reload error with failed config rollback
func TestRegenerateAndReload_ReloadError_ConfigRollbackFailed(t *testing.T) {
	// Setup
	ctx := context.Background()
	instanceID := uuid.New()
	userID := uuid.New()
	
	mockConfigGenerator := new(MockConfigGenerator)
	mockConfigWriter := new(MockConfigWriter)
	mockProcessManager := new(MockProcessManager)
	mockAuditRepo := new(MockAuditRepo)
	log := logger.New("info", "json")

	service := NewXrayProvisioningService(
		nil,
		nil,
		nil,
		nil,
		mockAuditRepo,
		mockConfigGenerator,
		mockConfigWriter,
		mockProcessManager,
		log,
	)

	// Expectations
	mockConfigGenerator.On("GenerateJSON", ctx, instanceID).Return([]byte("{}"), nil)
	mockConfigWriter.On("Backup", ctx, instanceID).Return(nil)
	mockConfigWriter.On("Write", ctx, instanceID, []byte("{}")).Return(nil)
	mockProcessManager.On("Reload", ctx, instanceID).Return(errors.New("reload failed"))
	mockConfigWriter.On("ListBackups", ctx, instanceID).Return(nil, errors.New("no backups"))
	mockAuditRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Execute
	err := service.RegenerateAndReload(ctx, instanceID, userID, "127.0.0.1", "test-agent")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reload failed and config rollback failed")
}

// Test: Config validation failure
func TestRegenerateAndReload_ConfigValidationFailed(t *testing.T) {
	// Setup
	ctx := context.Background()
	instanceID := uuid.New()
	userID := uuid.New()
	
	mockConfigGenerator := new(MockConfigGenerator)
	mockConfigWriter := new(MockConfigWriter)
	mockBinaryManager := new(MockBinaryManager)
	mockAuditRepo := new(MockAuditRepo)
	log := logger.New("info", "json")

	service := NewXrayProvisioningService(
		nil,
		nil,
		nil,
		nil,
		mockAuditRepo,
		mockConfigGenerator,
		mockConfigWriter,
		nil,
		mockBinaryManager,
		log,
	)

	// Expectations
	mockConfigGenerator.On("GenerateJSON", ctx, instanceID).Return([]byte("{}"), nil)
	mockConfigWriter.On("Backup", ctx, instanceID).Return(nil)
	mockConfigWriter.On("Write", ctx, instanceID, []byte("{}")).Return(nil)
	mockConfigWriter.On("GetPath", instanceID).Return("/etc/xray/test.json")
	mockBinaryManager.On("ValidateConfig", ctx, "/etc/xray/test.json").Return(errors.New("invalid config"))
	mockConfigWriter.On("ListBackups", ctx, instanceID).Return([]xrayConfig.BackupInfo{{
		InstanceID: instanceID,
		Timestamp:  time.Unix(123456, 0).UTC(),
		Path:       "/backup/test.json",
		Size:       100,
	}}, nil)
	mockConfigWriter.On("Restore", ctx, instanceID, mock.MatchedBy(func(t time.Time) bool {
		return t.Unix() == 123456
	})).Return(nil)
	mockAuditRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Execute
	err := service.RegenerateAndReload(ctx, instanceID, userID, "127.0.0.1", "test-agent")

	// Assert
	assert.Error(t, err)
	var provErr *ProvisioningError
	if errors.As(err, &provErr) {
		assert.Equal(t, ErrorClassNonRetryable, provErr.Class)
	}
	mockBinaryManager.AssertCalled(t, "ValidateConfig", ctx, "/etc/xray/test.json")
	mockConfigWriter.AssertCalled(t, "Restore", ctx, instanceID, int64(123456))
}

// Test: Reload timeout
func TestRegenerateAndReload_ReloadTimeout(t *testing.T) {
	// Setup
	ctx := context.Background()
	instanceID := uuid.New()
	userID := uuid.New()
	
	mockConfigGenerator := new(MockConfigGenerator)
	mockConfigWriter := new(MockConfigWriter)
	mockBinaryManager := new(MockBinaryManager)
	mockProcessManager := new(MockProcessManager)
	mockAuditRepo := new(MockAuditRepo)
	log := logger.New("info", "json")

	service := NewXrayProvisioningService(
		nil,
		nil,
		nil,
		nil,
		mockAuditRepo,
		mockConfigGenerator,
		mockConfigWriter,
		mockProcessManager,
		mockBinaryManager,
		log,
	)

	// Expectations
	mockConfigGenerator.On("GenerateJSON", ctx, instanceID).Return([]byte("{}"), nil)
	mockConfigWriter.On("Backup", ctx, instanceID).Return(nil)
	mockConfigWriter.On("Write", ctx, instanceID, []byte("{}")).Return(nil)
	mockConfigWriter.On("GetPath", instanceID).Return("/etc/xray/test.json")
	mockBinaryManager.On("ValidateConfig", ctx, "/etc/xray/test.json").Return(nil)
	mockProcessManager.On("Reload", mock.Anything, instanceID).Return(context.DeadlineExceeded)
	mockConfigWriter.On("ListBackups", ctx, instanceID).Return([]xrayConfig.BackupInfo{{
		InstanceID: instanceID,
		Timestamp:  time.Unix(123456, 0).UTC(),
		Path:       "/backup/test.json",
		Size:       100,
	}}, nil)
	mockConfigWriter.On("Restore", ctx, instanceID, mock.MatchedBy(func(t time.Time) bool {
		return t.Unix() == 123456
	})).Return(nil)
	mockProcessManager.On("Reload", mock.Anything, instanceID).Return(nil).Once()
	mockAuditRepo.On("Create", ctx, mock.Anything).Return(nil)

	// Execute
	err := service.RegenerateAndReload(ctx, instanceID, userID, "127.0.0.1", "test-agent")

	// Assert
	assert.Error(t, err)
	var provErr *ProvisioningError
	if errors.As(err, &provErr) {
		assert.Equal(t, ErrorClassRetryable, provErr.Class)
		assert.ErrorIs(t, provErr.Err, ErrReloadTimeout)
	}
}

// Test: Parallel provisioning (race condition prevention)
func TestProvisionUserToXray_ParallelRequests(t *testing.T) {
	// Setup
	ctx := context.Background()
	testUser := createTestUser()
	testInstance := createTestInstance()
	testInbound := createTestInbound(testInstance.ID)
	
	mockXrayInstanceRepo := new(MockXrayInstanceRepo)
	mockInboundRepo := new(MockInboundRepo)
	mockClientRepo := new(MockClientRepo)
	mockRealityRepo := new(MockRealityConfigRepo)
	mockAuditRepo := new(MockAuditRepo)
	mockConfigGenerator := new(MockConfigGenerator)
	mockConfigWriter := new(MockConfigWriter)
	mockProcessManager := new(MockProcessManager)
	mockBinaryManager := new(MockBinaryManager)
	log := logger.New("info", "json")

	service := NewXrayProvisioningService(
		mockXrayInstanceRepo,
		mockInboundRepo,
		mockClientRepo,
		mockRealityRepo,
		mockAuditRepo,
		mockConfigGenerator,
		mockConfigWriter,
		mockProcessManager,
		mockBinaryManager,
		log,
	)

	// Expectations - first call creates, second finds existing
	mockClientRepo.On("FindByUserID", ctx, testUser.ID).Return([]*xray.Client{}, nil).Once()
	mockClientRepo.On("FindByUserID", ctx, testUser.ID).Return([]*xray.Client{createTestClient()}, nil).Once()
	mockXrayInstanceRepo.On("FindRunning", ctx).Return([]*xray.XrayInstance{testInstance}, nil)
	mockInboundRepo.On("FindEnabledByInstanceID", ctx, testInstance.ID).Return([]*xray.Inbound{testInbound}, nil)
	mockClientRepo.On("Create", ctx, mock.AnythingOfType("*xray.Client")).Return(nil)
	mockProcessManager.On("IsRunning", mock.Anything, testInstance.ID).Return(true, nil)
	mockProcessManager.On("Status", mock.Anything, testInstance.ID).Return(nil, nil)
	mockAuditRepo.On("Create", ctx, mock.Anything).Return(nil)
	mockConfigGenerator.On("GenerateJSON", ctx, testInstance.ID).Return([]byte("{}"), nil)
	mockConfigWriter.On("Backup", ctx, testInstance.ID).Return(nil)
	mockConfigWriter.On("Write", ctx, testInstance.ID, []byte("{}")).Return(nil)
	mockConfigWriter.On("GetPath", testInstance.ID).Return("/etc/xray/test.json")
	mockBinaryManager.On("ValidateConfig", ctx, "/etc/xray/test.json").Return(nil)
	mockProcessManager.On("Reload", mock.Anything, testInstance.ID).Return(nil)
	mockConfigWriter.On("ListBackups", ctx, testInstance.ID).Return([]xrayConfig.BackupInfo{}, nil)

	// Execute parallel requests
	var wg sync.WaitGroup
	errors := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		errors <- service.ProvisionUserToXray(ctx, testUser, "127.0.0.1", "test-agent-1")
	}()
	go func() {
		defer wg.Done()
		errors <- service.ProvisionUserToXray(ctx, testUser, "127.0.0.1", "test-agent-2")
	}()

	wg.Wait()
	close(errors)

	// Assert - both should succeed, but only one creates client
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
		}
	}

	assert.Equal(t, 0, errorCount, "Both parallel requests should succeed")
	mockClientRepo.AssertNumberOfCalls(t, "Create", 1) // Only one create call
}

// Helper to create test client
func createTestClient() *xray.Client {
	inboundID := uuid.New()
	userID := uuid.New()
	client, _ := xray.NewClient(inboundID, userID, uuid.New().String(), "xtls-rprx-vision", "test@example.com")
	return client
}

// Test: Error classification
func TestClassifyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorClass
	}{
		{"nil error", nil, ErrorClassSkippable},
		{"no running instances", ErrNoRunningInstances, ErrorClassSkippable},
		{"no enabled inbounds", ErrNoEnabledInbounds, ErrorClassSkippable},
		{"reload timeout", ErrReloadTimeout, ErrorClassRetryable},
		{"instance unhealthy", ErrInstanceUnhealthy, ErrorClassRetryable},
		{"config validation failed", ErrConfigValidationFailed, ErrorClassNonRetryable},
		{"generic error", errors.New("some error"), ErrorClassRetryable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			class := ClassifyError(tt.err)
			assert.Equal(t, tt.expected, class)
		})
	}
}

type MockBinaryManager struct {
	mock.Mock
}

func (m *MockBinaryManager) Detect(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockBinaryManager) Validate(ctx context.Context, binaryPath string) error {
	args := m.Called(ctx, binaryPath)
	return args.Error(0)
}

func (m *MockBinaryManager) ValidateConfig(ctx context.Context, configPath string) error {
	args := m.Called(ctx, configPath)
	return args.Error(0)
}

func (m *MockBinaryManager) CurrentVersion(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockBinaryManager) LatestVersion(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockBinaryManager) Download(ctx context.Context, version string) error {
	args := m.Called(ctx, version)
	return args.Error(0)
}

func (m *MockBinaryManager) Upgrade(ctx context.Context, version string) error {
	args := m.Called(ctx, version)
	return args.Error(0)
}

func (m *MockBinaryManager) GetPath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockBinaryManager) IsInstalled(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}
