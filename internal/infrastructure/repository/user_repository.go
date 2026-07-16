package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/suproxy/backend/internal/domain/user"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	if u == nil {
		return errors.New("user cannot be nil")
	}

	userModel := toUserModel(u)
	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return user.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return toDomainUser(&model)
}

func (r *userRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
	var model UserModel
	if err := r.db.WithContext(ctx).Where("email = ?", email.String()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return toDomainUser(&model)
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	userModel := toUserModel(u)
	if err := r.db.WithContext(ctx).Save(userModel).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&UserModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return user.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*user.User, error) {
	var models []UserModel
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]*user.User, 0, len(models))
	for _, model := range models {
		u, err := toDomainUser(&model)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email user.Email) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("email = ?", email.String()).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) ListWithFilters(ctx context.Context, filters user.UserFilters) ([]*user.User, int64, error) {
	query := r.db.WithContext(ctx).Model(&UserModel{})

	// Apply filters
	if filters.Role != nil {
		query = query.Where("role = ?", string(*filters.Role))
	}

	if filters.Status != nil {
		query = query.Where("status = ?", string(*filters.Status))
	}

	if filters.Email != "" {
		query = query.Where("email LIKE ?", "%"+filters.Email+"%")
	}

	// Count total before pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortField := "created_at"
	if filters.SortBy != "" {
		sortField = filters.SortBy
	}

	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query = query.Order(sortField + " " + sortOrder)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// Execute query
	var models []UserModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain users
	users := make([]*user.User, 0, len(models))
	for _, model := range models {
		u, err := toDomainUser(&model)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}
