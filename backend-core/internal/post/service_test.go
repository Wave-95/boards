package post

import (
	"context"
	"testing"

	"github.com/Wave-95/boards/backend-core/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	mockPostRepo := NewMockRepository()
	service := NewService(mockPostRepo)
	assert.NotNil(t, service)

	t.Run("Create, get, update, and delete post", func(t *testing.T) {
		// Create
		createInput := CreatePostInput{
			UserID:  uuid.New().String(),
			BoardID: uuid.New().String(),
			Content: "This is great content right here",
			PosX:    10,
			PosY:    10,
			Color:   models.PostColorLightPink,
			ZIndex:  1,
		}
		post, err := service.CreatePost(context.Background(), createInput)
		assert.NoError(t, err)
		assert.NotEmpty(t, post.ID)

		// Update
		updatedContent := "This content has been updated"
		updateInput := UpdatePostInput{
			ID:      post.ID.String(),
			Content: &updatedContent,
		}
		_, err = service.UpdatePost(context.Background(), updateInput)
		assert.NoError(t, err)

		// Get
		updatedPost, err := service.GetPost(context.Background(), post.ID.String())
		assert.NoError(t, err)
		assert.Equal(t, updatedContent, updatedPost.Content)

		// Delete post
		err = service.DeletePost(context.Background(), post.ID.String())
		assert.NoError(t, err)
	})
}
