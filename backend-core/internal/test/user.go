package test

import (
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

func NewUser() models.User {
	email := uuid.New().String() + "email.com"
	password := "password123"
	user := models.User{
		Id:        uuid.New(),
		Name:      "testname",
		Email:     &email,
		Password:  &password,
		IsGuest:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return user
}