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
	authorID := uuid.New()
	ctx := context.Background()

	t.Run("Successful Add Post", func(t *testing.T) {
		newPost := &model.PostInput{
			AuthorID:        authorID,
			Title:           "Test Title",
			Text:            "Test Text",
			CommentsEnabled: true,
		}

		createdPost, err := mockAccessor.AddPost(ctx, newPost)
		assertions.Nil(err)
		assertions.NotNil(createdPost)
		assertions.Equal(&model.Post{
			ID:              0,
			AuthorID:        newPost.AuthorID,
			Title:           newPost.Title,
			Text:            newPost.Text,
			CommentsEnabled: newPost.CommentsEnabled,
			CreateDate:      createdPost.CreateDate,
		}, createdPost)
	})

	t.Run("Successful Update CommentsEnabled", func(t *testing.T) {
		updatedPost, err := mockAccessor.UpdateCommentsEnabled(ctx, 0, authorID, false)
		assertions.Nil(err)
		assertions.NotNil(updatedPost)
		assertions.Equal(updatedPost.CommentsEnabled, false)
	})

	t.Run("Unsuccessful Update CommentsEnabled Incorrect PostID", func(t *testing.T) {
		updatedPost, err := mockAccessor.UpdateCommentsEnabled(ctx, -1, authorID, false)
		assertions.NotNil(err)
		assertions.Nil(updatedPost)
	})

	t.Run("Unsuccessful Update CommentsEnabled Incorrect AuthorID", func(t *testing.T) {
		updatedPost, err := mockAccessor.UpdateCommentsEnabled(ctx, 0, uuid.New(), false)
		assertions.NotNil(err)
		assertions.Nil(updatedPost)
	})

	t.Run("Successful Get Post", func(t *testing.T) {
		existingPost := &model.PostInput{
			AuthorID:        authorID,
			Title:           "Test Title",
			Text:            "Test Text",
			CommentsEnabled: false,
		}

		post, err := mockAccessor.GetPost(ctx, 0)
		assertions.Nil(err)
		assertions.Equal(&model.Post{
			ID:              0,
			AuthorID:        existingPost.AuthorID,
			Title:           existingPost.Title,
			Text:            existingPost.Text,
			CommentsEnabled: existingPost.CommentsEnabled,
			CreateDate:      post.CreateDate,
		}, post)
	})

	t.Run("Unsuccessful Get Post Incorrect PostID", func(t *testing.T) {
		post, err := mockAccessor.GetPost(ctx, -1)
		assertions.NotNil(err)
		assertions.Nil(post)
	})

	t.Run("Successful Another Post", func(t *testing.T) {
		newPost := &model.PostInput{
			AuthorID:        authorID,
			Title:           "Test Title",
			Text:            "Test Text",
			CommentsEnabled: true,
		}

		createdPost, err := mockAccessor.AddPost(ctx, newPost)
		assertions.Nil(err)
		assertions.NotNil(createdPost)
		assertions.Equal(&model.Post{
			ID:              1,
			AuthorID:        newPost.AuthorID,
			Title:           newPost.Title,
			Text:            newPost.Text,
			CommentsEnabled: newPost.CommentsEnabled,
			CreateDate:      createdPost.CreateDate,
		}, createdPost)
	})

	t.Run("Successful Get All Posts", func(t *testing.T) {
		posts, err := mockAccessor.GetAllPosts(ctx)
		assertions.Nil(err)
		assertions.Equal([]*model.Post{
			{
				ID:              0,
				AuthorID:        authorID,
				Title:           "Test Title",
				Text:            "Test Text",
				CommentsEnabled: false,
				CreateDate:      posts[0].CreateDate,
			},
			{
				ID:              1,
				AuthorID:        authorID,
				Title:           "Test Title",
				Text:            "Test Text",
				CommentsEnabled: true,
				CreateDate:      posts[1].CreateDate,
			},
		}, posts)
	})

	t.Run("Successful Add Comment", func(t *testing.T) {
		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   1,
			Text:     "Test Text",
			ParentID: nil,
		}

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.Nil(err)
		assertions.NotNil(createdComment)
		assertions.Equal(&model.Comment{
			ID:         0,
			AuthorID:   newComment.AuthorID,
			PostID:     newComment.PostID,
			ParentID:   nil,
			Text:       newComment.Text,
			CreateDate: createdComment.CreateDate,
		}, createdComment)
	})

	t.Run("Unsuccessful Add Comment Comments Disabled", func(t *testing.T) {
		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   0,
			Text:     "Test Text",
			ParentID: nil,
		}

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.NotNil(err)
		assertions.Nil(createdComment)
	})

	t.Run("Unsuccessful Add Comment Post Does Not Exist", func(t *testing.T) {
		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   -1,
			Text:     "Test Text",
			ParentID: nil,
		}

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.NotNil(err)
		assertions.Nil(createdComment)
	})

	t.Run("Unsuccessful Add Comment Comments Parent Comment Does Not Exist", func(t *testing.T) {
		incorrectParentID := int64(123)
		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   1,
			Text:     "Test Text",
			ParentID: &incorrectParentID,
		}

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.NotNil(err)
		assertions.Nil(createdComment)
	})

	t.Run("Successful Add Another Comment", func(t *testing.T) {
		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   1,
			Text:     "Test Text",
			ParentID: nil,
		}

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.Nil(err)
		assertions.NotNil(createdComment)
		assertions.Equal(&model.Comment{
			ID:         1,
			AuthorID:   newComment.AuthorID,
			PostID:     newComment.PostID,
			ParentID:   nil,
			Text:       newComment.Text,
			CreateDate: createdComment.CreateDate,
		}, createdComment)
	})

	t.Run("Successful Add Child Comment", func(t *testing.T) {
		parentID := int64(0)
		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   1,
			Text:     "Test Text",
			ParentID: &parentID,
		}

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.Nil(err)
		assertions.NotNil(createdComment)
		assertions.Equal(&model.Comment{
			ID:         2,
			AuthorID:   newComment.AuthorID,
			PostID:     newComment.PostID,
			ParentID:   &parentID,
			Text:       newComment.Text,
			CreateDate: createdComment.CreateDate,
		}, createdComment)
	})

	t.Run("Successful Get Root Comment Path", func(t *testing.T) {
		commentPath, err := mockAccessor.GetCommentPath(ctx, 1, nil)

		assertions.Nil(err)
		assertions.Equal("1", commentPath)
	})

	t.Run("Successful Get Child Comment Path", func(t *testing.T) {
		parentID := int64(0)
		commentPath, err := mockAccessor.GetCommentPath(ctx, 1, &parentID)

		assertions.Nil(err)
		assertions.Equal("1.0", commentPath)
	})

	t.Run("Unsuccessful Get Comment Path Post Does Not Exist", func(t *testing.T) {
		commentPath, err := mockAccessor.GetCommentPath(ctx, -1, nil)

		assertions.NotNil(err)
		assertions.Equal("", commentPath)
	})

	t.Run("Unsuccessful Get Comment Path Parent Comment Does Not Exist", func(t *testing.T) {
		parentID := int64(123)
		commentPath, err := mockAccessor.GetCommentPath(ctx, 1, &parentID)

		assertions.NotNil(err)
		assertions.Equal("", commentPath)
	})

	t.Run("Successful Get Root Comments", func(t *testing.T) {
		comments, err := mockAccessor.GetCommentsLevel(ctx, 1, "1")

		assertions.Nil(err)
		assertions.Equal([]*model.Comment{
			{
				ID:         0,
				AuthorID:   authorID,
				PostID:     1,
				ParentID:   nil,
				Text:       "Test Text",
				CreateDate: comments[0].CreateDate,
			},
			{
				ID:         1,
				AuthorID:   authorID,
				PostID:     1,
				ParentID:   nil,
				Text:       "Test Text",
				CreateDate: comments[1].CreateDate,
			},
		}, comments)
	})

	t.Run("Successful Get Child Comments", func(t *testing.T) {
		comments, err := mockAccessor.GetCommentsLevel(ctx, 1, "1.0")

		parentID := int64(0)
		assertions.Nil(err)
		assertions.Equal([]*model.Comment{
			{
				ID:         2,
				AuthorID:   authorID,
				PostID:     1,
				ParentID:   &parentID,
				Text:       "Test Text",
				CreateDate: comments[0].CreateDate,
			},
		}, comments)
	})
}
