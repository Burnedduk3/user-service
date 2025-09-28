package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/internal/application/dto"
	"user-service/internal/domain/entities"
	domainErrors "user-service/internal/domain/errors"
	"user-service/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserUseCases implements the UserUseCases interface for testing
type MockUserUseCases struct {
	mock.Mock
}

func (m *MockUserUseCases) CreateUser(ctx context.Context, request *dto.CreateUserRequestDTO) (*dto.UserResponseDTO, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUseCases) GetUserByID(ctx context.Context, id uint) (*dto.UserResponseDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUseCases) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponseDTO, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUseCases) UpdateUser(ctx context.Context, id uint, request *dto.UpdateUserRequestDTO) (*dto.UserResponseDTO, error) {
	args := m.Called(ctx, id, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponseDTO), args.Error(1)
}

func (m *MockUserUseCases) ListUsers(ctx context.Context, page, pageSize int) (*dto.UserListResponseDTO, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserListResponseDTO), args.Error(1)
}

func setupTestHandler() (*UserHandler, *MockUserUseCases) {
	mockUseCases := new(MockUserUseCases)
	log := logger.New("test")
	handler := NewUserHandler(mockUseCases, log)
	return handler, mockUseCases
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	// Setup
	handler, mockUseCases := setupTestHandler()

	requestBody := dto.CreateUserRequestDTO{
		Email:     "test@example.com",
		Password:  "SecurePass123",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
	}

	expectedResponse := &dto.UserResponseDTO{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		FullName:  "John Doe",
		Phone:     "1234567890",
		Status:    entities.UserStatusActive,
	}

	mockUseCases.On("CreateUser", mock.Anything, &requestBody).Return(expectedResponse, nil)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Execute
	err := handler.CreateUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.UserResponseDTO
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.Email, response.Email)
	assert.Equal(t, expectedResponse.FirstName, response.FirstName)

	mockUseCases.AssertExpectations(t)
}

func TestUserHandler_CreateUser_ValidationError(t *testing.T) {
	// Setup
	handler, _ := setupTestHandler()

	requestBody := dto.CreateUserRequestDTO{
		Email:     "invalid-email", // Invalid email format
		Password:  "short",         // Too short
		FirstName: "",              // Required field missing
		LastName:  "Doe",
	}

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Execute
	err := handler.CreateUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "VALIDATION_ERROR", response.Error)
	assert.NotNil(t, response.Details)
}

func TestUserHandler_CreateUser_UserAlreadyExists(t *testing.T) {
	// Setup
	handler, mockUseCases := setupTestHandler()

	requestBody := dto.CreateUserRequestDTO{
		Email:     "existing@example.com",
		Password:  "SecurePass123",
		FirstName: "John",
		LastName:  "Doe",
	}

	mockUseCases.On("CreateUser", mock.Anything, &requestBody).Return(nil, domainErrors.ErrUserAlreadyExists)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Execute
	err := handler.CreateUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)

	var response ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "USER_ALREADY_EXISTS", response.Error)
	mockUseCases.AssertExpectations(t)
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	// Setup
	handler, mockUseCases := setupTestHandler()

	expectedResponse := &dto.UserResponseDTO{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		FullName:  "John Doe",
		Status:    entities.UserStatusActive,
	}

	mockUseCases.On("GetUserByID", mock.Anything, uint(1)).Return(expectedResponse, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/1", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Execute
	err := handler.GetUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.UserResponseDTO
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.Email, response.Email)

	mockUseCases.AssertExpectations(t)
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	// Setup
	handler, mockUseCases := setupTestHandler()

	mockUseCases.On("GetUserByID", mock.Anything, uint(999)).Return(nil, domainErrors.ErrUserNotFound)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/999", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	// Execute
	err := handler.GetUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var response ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "USER_NOT_FOUND", response.Error)
	mockUseCases.AssertExpectations(t)
}

func TestUserHandler_GetUser_InvalidID(t *testing.T) {
	// Setup
	handler, _ := setupTestHandler()

	// Create request with invalid ID
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/invalid", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	// Execute
	err := handler.GetUser(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "INVALID_ID", response.Error)
}

func TestUserHandler_ListUsers_Success(t *testing.T) {
	// Setup
	handler, mockUseCases := setupTestHandler()

	expectedUsers := []*dto.UserResponseDTO{
		{
			ID:        1,
			Email:     "user1@example.com",
			FirstName: "User",
			LastName:  "One",
			FullName:  "User One",
			Status:    entities.UserStatusActive,
		},
		{
			ID:        2,
			Email:     "user2@example.com",
			FirstName: "User",
			LastName:  "Two",
			FullName:  "User Two",
			Status:    entities.UserStatusActive,
		},
	}

	expectedResponse := &dto.UserListResponseDTO{
		Users:      expectedUsers,
		Total:      2,
		Page:       1,
		PageSize:   10,
		TotalPages: 1,
	}

	mockUseCases.On("ListUsers", mock.Anything, 1, 10).Return(expectedResponse, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Execute
	err := handler.ListUsers(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.UserListResponseDTO
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response.Users, 2)
	assert.Equal(t, 2, response.Total)
	assert.Equal(t, 1, response.Page)

	mockUseCases.AssertExpectations(t)
}

func TestUserHandler_ListUsers_WithPagination(t *testing.T) {
	// Setup
	handler, mockUseCases := setupTestHandler()

	expectedResponse := &dto.UserListResponseDTO{
		Users:      []*dto.UserResponseDTO{},
		Total:      0,
		Page:       2,
		PageSize:   5,
		TotalPages: 1,
	}

	mockUseCases.On("ListUsers", mock.Anything, 2, 5).Return(expectedResponse, nil)

	// Create request with pagination parameters
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?page=2&page_size=5", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	// Execute
	err := handler.ListUsers(c)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockUseCases.AssertExpectations(t)
}
