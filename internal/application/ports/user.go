package ports

import (
	"context"
	"user-service/internal/domain/entities"
)

// UserRepository defines the contract for user persistence
type UserRepository interface {
	// Create a new user
	Create(ctx context.Context, user *entities.User) (*entities.User, error)

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id uint) (*entities.User, error)

	// GetByEmail retrieves a user by their email (useful for login)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update an existing user
	Update(ctx context.Context, user *entities.User) (*entities.User, error)

	// Delete a user (soft delete recommended)
	Delete(ctx context.Context, id uint) error

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// List users with pagination (useful for admin features)
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)
}
