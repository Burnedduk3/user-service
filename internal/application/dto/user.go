package dto

import (
	"time"
	"user-service/internal/domain/entities"
)

// CreateUserRequestDTO for user creation
type CreateUserRequestDTO struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=15"`
}

// UpdateUserRequestDTO for user updates
type UpdateUserRequestDTO struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=15"`
}

// UserResponseDTO for user responses (excludes sensitive data)
type UserResponseDTO struct {
	ID        uint                `json:"id"`
	Email     string              `json:"email"`
	FirstName string              `json:"first_name"`
	LastName  string              `json:"last_name"`
	FullName  string              `json:"full_name"`
	Phone     string              `json:"phone"`
	Status    entities.UserStatus `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// UserListResponseDTO for paginated user lists
type UserListResponseDTO struct {
	Users    []*UserResponseDTO `json:"users"`
	Total    int                `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// Conversion methods
func (dto *CreateUserRequestDTO) ToEntity() (*entities.User, error) {
	return entities.NewUser(
		dto.Email,
		dto.Password,
		dto.FirstName,
		dto.LastName,
		dto.Phone,
	)
}

func UserToResponseDTO(user *entities.User) *UserResponseDTO {
	return &UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.FullName(),
		Phone:     user.Phone,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UsersToResponseDTOs(users []*entities.User) []*UserResponseDTO {
	dtos := make([]*UserResponseDTO, 0, len(users))
	for _, user := range users {
		dtos = append(dtos, UserToResponseDTO(user))
	}
	return dtos
}
