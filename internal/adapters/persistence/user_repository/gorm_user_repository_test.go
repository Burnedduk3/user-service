package user_repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"user-service/internal/application/ports"
	"user-service/internal/domain/entities"
	domainErrors "user-service/internal/domain/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// UserRepositoryTestSuite defines the test suite for user repository
type UserRepositoryTestSuite struct {
	suite.Suite
	repo ports.UserRepository
	ctx  context.Context
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.ctx = context.Background()
	// TODO: Initialize your repository implementation here
	// suite.repo = NewGormUserRepository(db)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	// TODO: Clean up test data
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (suite *UserRepositoryTestSuite) TestCreate_Success() {
	// Given
	user, err := entities.NewUser(
		"test@example.com",
		"SecurePass123",
		"John",
		"Doe",
		"1234567890",
	)
	require.NoError(suite.T(), err)

	// When
	createdUser, err := suite.repo.Create(suite.ctx, user)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), createdUser)
	assert.NotZero(suite.T(), createdUser.ID)
	assert.Equal(suite.T(), user.Email, createdUser.Email)
	assert.Equal(suite.T(), user.FirstName, createdUser.FirstName)
	assert.Equal(suite.T(), user.LastName, createdUser.LastName)
	assert.Equal(suite.T(), user.Phone, createdUser.Phone)
	assert.Equal(suite.T(), entities.UserStatusActive, createdUser.Status)
	assert.WithinDuration(suite.T(), time.Now(), createdUser.CreatedAt, time.Second)
	assert.WithinDuration(suite.T(), time.Now(), createdUser.UpdatedAt, time.Second)
}

func (suite *UserRepositoryTestSuite) TestCreate_DuplicateEmail() {
	// Given
	user1, err := entities.NewUser("duplicate@example.com", "SecurePass123", "John", "Doe", "1234567890")
	require.NoError(suite.T(), err)

	user2, err := entities.NewUser("duplicate@example.com", "SecurePass123", "Jane", "Smith", "0987654321")
	require.NoError(suite.T(), err)

	// When
	_, err = suite.repo.Create(suite.ctx, user1)
	require.NoError(suite.T(), err)

	_, err = suite.repo.Create(suite.ctx, user2)

	// Then
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), domainErrors.ErrUserAlreadyExists.Code, err.(*domainErrors.DomainError).Code)
}

func (suite *UserRepositoryTestSuite) TestGetByID_Success() {
	// Given
	user, err := entities.NewUser("getbyid@example.com", "SecurePass123", "John", "Doe", "1234567890")
	require.NoError(suite.T(), err)

	createdUser, err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// When
	foundUser, err := suite.repo.GetByID(suite.ctx, createdUser.ID)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), createdUser.ID, foundUser.ID)
	assert.Equal(suite.T(), createdUser.Email, foundUser.Email)
	assert.Equal(suite.T(), createdUser.FirstName, foundUser.FirstName)
	assert.Equal(suite.T(), createdUser.LastName, foundUser.LastName)
}

func (suite *UserRepositoryTestSuite) TestGetByID_NotFound() {
	// When
	foundUser, err := suite.repo.GetByID(suite.ctx, 99999)

	// Then
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), foundUser)
	assert.Equal(suite.T(), domainErrors.ErrUserNotFound.Code, err.(*domainErrors.DomainError).Code)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail_Success() {
	// Given
	user, err := entities.NewUser("getbyemail@example.com", "SecurePass123", "John", "Doe", "1234567890")
	require.NoError(suite.T(), err)

	createdUser, err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// When
	foundUser, err := suite.repo.GetByEmail(suite.ctx, "getbyemail@example.com")

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundUser)
	assert.Equal(suite.T(), createdUser.ID, foundUser.ID)
	assert.Equal(suite.T(), createdUser.Email, foundUser.Email)
}

func (suite *UserRepositoryTestSuite) TestGetByEmail_NotFound() {
	// When
	foundUser, err := suite.repo.GetByEmail(suite.ctx, "nonexistent@example.com")

	// Then
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), foundUser)
	assert.Equal(suite.T(), domainErrors.ErrUserNotFound.Code, err.(*domainErrors.DomainError).Code)
}

