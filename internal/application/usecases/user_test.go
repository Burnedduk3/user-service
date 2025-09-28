package usecases

import (
	"context"
	"testing"
	"time"
	"user-service/internal/application/dto"
	"user-service/internal/domain/entities"
	domainErrors "user-service/internal/domain/errors"
	"user-service/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository implements the UserRepository interface for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func setupTestUseCases() (UserUseCases, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	log := logger.New("test")
	useCases := NewUserUseCases(mockRepo, log)
	return useCases, mockRepo
}

// CreateUser Tests
func TestUserUseCases_CreateUser_Success(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	request := &dto.CreateUserRequestDTO{
		Email:     "test@example.com",
		Password:  "SecurePass123",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
	}

	// Mock repository calls
	mockRepo.On("ExistsByEmail", ctx, "test@example.com").Return(false, nil)

	// Expected user to be created
	expectedCreatedUser := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "hashed_password", // This will be the hashed version
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Status:    entities.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("Create", ctx, mock.MatchedBy(func(user *entities.User) bool {
		return user.Email == "test@example.com" &&
			user.FirstName == "John" &&
			user.LastName == "Doe" &&
			user.Phone == "1234567890" &&
			user.Status == entities.UserStatusActive &&
			user.Password != "SecurePass123" // Password should be hashed
	})).Return(expectedCreatedUser, nil)

	// When
	result, err := useCases.CreateUser(ctx, request)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "John", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)
	assert.Equal(t, "John Doe", result.FullName)
	assert.Equal(t, "1234567890", result.Phone)
	assert.Equal(t, entities.UserStatusActive, result.Status)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_CreateUser_EmailAlreadyExists(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	request := &dto.CreateUserRequestDTO{
		Email:     "existing@example.com",
		Password:  "SecurePass123",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock repository to return true for existing email
	mockRepo.On("ExistsByEmail", ctx, "existing@example.com").Return(false, domainErrors.ErrUserAlreadyExists)

	// When
	result, err := useCases.CreateUser(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domainErrors.ErrUserAlreadyExists, err)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_CreateUser_InvalidUserData(t *testing.T) {
	// Given
	useCases, _ := setupTestUseCases()
	ctx := context.Background()

	request := &dto.CreateUserRequestDTO{
		Email:     "invalid-email", // Invalid email format
		Password:  "SecurePass123",
		FirstName: "John",
		LastName:  "Doe",
	}

	// When
	result, err := useCases.CreateUser(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Invalid email format")
}

func TestUserUseCases_CreateUser_RepositoryExistsError(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	hashInBytes, err := bcrypt.GenerateFromPassword([]byte("SecurePass123"), bcrypt.MinCost)

	request := &dto.CreateUserRequestDTO{
		Email:     "test@example.com",
		Password:  string(hashInBytes),
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock repository to return error when checking if email exists
	mockRepo.On("ExistsByEmail", ctx, "test@example.com").Return(true, nil)
	mockRepo.On("Create", ctx, mock.MatchedBy(func(user *entities.User) bool {
		return user.Email == "test@example.com"
	})).Return(nil, domainErrors.ErrFailedToCheckUserExistance)

	// When
	result, err := useCases.CreateUser(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to check user existence")

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_CreateUser_RepositoryCreateError(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	request := &dto.CreateUserRequestDTO{
		Email:     "test@example.com",
		Password:  "SecurePass123",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Mock successful email check but failed create
	mockRepo.On("ExistsByEmail", ctx, "test@example.com").Return(true, nil)
	mockRepo.On("Create", ctx, mock.Anything).Return(nil, assert.AnError)

	// When
	result, err := useCases.CreateUser(ctx, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create user")

	mockRepo.AssertExpectations(t)
}

// GetUserByID Tests
func TestUserUseCases_GetUserByID_Success(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	expectedUser := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Status:    entities.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", ctx, uint(1)).Return(expectedUser, nil)

	// When
	result, err := useCases.GetUserByID(ctx, 1)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "John Doe", result.FullName)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_GetUserByID_NotFound(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, uint(999)).Return(nil, domainErrors.ErrUserNotFound)

	// When
	result, err := useCases.GetUserByID(ctx, 999)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domainErrors.ErrUserNotFound, err)

	mockRepo.AssertExpectations(t)
}

// GetUserByEmail Tests
func TestUserUseCases_GetUserByEmail_Success(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	expectedUser := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Status:    entities.UserStatusActive,
	}

	mockRepo.On("GetByEmail", ctx, "test@example.com").Return(expectedUser, nil)

	// When
	result, err := useCases.GetUserByEmail(ctx, "test@example.com")

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "John", result.FirstName)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_GetUserByEmail_NotFound(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	mockRepo.On("GetByEmail", ctx, "notfound@example.com").Return(nil, domainErrors.ErrUserNotFound)

	// When
	result, err := useCases.GetUserByEmail(ctx, "notfound@example.com")

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domainErrors.ErrUserNotFound, err)

	mockRepo.AssertExpectations(t)
}

// UpdateUser Tests
func TestUserUseCases_UpdateUser_Success(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	existingUser := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Status:    entities.UserStatusActive,
	}

	updatedUser := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Johnny",
		LastName:  "Smith",
		Phone:     "0987654321",
		Status:    entities.UserStatusActive,
	}

	request := &dto.UpdateUserRequestDTO{
		FirstName: "Johnny",
		LastName:  "Smith",
		Phone:     "0987654321",
	}

	mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(user *entities.User) bool {
		return user.ID == 1 &&
			user.FirstName == "Johnny" &&
			user.LastName == "Smith" &&
			user.Phone == "0987654321"
	})).Return(updatedUser, nil)

	// When
	result, err := useCases.UpdateUser(ctx, 1, request)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Johnny", result.FirstName)
	assert.Equal(t, "Smith", result.LastName)
	assert.Equal(t, "0987654321", result.Phone)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_UpdateUser_UserNotFound(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	request := &dto.UpdateUserRequestDTO{
		FirstName: "Johnny",
	}

	mockRepo.On("GetByID", ctx, uint(999)).Return(nil, domainErrors.ErrUserNotFound)

	// When
	result, err := useCases.UpdateUser(ctx, 999, request)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domainErrors.ErrUserNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_UpdateUser_PartialUpdate(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	existingUser := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Status:    entities.UserStatusActive,
	}

	updatedUser := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "Johnny",     // Only first name updated
		LastName:  "Doe",        // Should remain the same
		Phone:     "1234567890", // Should remain the same
		Status:    entities.UserStatusActive,
	}

	request := &dto.UpdateUserRequestDTO{
		FirstName: "Johnny",
		// LastName and Phone are empty - should not be updated
	}

	mockRepo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.MatchedBy(func(user *entities.User) bool {
		return user.ID == 1 &&
			user.FirstName == "Johnny" &&
			user.LastName == "Doe" && // Should remain unchanged
			user.Phone == "1234567890" // Should remain unchanged
	})).Return(updatedUser, nil)

	// When
	result, err := useCases.UpdateUser(ctx, 1, request)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Johnny", result.FirstName)
	assert.Equal(t, "Doe", result.LastName)     // Should not have changed
	assert.Equal(t, "1234567890", result.Phone) // Should not have changed

	mockRepo.AssertExpectations(t)
}

