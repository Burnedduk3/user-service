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
	// TODO: Implement using TDD
	// 1. Check if user with email already exists
	// 2. Convert domain entity to GORM model
	// 3. Create in database
	// 4. Convert back to domain entity
	panic("implement me")
}

// GetByID implements ports.UserRepository
func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	// TODO: Implement using TDD
	// 1. Find user by ID in database
	// 2. Handle not found case
	// 3. Convert GORM model to domain entity
	panic("implement me")
}

// GetByEmail implements ports.UserRepository
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	// TODO: Implement using TDD
	// 1. Find user by email in database
	// 2. Handle not found case
	// 3. Convert GORM model to domain entity
	panic("implement me")
}

// Update implements ports.UserRepository
func (r *GormUserRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	// TODO: Implement using TDD
	// 1. Find existing user
	// 2. Update fields
	// 3. Save to database
	// 4. Return updated entity
	panic("implement me")
}

// Delete implements ports.UserRepository
func (r *GormUserRepository) Delete(ctx context.Context, id uint) error {
	// TODO: Implement using TDD
	// 1. Find user by ID
	// 2. Perform soft delete
	// 3. Handle not found case
	panic("implement me")
}

// ExistsByEmail implements ports.UserRepository
func (r *GormUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	// TODO: Implement using TDD
	// 1. Count users with given email
	// 2. Return true if count > 0
	panic("implement me")
}

// List implements ports.UserRepository
func (r *GormUserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	// TODO: Implement using TDD
	// 1. Query users with pagination
	// 2. Convert to domain entities
	panic("implement me")
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
	entities := make([]*entities.User, 0, len(models))
	for _, model := range models {
		entities = append(entities, r.toEntity(&model))
	}
	return entities
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
