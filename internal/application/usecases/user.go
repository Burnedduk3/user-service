// internal/application/usecases/user_usecases.go
package usecases

import (
	"context"
	"errors"
	"net/mail"
	"user-service/internal/application/dto"
	"user-service/internal/application/ports"
	userErrors "user-service/internal/domain/errors"
	"user-service/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// UserUseCases defines the interface for user business operations
type UserUseCases interface {
	CreateUser(ctx context.Context, request *dto.CreateUserRequestDTO) (*dto.UserResponseDTO, error)
	GetUserByID(ctx context.Context, id uint) (*dto.UserResponseDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponseDTO, error)
	UpdateUser(ctx context.Context, id uint, request *dto.UpdateUserRequestDTO) (*dto.UserResponseDTO, error)
	ListUsers(ctx context.Context, page, pageSize int) (*dto.UserListResponseDTO, error)
}

// userUseCasesImpl implements UserUseCases interface
type userUseCasesImpl struct {
	userRepo ports.UserRepository
	logger   logger.Logger
}

// NewUserUseCases creates a new instance of user use cases
func NewUserUseCases(userRepo ports.UserRepository, log logger.Logger) UserUseCases {
	return &userUseCasesImpl{
		userRepo: userRepo,
		logger:   log.With("component", "user_usecases"),
	}
}

func (uc *userUseCasesImpl) CreateUser(ctx context.Context, request *dto.CreateUserRequestDTO) (*dto.UserResponseDTO, error) {
	uc.logger.Info("CreateUser use case called", "email", request.Email)
	if _, err := mail.ParseAddress(request.Email); err != nil {
		return nil, userErrors.ErrInvalidUserEmail
	}

	if _, err := uc.userRepo.ExistsByEmail(ctx, request.Email); err != nil {
		return nil, userErrors.ErrUserAlreadyExists
	}

	domainEntity, err := request.ToEntity()

	if err != nil {
		return nil, err
	}

	domainEntity.Password, err = hashPassword(domainEntity.Password)

	if err != nil {
		return nil, err
	}

	createUser, err := uc.userRepo.Create(ctx, domainEntity)

	if err != nil {
		switch {
		case errors.Is(err, userErrors.ErrFailedToCheckUserExistance):
			return nil, userErrors.ErrFailedToCheckUserExistance
		default:
			return nil, userErrors.ErrFailedToCreateUser

		}
	}

	uc.logger.Info("CreateUser success", "email", request.Email)

	return dto.UserToResponseDTO(createUser), nil
}

// GetUserByID retrieves a user by their ID
func (uc *userUseCasesImpl) GetUserByID(ctx context.Context, id uint) (*dto.UserResponseDTO, error) {
	uc.logger.Info("GetUserByID use case called", "user_id", id)

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	uc.logger.Info("GetUserByID success", "user_id", id)
	return dto.UserToResponseDTO(user), nil
}

// GetUserByEmail retrieves a user by their email address
func (uc *userUseCasesImpl) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponseDTO, error) {
	uc.logger.Info("GetUserByEmail use case called", "email", email)

	user, err := uc.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}
	uc.logger.Info("GetUserByEmail success", "user_id", user.ID)
	return dto.UserToResponseDTO(user), nil
}

// UpdateUser updates an existing user
func (uc *userUseCasesImpl) UpdateUser(ctx context.Context, id uint, request *dto.UpdateUserRequestDTO) (*dto.UserResponseDTO, error) {
	uc.logger.Info("UpdateUser use case called", "user_id", id)

	existingUser, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if request.Phone != "" {
		existingUser.Phone = request.Phone
	}

	if request.LastName != "" {
		existingUser.LastName = request.LastName
	}

	if request.FirstName != "" {
		existingUser.FirstName = request.FirstName
	}
	uc.userRepo.Update(ctx, existingUser)

	uc.logger.Info("UpdateUser success", "user_id", id)

	return dto.UserToResponseDTO(existingUser), nil
}

// ListUsers retrieves a paginated list of users
func (uc *userUseCasesImpl) ListUsers(ctx context.Context, page, pageSize int) (*dto.UserListResponseDTO, error) {
	uc.logger.Info("ListUsers use case called", "page", page, "page_size", pageSize)

	if page < 0 {
		page = 0
	}

	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, err := uc.userRepo.List(ctx, pageSize, page)

	if err != nil {
		return nil, err
	}

	response := dto.UsersToResponseDTOs(users)

	uc.logger.Info("ListUsers success", "page", page, "page_size", pageSize)

	return &dto.UserListResponseDTO{
		Users:    response,
		Page:     page,
		PageSize: pageSize,
		Total:    len(users),
	}, nil
}

// hashPassword hashes a plain text password using bcrypt
func hashPassword(password string) (string, error) {
	hashInBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashInBytes), nil
}
