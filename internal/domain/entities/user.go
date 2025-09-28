package entities

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

type User struct {
	ID        uint       `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"` // Never expose in JSON
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Phone     string     `json:"phone"`
	Status    UserStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Domain methods for business logic
func (u *User) FullName() string {
	return strings.TrimSpace(strings.TrimSpace(u.FirstName) + " " + strings.TrimSpace(u.LastName))
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
}

func (u *User) Suspend() {
	u.Status = UserStatusSuspended
	u.UpdatedAt = time.Now()
}

// Factory function for creating new users
func NewUser(email, password, firstName, lastName, phone string) (*User, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	if strings.TrimSpace(firstName) == "" {
		return nil, errors.New("first name is required")
	}

	now := time.Now()

	return &User{
		Email:     strings.ToLower(strings.TrimSpace(email)),
		Password:  password, // Should be hashed before saving
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Phone:     strings.TrimSpace(phone),
		Status:    UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Domain validation functions
func validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Add more password validation as needed
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}
