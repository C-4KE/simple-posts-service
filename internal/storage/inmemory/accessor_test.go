package inmemory

import (
	"context"
	"testing"

	"github.com/C-4KE/simple-posts-service/graph/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddPost(t *testing.T) {
	mockStorage := NewInMemoryStorage()
	mockAccessor := NewInMemoryAccessor(mockStorage)
	defer mockAccessor.CloseStorage()

	assertions := assert.New(t)

	t.Run("Successful Add Post", func(t *testing.T) {
		ctx := context.Background()
		authorID := uuid.New()

		newPost := &model.PostInput{
			AuthorID:        authorID,
			Title:           "Test Title",
			Text:            "Test Text",
			CommentsEnabled: true,
		}

		createdPost, err := mockAccessor.AddPost(ctx, newPost)
		assertions.Nil(err)
		assertions.Equal(createdPost, &model.Post{
			ID:              0,
			AuthorID:        newPost.AuthorID,
			Title:           newPost.Title,
			Text:            newPost.Text,
			CommentsEnabled: newPost.CommentsEnabled,
			CreateDate:      createdPost.CreateDate,
		})
	})
}
