package post

import (
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/Wave-95/boards/backend-core/pkg/validator"
	"github.com/google/uuid"
)

// CreatePostInput defines the structure of a request to create a post
type CreatePostInput struct {
	UserID      string  `json:"user_id" validate:"required,uuid"`
	BoardID     string  `json:"board_id" validate:"required,uuid"`
	Content     string  `json:"content"`
	PosX        int     `json:"pos_x"`
	PosY        int     `json:"pos_y"`
	Color       string  `json:"color" validate:"required,min=7,max=7"`
	Height      int     `json:"height" validate:"min=0"`
	ZIndex      int     `json:"z_index"`
	PostOrder   float64 `json:"post_order"`
	PostGroupID string  `json:"post_group_id"`
}

// Validate validates the create post input.
func (i *CreatePostInput) Validate() error {
	validator := validator.New()
	return validator.Struct(i)
}

// UpdatePostInput defines the structure of a request to update a post.
type UpdatePostInput struct {
	ID          string   `json:"id" validate:"required,uuid"`
	Content     *string  `json:"content"`
	Color       *string  `json:"color" validate:"omitempty,min=7,max=7"`
	Height      *int     `json:"height" validate:"omitempty,min=0"`
	PostOrder   *float64 `json:"post_order"`
	PostGroupID *string  `json:"post_group_id" validate:"omitempty,uuid"`
}

// Validate validates the update post payload.
func (i *UpdatePostInput) Validate() error {
	validator := validator.New()
	return validator.Struct(i)
}

// CreatePostgroupInput defines the structure of a request to create a post group.
type CreatePostGroupInput struct {
	BoardID string `json:"board_id" validate:"required,uuid"`
	PosX    int    `json:"pos_x" validate:"required"`
	PosY    int    `json:"pos_y" validate:"required"`
	ZIndex  int    `json:"z_index"`
}

// Validate validates the create post group payload.
func (i CreatePostGroupInput) Validate() error {
	validator := validator.New()
	return validator.Struct(i)
}

// UpdatePostGroupInput defines the structure of a request to update a post group.
type UpdatePostGroupInput struct {
	ID     string  `json:"id" validate:"required,uuid"`
	Title  *string `json:"title"`
	PosX   *int    `json:"pos_x"`
	PosY   *int    `json:"pos_y"`
	ZIndex *int    `json:"z_index"`
}

// Validate validates the update post group payload.
func (i *UpdatePostGroupInput) Validate() error {
	validator := validator.New()
	return validator.Struct(i)
}

// GroupAndPost is a struct that encapsulates data returned from a joined post group and child post.
type GroupAndPost struct {
	PostGroup models.PostGroup
	Post      models.Post
}

// GroupWithPostsDTO is a nested struct describing a post group with associated child posts.
type GroupWithPostsDTO struct {
	ID        uuid.UUID     `json:"id"`
	BoardID   uuid.UUID     `json:"board_id"`
	Title     string        `json:"title"`
	PosX      int           `json:"pos_x"`
	PosY      int           `json:"pos_y"`
	ZIndex    int           `json:"z_index"`
	Posts     []models.Post `json:"posts"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
