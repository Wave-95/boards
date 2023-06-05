package post

import (
	"context"

	"github.com/Wave-95/boards/server/internal/models"
	"github.com/google/uuid"
)

type mockRepository struct {
	posts map[uuid.UUID]models.Post
}

func NewMockRepository() *mockRepository {
	posts := make(map[uuid.UUID]models.Post)
	return &mockRepository{posts: posts}
}

func (r *mockRepository) CreatePost(ctx context.Context, post models.Post) error {
	r.posts[post.Id] = post
	return nil
}

func (r *mockRepository) GetPost(ctx context.Context, postId uuid.UUID) (models.Post, error) {
	if post, ok := r.posts[postId]; ok {
		return post, nil
	}
	return models.Post{}, ErrPostNotFound
}

func (r *mockRepository) UpdatePost(ctx context.Context, post models.Post) error {
	r.posts[post.Id] = post
	return nil
}

func (r *mockRepository) DeletePost(ctx context.Context, postId uuid.UUID) error {
	delete(r.posts, postId)
	return nil
}