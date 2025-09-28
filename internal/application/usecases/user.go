// internal/application/usecases/user_usecases.go
package usecases

import (
	"context"
	"errors"
	"user-service/internal/application/dto"
	"user-service/internal/application/ports"
	"user-service/pkg/logger"
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

	// TODO: Implement CreateUser use case
	// Steps to implement:
	// 1. Log the operation with email (already done)
	// 2. Check if user already exists using uc.userRepo.ExistsByEmail()
	// 3. If exists, return domainErrors.ErrUserAlreadyExists
	// 4. Convert DTO to domain entity using request.ToEntity()
	// 5. Hash the password using hashPassword() helper function
	// 6. Create user in repository using uc.userRepo.Create()
	// 7. Convert created user to response DTO using dto.UserToResponseDTO()
	// 8. Log success and return response

	return nil, errors.New("CreateUser not implemented yet")
}

// GetUserByID retrieves a user by their ID
func (uc *userUseCasesImpl) GetUserByID(ctx context.Context, id uint) (*dto.UserResponseDTO, error) {
	uc.logger.Info("GetUserByID use case called", "user_id", id)

	// TODO: Implement GetUserByID use case
	// Steps to implement:
	// 1. Log the operation with user_id (already done)
	// 2. Get user from repository using uc.userRepo.GetByID()
	// 3. Handle any errors from repository (return as-is)
	// 4. Convert user entity to response DTO using dto.UserToResponseDTO()
	// 5. Log success and return response

	return nil, errors.New("GetUserByID not implemented yet")
}

// GetUserByEmail retrieves a user by their email address
func (uc *userUseCasesImpl) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponseDTO, error) {
	uc.logger.Info("GetUserByEmail use case called", "email", email)

	// TODO: Implement GetUserByEmail use case
	// Steps to implement:
	// 1. Log the operation with email (already done)
	// 2. Get user from repository using uc.userRepo.GetByEmail()
	// 3. Handle any errors from repository (return as-is)
	// 4. Convert user entity to response DTO using dto.UserToResponseDTO()
	// 5. Log success and return response

	return nil, errors.New("GetUserByEmail not implemented yet")
}

// UpdateUser updates an existing user
func (uc *userUseCasesImpl) UpdateUser(ctx context.Context, id uint, request *dto.UpdateUserRequestDTO) (*dto.UserResponseDTO, error) {
	uc.logger.Info("UpdateUser use case called", "user_id", id)

	// TODO: Implement UpdateUser use case
	// Steps to implement:
	// 1. Log the operation with user_id (already done)
	// 2. Get existing user using uc.userRepo.GetByID()
	// 3. Handle error if user not found
	// 4. Update fields only if they are provided in request (not empty)
	//    - if request.FirstName != "" then existingUser.FirstName = request.FirstName
	//    - if request.LastName != "" then existingUser.LastName = request.LastName
	//    - if request.Phone != "" then existingUser.Phone = request.Phone
	// 5. Update user in repository using uc.userRepo.Update()
	// 6. Convert updated user to response DTO using dto.UserToResponseDTO()
	// 7. Log success and return response

	return nil, errors.New("UpdateUser not implemented yet")
}

// ListUsers retrieves a paginated list of users
func (uc *userUseCasesImpl) ListUsers(ctx context.Context, page, pageSize int) (*dto.UserListResponseDTO, error) {
	uc.logger.Info("ListUsers use case called", "page", page, "page_size", pageSize)

	// TODO: Implement ListUsers use case
	// Steps to implement:
	// 1. Log the operation with page and pageSize (already done)
	// 2. Validate and correct pagination parameters:
	//    - if page < 1 then page = 1
	//    - if pageSize < 1 or pageSize > 100 then pageSize = 10
	// 3. Calculate offset: offset = (page - 1) * pageSize
	// 4. Get users from repository using uc.userRepo.List(pageSize, offset)
	// 5. Handle any errors from repository
	// 6. Convert users to response DTOs using dto.UsersToResponseDTOs()
	// 7. Create UserListResponseDTO with:
	//    - Users: converted DTOs
	//    - Total: len(userDTOs) (for now, in real app you'd get total count)
	//    - Page: page
	//    - PageSize: pageSize
	//    - TotalPages: 1 (for now, in real app: (total + pageSize - 1) / pageSize)
	// 8. Log success and return response

	return nil, errors.New("ListUsers not implemented yet")
}

// Helper functions to implement

// TODO: Implement hashPassword function
// hashPassword hashes a plain text password using bcrypt
// func hashPassword(password string) (string, error) {
//     // Use bcrypt.GenerateFromPassword() with bcrypt.DefaultCost
//     // Return the hashed password as string
// }

// TODO: Implement verifyPassword function (for future use)
// verifyPassword verifies a password against its hash
// func verifyPassword(hashedPassword, password string) error {
//     // Use bcrypt.CompareHashAndPassword()
//     // Return error if passwords don't match
// }
