package user

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/security"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/Wave-95/boards/backend-notification/constants/payloads"
	"github.com/Wave-95/boards/backend-notification/constants/queues"
	"github.com/Wave-95/boards/backend-notification/constants/tasks"
	"github.com/Wave-95/boards/wrappers/amqp"
	"github.com/google/uuid"
)

var (
	ErrInvalidVerificationCode = errors.New("Invalid email verification code.")
)

// Service is an interface that describes all the methods pertaining to the user service.
type Service interface {
	CreateUser(ctx context.Context, input CreateUserInput) (models.User, error)
	CreateEmailVerification(ctx context.Context, userID string) (models.Verification, error)

	GetUser(ctx context.Context, userID string) (models.User, error)
	ListUsersByEmail(ctx context.Context, email string) ([]models.User, error)

	VerifyEmail(ctx context.Context, input VerifyEmailInput) error
}

type service struct {
	userRepo  Repository
	amqp      amqp.Amqp
	validator validator.Validate
}

// NewService initializes a service struct with dependencies
func NewService(repo Repository, amqp amqp.Amqp, validator validator.Validate) *service {
	return &service{userRepo: repo, amqp: amqp, validator: validator}
}

// CreateUser takes a user input and standardizes the user name, hashes the password (if provided), and stores the
// user details into the database.
func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (models.User, error) {
	// Validate input
	if err := s.validator.Struct(input); err != nil {
		return models.User{}, fmt.Errorf("service: failed to validate input: %w", err)
	}

	// Prepare user input
	id := uuid.New()
	name := toNameCase(input.Name)
	now := time.Now()
	user := models.User{
		ID:        id,
		Name:      name,
		Email:     input.Email,
		IsGuest:   input.IsGuest,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Hash the password
	if input.Password != nil {
		hashedPassword, err := security.HashPassword(*input.Password)
		if err != nil {
			return models.User{}, fmt.Errorf("service: faled to hash password: %w", err)
		}
		user.Password = &hashedPassword
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return models.User{}, fmt.Errorf("service: failed to create user: %w", err)
	}

	// Publish email verification task
	if input.Email != nil {
		_, err := s.CreateEmailVerification(ctx, id.String())
		if err != nil {
			return models.User{}, fmt.Errorf("service: faled to create email verification: %w", err)
		}
	}

	// Hide password
	user.Password = nil
	return user, nil
}

// CreateEmailVerification validates an input and inserts an email verification record into the database.
func (s *service) CreateEmailVerification(ctx context.Context, userID string) (models.Verification, error) {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return models.Verification{}, fmt.Errorf("service: failed to get user: %w", err)
	}

	// Generate 4 digit code
	randNum := rand.Intn(10000)
	code := fmt.Sprintf("%04d", randNum)

	// Prepare verification entity
	id := uuid.New()
	now := time.Now()
	verification := models.Verification{
		ID:        id,
		Code:      code,
		UserID:    user.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.CreateEmailVerification(ctx, verification); err != nil {
		return models.Verification{}, fmt.Errorf("service: failed to create email verification: %w", err)
	}

	s.amqp.Publish(queues.Notification, tasks.EmailVerification, payloads.EmailVerification{Email: *user.Email, Name: user.Name, Code: code})

	verification.Code = ""
	return verification, nil
}

// GetUser returns a single user for a given user ID
func (s *service) GetUser(ctx context.Context, userID string) (models.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return models.User{}, fmt.Errorf("service: issue parsing userID into UUID: %w", err)
	}

	user, err := s.userRepo.GetUser(ctx, userUUID)
	if err != nil {
		return models.User{}, fmt.Errorf("service: failed to get user: %w", err)
	}

	return user, nil
}

// ListUsersByEmail returns a list of the top 10 users ranked by email similarity
func (s *service) ListUsersByEmail(ctx context.Context, email string) ([]models.User, error) {
	users, err := s.userRepo.ListUsersByFuzzyEmail(ctx, email)
	if err != nil {
		return []models.User{}, fmt.Errorf("service: failed to list users by fuzzy email: %w", err)
	}

	for _, user := range users {
		user.Password = nil
	}

	return users, nil
}

// VerifyEmail looks up the first non-null verification record for a given user ID and compares the input code
// to the verification code. Depending on the match, the users and verification table will be updated.
func (s *service) VerifyEmail(ctx context.Context, input VerifyEmailInput) error {
	// Validate input
	if err := s.validator.Struct(input); err != nil {
		return fmt.Errorf("service: failed to validate input: %w", err)
	}

	userUUID, err := uuid.Parse(input.UserID)
	if err != nil {
		return fmt.Errorf("service: failed to parse user ID into UUID: %w", err)
	}

	verification, err := s.userRepo.GetEmailVerification(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("service: failed to get email verification: %w", err)
	}

	if verification.Code != input.Code {
		return ErrInvalidVerificationCode
	}

	// Update users and verification table
	err = s.userRepo.UpdateUserVerification(ctx, UpdateUserVerificationInput{UserID: userUUID, IsVerified: true})
	if err != nil {
		return fmt.Errorf("service: failed to update user's verification status: %w", err)
	}

	err = s.userRepo.UpdateEmailVerification(ctx, UpdateEmailVerificationInput{UserID: userUUID, IsVerified: true})
	if err != nil {
		return fmt.Errorf("service: failed to update email verification: %w", err)
	}

	return nil
}

// toNameCase creates a regular expression to match word boundaries and convert them to name case
func toNameCase(word string) string {
	re := regexp.MustCompile(`\b\w`)
	nameCase := re.ReplaceAllStringFunc(word, strings.ToUpper)

	return nameCase
}
