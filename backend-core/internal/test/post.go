package test

import (
	"time"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
)

// NewPost generates a new test post.
func NewPost(userID uuid.UUID, postGroupID uuid.UUID) models.Post {
	postID := uuid.New()
	now := time.Now()
	testPost := models.Post{
		ID:          postID,
		UserID:      userID,
		Content:     "This is a post!",
		Color:       models.PostColorLightPink,
		CreatedAt:   now,
		UpdatedAt:   now,
		PostOrder:   float64(1),
		PostGroupID: postGroupID,
	}
	return testPost
}

// NewPostGroup generates a new test post group.
func NewPostGroup(boardID uuid.UUID) models.PostGroup {
	ID := uuid.New()
	now := time.Now()
	return models.PostGroup{
		ID:        ID,
		BoardID:   boardID,
		PosX:      10,
		PosY:      10,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