// ListUsers Tests
func TestUserUseCases_ListUsers_Success(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	expectedUsers := []*entities.User{
		{
			ID:        1,
			Email:     "user1@example.com",
			FirstName: "User",
			LastName:  "One",
			Status:    entities.UserStatusActive,
		},
		{
			ID:        2,
			Email:     "user2@example.com",
			FirstName: "User",
			LastName:  "Two",
			Status:    entities.UserStatusActive,
		},
	}

	mockRepo.On("List", ctx, 10, 0).Return(expectedUsers, nil)

	// When
	result, err := useCases.ListUsers(ctx, 1, 10)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Users, 2)
	assert.Equal(t, 2, result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, 1, result.TotalPages)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_ListUsers_InvalidPagination(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	// Mock for corrected pagination parameters
	mockRepo.On("List", ctx, 10, 0).Return([]*entities.User{}, nil)

	// When - Pass invalid pagination parameters
	result, err := useCases.ListUsers(ctx, -1, 150) // Invalid page and page_size

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Page)      // Should default to 1
	assert.Equal(t, 10, result.PageSize) // Should default to 10

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_ListUsers_SecondPage(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	expectedUsers := []*entities.User{
		{
			ID:        3,
			Email:     "user3@example.com",
			FirstName: "User",
			LastName:  "Three",
		},
	}

	// For page 2 with page_size 5, offset should be 5
	mockRepo.On("List", ctx, 5, 5).Return(expectedUsers, nil)

	// When
	result, err := useCases.ListUsers(ctx, 2, 5)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 5, result.PageSize)
	assert.Len(t, result.Users, 1)

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_ListUsers_RepositoryError(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	mockRepo.On("List", ctx, 10, 0).Return(nil, assert.AnError)

	// When
	result, err := useCases.ListUsers(ctx, 1, 10)

	// Then
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to list users")

	mockRepo.AssertExpectations(t)
}

func TestUserUseCases_ListUsers_EmptyResult(t *testing.T) {
	// Given
	useCases, mockRepo := setupTestUseCases()
	ctx := context.Background()

	mockRepo.On("List", ctx, 10, 0).Return([]*entities.User{}, nil)

	// When
	result, err := useCases.ListUsers(ctx, 1, 10)

	// Then
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Users, 0)
	assert.Equal(t, 0, result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)

	mockRepo.AssertExpectations(t)
}