func (suite *UserRepositoryTestSuite) TestUpdate_Success() {
	// Given
	user, err := entities.NewUser("update@example.com", "SecurePass123", "John", "Doe", "1234567890")
	require.NoError(suite.T(), err)

	createdUser, err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// Modify the user
	createdUser.FirstName = "Johnny"
	createdUser.LastName = "Smith"
	createdUser.Phone = "9876543210"

	// When
	updatedUser, err := suite.repo.Update(suite.ctx, createdUser)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedUser)
	assert.Equal(suite.T(), createdUser.ID, updatedUser.ID)
	assert.Equal(suite.T(), "Johnny", updatedUser.FirstName)
	assert.Equal(suite.T(), "Smith", updatedUser.LastName)
	assert.Equal(suite.T(), "9876543210", updatedUser.Phone)
	assert.True(suite.T(), updatedUser.UpdatedAt.After(updatedUser.CreatedAt))
}

func (suite *UserRepositoryTestSuite) TestUpdate_NotFound() {
	// Given
	user := &entities.User{
		ID:        99999,
		Email:     "notfound@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// When
	updatedUser, err := suite.repo.Update(suite.ctx, user)

	// Then
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), updatedUser)
	assert.Equal(suite.T(), domainErrors.ErrUserNotFound.Code, err.(*domainErrors.DomainError).Code)
}

func (suite *UserRepositoryTestSuite) TestExistsByEmail_True() {
	// Given
	user, err := entities.NewUser("exists@example.com", "SecurePass123", "John", "Doe", "1234567890")
	require.NoError(suite.T(), err)

	_, err = suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// When
	exists, err := suite.repo.ExistsByEmail(suite.ctx, "exists@example.com")

	// Then
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *UserRepositoryTestSuite) TestExistsByEmail_False() {
	// When
	exists, err := suite.repo.ExistsByEmail(suite.ctx, "doesnotexist@example.com")

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *UserRepositoryTestSuite) TestList_Success() {
	// Given - Create multiple users
	users := []*entities.User{}
	for i := 0; i < 5; i++ {
		user, err := entities.NewUser(
			fmt.Sprintf("user%d@example.com", i),
			"SecurePass123",
			fmt.Sprintf("User%d", i),
			"Test",
			"1234567890",
		)
		require.NoError(suite.T(), err)

		createdUser, err := suite.repo.Create(suite.ctx, user)
		require.NoError(suite.T(), err)
		users = append(users, createdUser)
	}

	// When
	foundUsers, err := suite.repo.List(suite.ctx, 3, 0)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), foundUsers, 3)
}

func (suite *UserRepositoryTestSuite) TestList_WithPagination() {
	// Given - Create multiple users (assuming they exist from previous test)

	// When - Get second page
	foundUsers, err := suite.repo.List(suite.ctx, 2, 2)

	// Then
	assert.NoError(suite.T(), err)
	assert.LessOrEqual(suite.T(), len(foundUsers), 2)
}

func (suite *UserRepositoryTestSuite) TestDelete_Success() {
	// Given
	user, err := entities.NewUser("delete@example.com", "SecurePass123", "John", "Doe", "1234567890")
	require.NoError(suite.T(), err)

	createdUser, err := suite.repo.Create(suite.ctx, user)
	require.NoError(suite.T(), err)

	// When
	err = suite.repo.Delete(suite.ctx, createdUser.ID)

	// Then
	assert.NoError(suite.T(), err)

	// Verify user is deleted
	_, err = suite.repo.GetByID(suite.ctx, createdUser.ID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), domainErrors.ErrUserNotFound.Code, err.(*domainErrors.DomainError).Code)
}

func (suite *UserRepositoryTestSuite) TestDelete_NotFound() {
	// When
	err := suite.repo.Delete(suite.ctx, 99999)

	// Then
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), domainErrors.ErrUserNotFound.Code, err.(*domainErrors.DomainError).Code)
}
