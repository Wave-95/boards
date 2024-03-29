package models

import (
	"time"

	"github.com/google/uuid"
)

// User defines a domain model for the user entity.
type User struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      *string   `json:"email"`
	Password   *string   `json:"password,omitempty"`
	IsGuest    bool      `json:"is_guest"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type Verification struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code,omitempty"`
	UserID     uuid.UUID `json:"user_id"`
	IsVerified *bool     `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}
