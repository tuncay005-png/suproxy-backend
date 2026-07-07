package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/audit"
	"github.com/suproxy/backend/internal/infrastructure/metrics"
	"gorm.io/gorm"
)

type AuditLogModel struct {
	ID         uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID     uuid.UUID              `gorm:"type:uuid;index"`
	Action     string                 `gorm:"type:varchar(50);not null;index"`
	EntityType string                 `gorm:"type:varchar(50);not null"`
	EntityID   uuid.UUID              `gorm:"type:uuid;index"`
	IPAddress  string                 `gorm:"type:varchar(45)"`
	UserAgent  string                 `gorm:"type:text"`
	Metadata   map[string]interface{} `gorm:"type:jsonb"`
	CreatedAt  time.Time              `gorm:"not null;default:CURRENT_TIMESTAMP;index"`
}

func (AuditLogModel) TableName() string {
	return "audit_logs"
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) audit.Repository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *audit.Log) error {
	model := &AuditLogModel{
		ID:         log.ID,
		UserID:     log.UserID,
		Action:     string(log.Action),
		EntityType: log.EntityType,
		EntityID:   log.EntityID,
		IPAddress:  log.IPAddress,
		UserAgent:  log.UserAgent,
		Metadata:   log.Metadata,
		CreatedAt:  log.CreatedAt,
	}
	
	// Record metric
	metrics.IncAuditLogs()
	
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *auditLogRepository) FindByID(ctx context.Context, id uuid.UUID) (*audit.Log, error) {
	var model AuditLogModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return toDomainAuditLog(&model), nil
}

func (r *auditLogRepository) FindByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*audit.Log, error) {
	var models []AuditLogModel
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, from, to).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	logs := make([]*audit.Log, 0, len(models))
	for _, model := range models {
		logs = append(logs, toDomainAuditLog(&model))
	}
	return logs, nil
}

func (r *auditLogRepository) FindByEntityID(ctx context.Context, entityType string, entityID uuid.UUID) ([]*audit.Log, error) {
	var models []AuditLogModel
	if err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	logs := make([]*audit.Log, 0, len(models))
	for _, model := range models {
		logs = append(logs, toDomainAuditLog(&model))
	}
	return logs, nil
}

func (r *auditLogRepository) List(ctx context.Context, offset, limit int) ([]*audit.Log, error) {
	var models []AuditLogModel
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, err
	}

	logs := make([]*audit.Log, 0, len(models))
	for _, model := range models {
		logs = append(logs, toDomainAuditLog(&model))
	}
	return logs, nil
}

func (r *auditLogRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&AuditLogModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *auditLogRepository) ListWithFilters(ctx context.Context, filters audit.AuditFilters) ([]*audit.Log, int64, error) {
	query := r.db.WithContext(ctx).Model(&AuditLogModel{})

	// Apply filters
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.Action != nil && *filters.Action != "" {
		query = query.Where("action = ?", *filters.Action)
	}
	if filters.EntityType != nil && *filters.EntityType != "" {
		query = query.Where("entity_type = ?", *filters.EntityType)
	}
	if filters.EntityID != nil {
		query = query.Where("entity_id = ?", *filters.EntityID)
	}
	if filters.IPAddress != nil && *filters.IPAddress != "" {
		query = query.Where("ip_address = ?", *filters.IPAddress)
	}
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}

	// Count total (before pagination)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := filters.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// Execute query
	var models []AuditLogModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	logs := make([]*audit.Log, 0, len(models))
	for _, model := range models {
		logs = append(logs, toDomainAuditLog(&model))
	}

	return logs, total, nil
}

func (r *auditLogRepository) CountByAction(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		Action string
		Count  int64
	}

	if err := r.db.WithContext(ctx).
		Model(&AuditLogModel{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Find(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Action] = result.Count
	}
	return counts, nil
}

func (r *auditLogRepository) CountByEntityType(ctx context.Context) (map[string]int64, error) {
	var results []struct {
		EntityType string
		Count      int64
	}

	if err := r.db.WithContext(ctx).
		Model(&AuditLogModel{}).
		Select("entity_type, COUNT(*) as count").
		Group("entity_type").
		Find(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.EntityType] = result.Count
	}
	return counts, nil
}

func (r *auditLogRepository) CountUniqueUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&AuditLogModel{}).
		Distinct("user_id").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *auditLogRepository) CountUniqueIPs(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&AuditLogModel{}).
		Distinct("ip_address").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *auditLogRepository) GetOldestLogDate(ctx context.Context) (*time.Time, error) {
	var model AuditLogModel
	if err := r.db.WithContext(ctx).
		Model(&AuditLogModel{}).
		Order("created_at ASC").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &model.CreatedAt, nil
}

func (r *auditLogRepository) GetNewestLogDate(ctx context.Context) (*time.Time, error) {
	var model AuditLogModel
	if err := r.db.WithContext(ctx).
		Model(&AuditLogModel{}).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &model.CreatedAt, nil
}

func (r *auditLogRepository) DeleteOlderThan(ctx context.Context, date time.Time) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", date).
		Delete(&AuditLogModel{}).Error
}

func toDomainAuditLog(m *AuditLogModel) *audit.Log {
	return &audit.Log{
		ID:         m.ID,
		UserID:     m.UserID,
		Action:     audit.Action(m.Action),
		EntityType: m.EntityType,
		EntityID:   m.EntityID,
		IPAddress:  m.IPAddress,
		UserAgent:  m.UserAgent,
		Metadata:   m.Metadata,
		CreatedAt:  m.CreatedAt,
	}
}
