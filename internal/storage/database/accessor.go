package database

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/C-4KE/simple-posts-service/graph/model"
	"github.com/google/uuid"
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

	err := databaseAccessor.storage.QueryRowContext(ctx, queryInsertPost,
		post.AuthorID,
		post.Title,
		post.Text,
		post.CreateDate,
		post.CommentsEnabled).Scan(&post.ID)

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

	err := databaseAccessor.storage.QueryRowContext(ctx, querySelectPost, postID).Scan(&post.ID,
		&post.AuthorID,
		&post.Title,
		&post.Text,
		&post.CreateDate,
		&post.CommentsEnabled)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (databaseAccessor *DatabaseAccessor) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	posts := make([]*model.Post, 0)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	querySelectPosts := `SELECT post_id, author_id, title, text, create_date, comments_enabled
						FROM posts`

	rows, err := databaseAccessor.storage.QueryContext(ctx, querySelectPosts)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var post model.Post
		if err = rows.Scan(&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Text,
			&post.CreateDate,
			&post.CommentsEnabled); err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	return posts, nil
}

func (databaseAccessor *DatabaseAccessor) UpdateCommentsEnabled(ctx context.Context, postID int64, authorID uuid.UUID, newCommentsEnabled bool) (*model.Post, error) {
	var dbPostId int64
	var dbAuthorID uuid.UUID

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	querySelectPost := `SELECT post_id, author_id
						FROM posts
						WHERE post_id = $1 AND author_id = $2`

	err := databaseAccessor.storage.QueryRowContext(ctx, querySelectPost, postID, authorID).Scan(&dbPostId, &dbAuthorID)

	if err != nil {
		return nil, err
	}

	if dbPostId != postID {
		return nil, errors.New("Post with ID " + strconv.FormatInt(postID, 10) + " was not found")
	}

	if dbAuthorID != authorID {
		return nil, errors.New("User with ID " + strconv.FormatUint(uint64(authorID.ID()), 10) + " is not the author of the post with ID " + strconv.FormatInt(postID, 10) + ".")
	}

	var post model.Post
	queryUpdatePost := `UPDATE posts SET comments_enabled = $1
						WHERE post_id = $2
						RETURNING post_id, author_id, title, text, create_date, comments_enabled`
	err = databaseAccessor.storage.QueryRowContext(ctx, queryUpdatePost, postID, newCommentsEnabled).Scan(&post.ID,
		&post.AuthorID,
		&post.Title,
		&post.Text,
		&post.CreateDate,
		&post.CommentsEnabled)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (databaseAccessor *DatabaseAccessor) AddComment(ctx context.Context, newComment *model.CommentInput) (*model.Comment, error) {
	var commentsEnabled bool

	querySelectPost := `SELECT comments_enabled
						FROM posts
						WHERE post_id = $1`
	err := databaseAccessor.storage.QueryRowContext(ctx, querySelectPost, newComment.PostID).Scan(&commentsEnabled)

	if err != nil {
		return nil, err
	}

	if !commentsEnabled {
		return nil, errors.New("Comments on post " + strconv.FormatInt(newComment.PostID, 10) + " are disabled.")
	}

	comment := &model.Comment{
		AuthorID:   newComment.AuthorID,
		PostID:     newComment.PostID,
		ParentID:   newComment.ParentID,
		Text:       newComment.Text,
		CreateDate: time.Now(),
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	querySelectComment := `SELECT path, replies_level
							FROM comments
							WHERE comment_id = $1`

	var parentPath string
	var parentRepliesLevel int
	err = databaseAccessor.storage.QueryRowContext(ctx, querySelectComment, newComment.ParentID).Scan(&parentPath, &parentRepliesLevel)

	queryInsertComment := `INSERT INTO comments (author_id, post_id, parent_id, text, create_date, path, replies_level)
							VALUES ($1, $2, $3, $4, $5, $6, $7)
							RETURNING comment_id`

	if err == sql.ErrNoRows {
		err = databaseAccessor.storage.QueryRowContext(ctx, queryInsertComment,
			comment.AuthorID,
			comment.PostID,
			comment.ParentID,
			comment.Text,
			comment.CreateDate,
			comment.PostID,
			0).Scan(&comment.ID)
	} else if err == nil {

		err = databaseAccessor.storage.QueryRowContext(ctx, queryInsertComment,
			comment.AuthorID,
			comment.PostID,
			comment.ParentID,
			comment.Text,
			comment.CreateDate,
			strings.Join([]string{parentPath, strconv.FormatInt(*comment.ParentID, 10)}, "."),
			parentRepliesLevel+1).Scan(&comment.ID)
	}

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (databaseAccessor *DatabaseAccessor) GetCommentPath(ctx context.Context, postID int64, parentID *int64) (string, error) {
	var commentsEnabled bool

	querySelectPost := `SELECT post_id
						FROM posts
						WHERE post_id = $1`
	err := databaseAccessor.storage.QueryRowContext(ctx, querySelectPost, postID).Scan(&commentsEnabled)

	if err != nil {
		return "", err
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()

	default:
	}

	querySelectComment := `SELECT path
							FROM comments
							WHERE comment_id = $1`

	var parentPath string
	err = databaseAccessor.storage.QueryRowContext(ctx, querySelectComment, parentID).Scan(&parentPath)

	var path string
	switch err {
	case sql.ErrNoRows:
		path = strconv.FormatInt(postID, 10)
	case nil:
		path = strings.Join([]string{parentPath, strconv.FormatInt(*parentID, 10)}, ".")
	default:
		return "", err
	}

	return path, nil
}

func (databaseAccessor *DatabaseAccessor) GetCommentsLevel(ctx context.Context, postID int64, path string) ([]*model.Comment, error) {
	var commentsEnabled bool

	querySelectPost := `SELECT post_id
						FROM posts
						WHERE post_id = $1`
	err := databaseAccessor.storage.QueryRowContext(ctx, querySelectPost, postID).Scan(&commentsEnabled)

	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	comments := make([]*model.Comment, 0)

	querySelectComments := `SELECT comment_id, author_id, post_id, parent_id, text, create_date
							FROM comments
							WHERE path = $1`

	rows, err := databaseAccessor.storage.QueryContext(ctx, querySelectComments, path)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var comment model.Comment
		if err = rows.Scan(&comment.ID,
			&comment.AuthorID,
			&comment.PostID,
			&comment.ParentID,
			&comment.Text,
			&comment.CreateDate); err != nil {
			return nil, err
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (databaseAccessor *DatabaseAccessor) CloseStorage() {
	databaseAccessor.storage.Close()
}
