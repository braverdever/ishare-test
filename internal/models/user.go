package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"unique;not null;size:255"`
	PasswordHash string    `json:"-" gorm:"not null;size:255"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return nil
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// UserResponse represents the response body for user operations
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthorizationCode represents a temporary authorization code for OAuth flow
type AuthorizationCode struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code      string    `json:"code" gorm:"unique;not null;size:255"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ClientID  string    `json:"client_id" gorm:"not null;size:255"`
	Scope     string    `json:"scope" gorm:"size:255"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (code *AuthorizationCode) BeforeCreate(tx *gorm.DB) error {
	if code.ID == uuid.Nil {
		code.ID = uuid.New()
	}
	return nil
}

// AccessToken represents an OAuth access token
type AccessToken struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Token     string    `json:"token" gorm:"unique;not null;size:500"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ClientID  string    `json:"client_id" gorm:"not null;size:255"`
	Scope     string    `json:"scope" gorm:"size:255"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (token *AccessToken) BeforeCreate(tx *gorm.DB) error {
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	return nil
} 