package user_repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"user-service/internal/application/ports"
	"user-service/internal/domain/entities"
	domainErrors "user-service/internal/domain/errors"

	"gorm.io/gorm"
)

// UserModel represents the database model for users
type UserModel struct {
	ID        uint           `gorm:"primarykey"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Password  string         `gorm:"not null"`
	FirstName string         `gorm:"not null"`
	LastName  string         `gorm:"not null"`
	Phone     string         `gorm:""`
	Status    string         `gorm:"not null;default:'active'"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // For soft deletes
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}

// GormUserRepository implements the UserRepository interface using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GORM user repository
func NewGormUserRepository(db *gorm.DB) ports.UserRepository {
	return &GormUserRepository{db: db}
}

// Create implements ports.UserRepository
func (r *GormUserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	// Check if user already exists
	exists, err := r.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainErrors.ErrUserAlreadyExists
	}

	gormModel := r.toModel(user)

	// Create user in database
	if err := r.db.WithContext(ctx).Create(gormModel).Error; err != nil {
		return nil, r.handleError(err)
	}

	return r.toEntity(gormModel), nil
}

// GetByID implements ports.UserRepository
func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, r.handleError(err)
	}

	return r.toEntity(&model), nil
}

// GetByEmail implements ports.UserRepository
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, r.handleError(err)
	}

	return r.toEntity(&model), nil
}

// Update implements ports.UserRepository
func (r *GormUserRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	// First check if user exists
	var existingModel UserModel
	err := r.db.WithContext(ctx).Where("id = ?", user.ID).First(&existingModel).Error
	if err != nil {
		return nil, r.handleError(err)
	}

	// Convert entity to model for update
	updateModel := r.toModel(user)

	// Update the user
	err = r.db.WithContext(ctx).Model(&existingModel).Updates(updateModel).Error
	if err != nil {
		return nil, r.handleError(err)
	}

	// Fetch the updated user to return
	var updatedModel UserModel
	err = r.db.WithContext(ctx).Where("id = ?", user.ID).First(&updatedModel).Error
	if err != nil {
		return nil, r.handleError(err)
	}

	return r.toEntity(&updatedModel), nil
}

// Delete implements ports.UserRepository
func (r *GormUserRepository) Delete(ctx context.Context, id uint) error {
	// First check if user exists
	var model UserModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return r.handleError(err)
	}

	// Perform soft delete
	err = r.db.WithContext(ctx).Delete(&model).Error
	if err != nil {
		return r.handleError(err)
	}

	return nil
}

// ExistsByEmail implements ports.UserRepository
func (r *GormUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserModel{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, domainErrors.ErrFailedToCheckUserExistance
	}

	return count > 0, nil
}

// List implements ports.UserRepository
func (r *GormUserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var models []UserModel

	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, r.handleError(err)
	}

	return r.toEntities(models), nil
}

// Helper functions for conversion between domain entities and GORM models

func (r *GormUserRepository) toModel(user *entities.User) *UserModel {
	return &UserModel{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (r *GormUserRepository) toEntity(model *UserModel) *entities.User {
	return &entities.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		FirstName: model.FirstName,
		LastName:  model.LastName,
		Phone:     model.Phone,
		Status:    entities.UserStatus(model.Status),
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (r *GormUserRepository) toEntities(models []UserModel) []*entities.User {
	users := make([]*entities.User, 0, len(models))
	for _, model := range models {
		users = append(users, r.toEntity(&model))
	}
	return users
}

// Helper to convert GORM errors to domain errors
func (r *GormUserRepository) handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domainErrors.ErrUserNotFound
	}

	// Handle unique constraint violation for email
	if errors.Is(err, gorm.ErrDuplicatedKey) ||
		(err.Error() != "" && (strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint"))) {
		return domainErrors.ErrUserAlreadyExists
	}

	// Return original error for other cases
	return err
}
