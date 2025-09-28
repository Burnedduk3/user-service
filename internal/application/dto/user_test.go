// internal/application/dto/user_test.go
package dto

import (
	"encoding/json"
	"testing"
	"time"

	"user-service/internal/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserRequestDTO_ToEntity(t *testing.T) {
	tests := []struct {
		name          string
		dto           CreateUserRequestDTO
		expectError   bool
		errorContains string
	}{
		{
			name: "valid conversion",
			dto: CreateUserRequestDTO{
				Email:     "test@example.com",
				Password:  "SecurePass123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "1234567890",
			},
			expectError: false,
		},
		{
			name: "invalid email",
			dto: CreateUserRequestDTO{
				Email:     "invalid-email",
				Password:  "SecurePass123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "1234567890",
			},
			expectError:   true,
			errorContains: "invalid email format",
		},
		{
			name: "invalid password",
			dto: CreateUserRequestDTO{
				Email:     "test@example.com",
				Password:  "weak",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "1234567890",
			},
			expectError:   true,
			errorContains: "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := tt.dto.ToEntity()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, entity)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, entity)
				assert.Equal(t, tt.dto.Email, entity.Email)
				assert.Equal(t, tt.dto.Password, entity.Password)
				assert.Equal(t, tt.dto.FirstName, entity.FirstName)
				assert.Equal(t, tt.dto.LastName, entity.LastName)
				assert.Equal(t, tt.dto.Phone, entity.Phone)
				assert.Equal(t, entities.UserStatusActive, entity.Status)
			}
		})
	}
}

func TestUserToResponseDTO(t *testing.T) {
	// Given
	now := time.Now()
	user := &entities.User{
		ID:        1,
		Email:     "test@example.com",
		Password:  "hashedpassword", // Should not appear in DTO
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
		Status:    entities.UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// When
	dto := UserToResponseDTO(user)

	// Then
	assert.NotNil(t, dto)
	assert.Equal(t, user.ID, dto.ID)
	assert.Equal(t, user.Email, dto.Email)
	assert.Equal(t, user.FirstName, dto.FirstName)
	assert.Equal(t, user.LastName, dto.LastName)
	assert.Equal(t, "John Doe", dto.FullName)
	assert.Equal(t, user.Phone, dto.Phone)
	assert.Equal(t, user.Status, dto.Status)
	assert.Equal(t, user.CreatedAt, dto.CreatedAt)
	assert.Equal(t, user.UpdatedAt, dto.UpdatedAt)
}

func TestUsersToResponseDTOs(t *testing.T) {
	// Given
	now := time.Now()
	users := []*entities.User{
		{
			ID:        1,
			Email:     "user1@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Status:    entities.UserStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Email:     "user2@example.com",
			FirstName: "Jane",
			LastName:  "Smith",
			Status:    entities.UserStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// When
	dtos := UsersToResponseDTOs(users)

	// Then
	assert.Len(t, dtos, 2)

	assert.Equal(t, users[0].ID, dtos[0].ID)
	assert.Equal(t, users[0].Email, dtos[0].Email)
	assert.Equal(t, "John Doe", dtos[0].FullName)

	assert.Equal(t, users[1].ID, dtos[1].ID)
	assert.Equal(t, users[1].Email, dtos[1].Email)
	assert.Equal(t, "Jane Smith", dtos[1].FullName)
}

func TestCreateUserRequestDTO_JSONSerialization(t *testing.T) {
	// Given
	dto := CreateUserRequestDTO{
		Email:     "test@example.com",
		Password:  "SecurePass123",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "1234567890",
	}

	// When - Serialize to JSON
	jsonData, err := json.Marshal(dto)
	require.NoError(t, err)

	// Then - Deserialize back
	var decodedDTO CreateUserRequestDTO
	err = json.Unmarshal(jsonData, &decodedDTO)
	require.NoError(t, err)

	assert.Equal(t, dto.Email, decodedDTO.Email)
	assert.Equal(t, dto.Password, decodedDTO.Password)
	assert.Equal(t, dto.FirstName, decodedDTO.FirstName)
	assert.Equal(t, dto.LastName, decodedDTO.LastName)
	assert.Equal(t, dto.Phone, decodedDTO.Phone)
}

func TestUserResponseDTO_JSONSerialization(t *testing.T) {
	// Given
	now := time.Now()
	dto := UserResponseDTO{
		ID:        1,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		FullName:  "John Doe",
		Phone:     "1234567890",
		Status:    entities.UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// When - Serialize to JSON
	jsonData, err := json.Marshal(dto)
	require.NoError(t, err)

	// Verify password field is not in JSON
	assert.NotContains(t, string(jsonData), "password")

	// Then - Deserialize back
	var decodedDTO UserResponseDTO
	err = json.Unmarshal(jsonData, &decodedDTO)
	require.NoError(t, err)

	assert.Equal(t, dto.ID, decodedDTO.ID)
	assert.Equal(t, dto.Email, decodedDTO.Email)
	assert.Equal(t, dto.FirstName, decodedDTO.FirstName)
	assert.Equal(t, dto.LastName, decodedDTO.LastName)
	assert.Equal(t, dto.FullName, decodedDTO.FullName)
	assert.Equal(t, dto.Phone, decodedDTO.Phone)
	assert.Equal(t, dto.Status, decodedDTO.Status)
}

func TestUpdateUserRequestDTO_PartialUpdate(t *testing.T) {
	// Given - Only some fields provided
	dto := UpdateUserRequestDTO{
		FirstName: "Johnny",
		// LastName and Phone are omitted
	}

	// When - Serialize to JSON
	jsonData, err := json.Marshal(dto)
	require.NoError(t, err)

	// Then - Verify only non-empty fields are serialized
	var decoded map[string]interface{}
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "Johnny", decoded["first_name"])
	assert.Equal(t, "", decoded["last_name"]) // Empty string, not omitted
	assert.Equal(t, "", decoded["phone"])     // Empty string, not omitted
}

func TestUserListResponseDTO_Structure(t *testing.T) {
	// Given
	users := []*UserResponseDTO{
		{ID: 1, Email: "user1@example.com", FullName: "User One"},
		{ID: 2, Email: "user2@example.com", FullName: "User Two"},
	}

	dto := UserListResponseDTO{
		Users:    users,
		Total:    10,
		Page:     1,
		PageSize: 2,
	}

	// When - Serialize to JSON
	jsonData, err := json.Marshal(dto)
	require.NoError(t, err)

	// Then - Deserialize and verify structure
	var decoded UserListResponseDTO
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Len(t, decoded.Users, 2)
	assert.Equal(t, 10, decoded.Total)
	assert.Equal(t, 1, decoded.Page)
	assert.Equal(t, 2, decoded.PageSize)
}
