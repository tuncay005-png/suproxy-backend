package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/audit"
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
