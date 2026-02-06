package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/C-4KE/simple-posts-service/graph/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func getMockAccessor(t *testing.T) (*DatabaseAccessor, sqlmock.Sqlmock) {
	mockStorage, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	mockAccessor := NewDatabaseAccessor(mockStorage)
	return mockAccessor, mock
}

func TestAddPost(t *testing.T) {
	assertions := assert.New(t)
	authorID := uuid.New()
	ctx := context.Background()

	t.Run("Successful Add Post", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		newPost := &model.PostInput{
			AuthorID:        authorID,
			Title:           "Test Title",
			Text:            "Test Text",
			CommentsEnabled: true,
		}

		mock.ExpectQuery(`INSERT INTO posts \(author_id, title, text, create_date, comments_enabled\)
						VALUES \(\$1, \$2, \$3, \$4, \$5\)
						RETURNING post_id`).WithArgs(authorID,
			newPost.Title,
			newPost.Text,
			AnyTime{},
			newPost.CommentsEnabled).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))

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
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		mock.ExpectQuery(`SELECT post_id, author_id
						FROM posts
						WHERE post_id = \$1 AND author_id = \$2`).
			WithArgs(int64(0), authorID).
			WillReturnRows(sqlmock.NewRows([]string{"post_id", "author_id"}).AddRow(int64(0), authorID))

		mock.ExpectQuery(`UPDATE posts SET comments_enabled = \$1
						WHERE post_id = \$2
						RETURNING post_id, author_id, title, text, create_date, comments_enabled`).
			WithArgs(int64(0), false).
			WillReturnRows(sqlmock.
				NewRows([]string{"post_id", "author_id", "title", "text", "create_date", "comments_enabled"}).
				AddRow(int64(0), authorID, "Test Title", "Test Text", time.Now(), false))

		updatedPost, err := mockAccessor.UpdateCommentsEnabled(ctx, 0, authorID, false)
		assertions.Nil(err)
		assertions.NotNil(updatedPost)
		assertions.Equal(updatedPost.CommentsEnabled, false)
	})

	t.Run("Unsuccessful Update CommentsEnabled Incorrect PostID", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		mock.ExpectQuery(`SELECT post_id, author_id
						FROM posts
						WHERE post_id = \$1 AND author_id = \$2`).
			WithArgs(int64(0), authorID)

		err := mock.ExpectationsWereMet()
		assertions.NotNil(err)

		updatedPost, err := mockAccessor.UpdateCommentsEnabled(ctx, -1, authorID, false)
		assertions.NotNil(err)
		assertions.Nil(updatedPost)
	})

	t.Run("Unsuccessful Update CommentsEnabled Incorrect AuthorID", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		incorrectAuthorID := uuid.New()
		mock.ExpectQuery(`SELECT post_id, author_id
						FROM posts
						WHERE post_id = \$1 AND author_id = \$2`).
			WithArgs(int64(0), authorID)

		err := mock.ExpectationsWereMet()
		assertions.NotNil(err)

		updatedPost, err := mockAccessor.UpdateCommentsEnabled(ctx, 0, incorrectAuthorID, false)
		assertions.NotNil(err)
		assertions.Nil(updatedPost)
	})

	t.Run("Successful Get Post", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		existingPost := &model.PostInput{
			AuthorID:        authorID,
			Title:           "Test Title",
			Text:            "Test Text",
			CommentsEnabled: false,
		}

		mock.ExpectQuery(`SELECT post_id, author_id, title, text, create_date, comments_enabled
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(0)).
			WillReturnRows(sqlmock.
				NewRows([]string{"post_id", "author_id", "title", "text", "create_date", "comments_enabled"}).
				AddRow(int64(0), authorID, "Test Title", "Test Text", time.Now(), false))

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
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		mock.ExpectQuery(`SELECT post_id, author_id, title, text, create_date, comments_enabled
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(-1))

		err := mock.ExpectationsWereMet()
		assertions.NotNil(err)

		post, err := mockAccessor.GetPost(ctx, -1)
		assertions.NotNil(err)
		assertions.Nil(post)
	})

	t.Run("Successful Get All Posts", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		mock.ExpectQuery(`SELECT post_id, author_id, title, text, create_date, comments_enabled
						FROM posts`).
			WillReturnRows(sqlmock.
				NewRows([]string{"post_id", "author_id", "title", "text", "create_date", "comments_enabled"}).
				AddRow(int64(0), authorID, "Test Title", "Test Text", time.Now(), false).
				AddRow(int64(1), authorID, "Test Title", "Test Text", time.Now(), true))

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
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   1,
			Text:     "Test Text",
			ParentID: nil,
		}

		mock.ExpectQuery(`SELECT comments_enabled
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"comments_enabled"}).AddRow(true))

		mock.ExpectQuery(`SELECT path, replies_level
							FROM comments
							WHERE comment_id = \$1`).
			WithArgs(nil).WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery(`INSERT INTO comments \(author_id, post_id, parent_id, text, create_date, path, replies_level\)
							VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)
							RETURNING comment_id`).WithArgs(authorID,
			newComment.PostID,
			newComment.ParentID,
			newComment.Text,
			AnyTime{},
			newComment.PostID,
			0).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))

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
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   0,
			Text:     "Test Text",
			ParentID: nil,
		}

		mock.ExpectQuery(`SELECT comments_enabled
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(0)).
			WillReturnRows(sqlmock.NewRows([]string{"comments_enabled"}).AddRow(false))

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.NotNil(err)
		assertions.Nil(createdComment)
	})

	t.Run("Unsuccessful Add Comment Post Does Not Exist", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   -1,
			Text:     "Test Text",
			ParentID: nil,
		}

		mock.ExpectQuery(`SELECT comments_enabled
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(-1)).
			WillReturnError(errors.New("Test"))

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.NotNil(err)
		assertions.Nil(createdComment)
	})

	t.Run("Unsuccessful Add Comment Comments Parent Comment Does Not Exist", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		incorrectParentID := int64(123)
		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   1,
			Text:     "Test Text",
			ParentID: &incorrectParentID,
		}

		mock.ExpectQuery(`SELECT comments_enabled
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"comments_enabled"}).AddRow(true))

		mock.ExpectQuery(`SELECT path, replies_level
							FROM comments
							WHERE comment_id = \$1`).
			WithArgs(nil).WillReturnError(errors.New("Test"))

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.NotNil(err)
		assertions.Nil(createdComment)
	})

	t.Run("Successful Add Child Comment", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		parentID := int64(0)
		newComment := &model.CommentInput{
			AuthorID: authorID,
			PostID:   1,
			Text:     "Test Text",
			ParentID: &parentID,
		}

		mock.ExpectQuery(`SELECT comments_enabled
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"comments_enabled"}).AddRow(true))

		mock.ExpectQuery(`SELECT path, replies_level
							FROM comments
							WHERE comment_id = \$1`).
			WithArgs(newComment.ParentID).
			WillReturnRows(sqlmock.NewRows([]string{"path", "replies_level"}).AddRow("1", 0))

		mock.ExpectQuery(`INSERT INTO comments \(author_id, post_id, parent_id, text, create_date, path, replies_level\)
							VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)
							RETURNING comment_id`).WithArgs(authorID,
			newComment.PostID,
			newComment.ParentID,
			newComment.Text,
			AnyTime{},
			"1.0",
			1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		createdComment, err := mockAccessor.AddComment(ctx, newComment)
		assertions.Nil(err)
		assertions.NotNil(createdComment)
		assertions.Equal(&model.Comment{
			ID:         1,
			AuthorID:   newComment.AuthorID,
			PostID:     newComment.PostID,
			ParentID:   &parentID,
			Text:       newComment.Text,
			CreateDate: createdComment.CreateDate,
		}, createdComment)
	})

	t.Run("Successful Get Root Comment Path", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		mock.ExpectQuery(`SELECT post_id
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"post_id"}).AddRow(int64(1)))

		mock.ExpectQuery(`SELECT path
							FROM comments
							WHERE comment_id = \$1`).
			WithArgs(nil).
			WillReturnError(sql.ErrNoRows)

		commentPath, err := mockAccessor.GetCommentPath(ctx, 1, nil)

		assertions.Nil(err)
		assertions.Equal("1", commentPath)
	})

	t.Run("Successful Get Child Comment Path", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		parentID := int64(0)

		mock.ExpectQuery(`SELECT post_id
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"post_id"}).AddRow(int64(1)))

		mock.ExpectQuery(`SELECT path
							FROM comments
							WHERE comment_id = \$1`).
			WithArgs(parentID).
			WillReturnRows(sqlmock.NewRows([]string{"path"}).AddRow("1"))

		commentPath, err := mockAccessor.GetCommentPath(ctx, 1, &parentID)

		assertions.Nil(err)
		assertions.Equal("1.0", commentPath)
	})

	t.Run("Unsuccessful Get Comment Path Post Does Not Exist", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		mock.ExpectQuery(`SELECT post_id
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(-1)).
			WillReturnError(sql.ErrNoRows)

		commentPath, err := mockAccessor.GetCommentPath(ctx, -1, nil)

		assertions.NotNil(err)
		assertions.Equal("", commentPath)
	})

	t.Run("Unsuccessful Get Comment Path Parent Comment Does Not Exist", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		parentID := int64(123)
		mock.ExpectQuery(`SELECT post_id
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"post_id"}).AddRow(int64(1)))

		mock.ExpectQuery(`SELECT path
							FROM comments
							WHERE comment_id = \$1`).
			WithArgs(parentID).
			WillReturnError(sql.ErrNoRows)

		commentPath, err := mockAccessor.GetCommentPath(ctx, 1, &parentID)

		assertions.NotNil(err)
		assertions.Equal("", commentPath)
	})

	t.Run("Successful Get Root Comments", func(t *testing.T) {
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		mock.ExpectQuery(`SELECT post_id
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"post_id"}).AddRow(int64(1)))

		mock.ExpectQuery(`SELECT comment_id, author_id, post_id, parent_id, text, create_date
							FROM comments
							WHERE path = \$1`).
			WithArgs("1").
			WillReturnRows(sqlmock.
				NewRows([]string{"comment_id", "author_id", "post_id", "parent_id", "text", "create_date"}).
				AddRow(int64(0), authorID, int64(1), nil, "Test Text", time.Now()).
				AddRow(int64(1), authorID, int64(1), nil, "Test Text", time.Now()))

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
		mockAccessor, mock := getMockAccessor(t)
		defer mockAccessor.CloseStorage()

		parentID := int64(0)

		mock.ExpectQuery(`SELECT post_id
						FROM posts
						WHERE post_id = \$1`).
			WithArgs(int64(1)).
			WillReturnRows(sqlmock.NewRows([]string{"post_id"}).AddRow(int64(1)))

		mock.ExpectQuery(`SELECT comment_id, author_id, post_id, parent_id, text, create_date
							FROM comments
							WHERE path = \$1`).
			WithArgs("1.0").
			WillReturnRows(sqlmock.
				NewRows([]string{"comment_id", "author_id", "post_id", "parent_id", "text", "create_date"}).
				AddRow(int64(2), authorID, int64(1), &parentID, "Test Text", time.Now()))

		comments, err := mockAccessor.GetCommentsLevel(ctx, 1, "1.0")

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
