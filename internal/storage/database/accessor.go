package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/C-4KE/simple-posts-service/graph/model"
)

type DatabaseAccessor struct {
	storage *sql.DB
}

func NewDatabaseAccessor(database *sql.DB) *DatabaseAccessor {
	return &DatabaseAccessor{
		storage: database,
	}
}

func (databaseAccessor *DatabaseAccessor) AddPost(ctx context.Context, newPost *model.PostInput) (*model.Post, error) {
	post := &model.Post{
		AuthorID:        newPost.AuthorID,
		Title:           newPost.Title,
		Text:            newPost.Text,
		CommentsEnabled: newPost.CommentsEnabled,
		CreateDate:      time.Now(),
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	queryInsertPost := `INSERT INTO posts (author_id, title, text, create_date, comments_enabled)
						VALUES ($1, $2, $3, $4, $5)
						RETURNING post_id`

	err := databaseAccessor.storage.QueryRowContext(ctx, queryInsertPost, post.AuthorID, post.Title, post.Text, post.CreateDate, post.CommentsEnabled).Scan(&post.ID)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (databaseAccessor *DatabaseAccessor) GetPost(ctx context.Context, postID int64) (*model.Post, error) {
	var post model.Post

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	querySelectPost := `SELECT post_id, author_id, title, text, create_date, comments_enabled
						FROM posts
						WHERE post_id = $1`

	err := databaseAccessor.storage.QueryRowContext(ctx, querySelectPost, postID).Scan(&post.ID, &post.AuthorID, &post.Title, &post.Text, &post.CreateDate, &post.CommentsEnabled)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (databaseAccessor *DatabaseAccessor) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	posts := make([]*model.Post, 10)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	querySelectPost := `SELECT post_id, author_id, title, text, create_date, comments_enabled
						FROM posts`

	rows, err := databaseAccessor.storage.QueryContext(ctx, querySelectPost)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var post model.Post
		if err = rows.Scan(&post.ID, &post.AuthorID, &post.Title, &post.Text, &post.CreateDate, &post.CommentsEnabled); err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}
