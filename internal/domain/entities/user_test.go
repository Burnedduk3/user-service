package entities

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		firstName     string
		lastName      string
		phone         string
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid user creation",
			email:       "john.doe@example.com",
			password:    "SecurePass123",
			firstName:   "John",
			lastName:    "Doe",
			phone:       "1234567890",
			expectError: false,
		},
		{
			name:          "empty email",
			email:         "",
			password:      "SecurePass123",
			firstName:     "John",
			lastName:      "Doe",
			phone:         "1234567890",
			expectError:   true,
			errorContains: "email is required",
		},
		{
			name:          "invalid email format",
			email:         "invalid-email",
			password:      "SecurePass123",
			firstName:     "John",
			lastName:      "Doe",
			phone:         "1234567890",
			expectError:   true,
			errorContains: "invalid email format",
		},
		{
			name:          "password too short",
			email:         "john.doe@example.com",
			password:      "short",
			firstName:     "John",
			lastName:      "Doe",
			phone:         "1234567890",
			expectError:   true,
			errorContains: "password must be at least 8 characters",
		},
		{
			name:          "password missing requirements",
			email:         "john.doe@example.com",
			password:      "alllowercase",
			firstName:     "John",
			lastName:      "Doe",
			phone:         "1234567890",
			expectError:   true,
			errorContains: "password must contain",
		},
		{
			name:          "empty first name",
			email:         "john.doe@example.com",
			password:      "SecurePass123",
			firstName:     "",
			lastName:      "Doe",
			phone:         "1234567890",
			expectError:   true,
			errorContains: "first name is required",
		},
		{
			name:        "empty phone should be valid",
			email:       "john.doe@example.com",
			password:    "SecurePass123",
			firstName:   "John",
			lastName:    "Doe",
			phone:       "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.password, tt.firstName, tt.lastName, tt.phone)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, strings.ToLower(tt.email), user.Email)
				assert.Equal(t, tt.password, user.Password)
				assert.Equal(t, tt.firstName, user.FirstName)
				assert.Equal(t, tt.lastName, user.LastName)
				assert.Equal(t, tt.phone, user.Phone)
				assert.Equal(t, UserStatusActive, user.Status)
				assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
			}
		})
	}
}

func TestUser_FullName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		expected  string
	}{
		{
			name:      "both names provided",
			firstName: "John",
			lastName:  "Doe",
			expected:  "John Doe",
		},
		{
			name:      "only first name",
			firstName: "John",
			lastName:  "",
			expected:  "John",
		},
		{
			name:      "only last name",
			firstName: "",
			lastName:  "Doe",
			expected:  "Doe",
		},
		{
			name:      "names with spaces",
			firstName: " John ",
			lastName:  " Doe ",
			expected:  "John Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				FirstName: tt.firstName,
				LastName:  tt.lastName,
			}
			assert.Equal(t, tt.expected, user.FullName())
		})
	}
}

func TestUser_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   UserStatus
		expected bool
	}{
		{
			name:     "active user",
			status:   UserStatusActive,
			expected: true,
		},
		{
			name:     "inactive user",
			status:   UserStatusInactive,
			expected: false,
		},
		{
			name:     "suspended user",
			status:   UserStatusSuspended,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Status: tt.status}
			assert.Equal(t, tt.expected, user.IsActive())
		})
	}
}

func TestUser_Activate(t *testing.T) {
	user := &User{
		Status:    UserStatusInactive,
		UpdatedAt: time.Now().Add(-time.Hour),
	}
	oldUpdatedAt := user.UpdatedAt

	user.Activate()

	assert.Equal(t, UserStatusActive, user.Status)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt))
}

func TestUser_Suspend(t *testing.T) {
	user := &User{
		Status:    UserStatusActive,
		UpdatedAt: time.Now().Add(-time.Hour),
	}
	oldUpdatedAt := user.UpdatedAt

	user.Suspend()

	assert.Equal(t, UserStatusSuspended, user.Status)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt))
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{"valid email", "test@example.com", false},
		{"valid email with subdomain", "test@mail.example.com", false},
		{"valid email with numbers", "test123@example.com", false},
		{"empty email", "", true},
		{"email without @", "testexample.com", true},
		{"email without domain", "test@", true},
		{"email without TLD", "test@example", true},
		{"email with spaces", "test @example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{"valid password", "SecurePass123", false},
		{"valid password with special chars", "SecurePass123!", false},
		{"too short", "Short1", true},
		{"no uppercase", "securepass123", true},
		{"no lowercase", "SECUREPASS123", true},
		{"no numbers", "SecurePassword", true},
		{"only numbers", "12345678", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
