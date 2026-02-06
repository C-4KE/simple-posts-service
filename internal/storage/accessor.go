package storage

import (
	"context"

	"github.com/C-4KE/simple-posts-service/graph/model"
	"github.com/google/uuid"
)

type Accessor interface {
	AddPost(ctx context.Context, newPost *model.PostInput) (*model.Post, error)
	GetPost(ctx context.Context, postID int64) (*model.Post, error)
	GetAllPosts(ctx context.Context) ([]*model.Post, error)
	UpdateCommentsEnabled(ctx context.Context, postID int64, authorID uuid.UUID, newCommentsEnabled bool) (*model.Post, error)

	AddComment(ctx context.Context, newComment *model.CommentInput) (*model.Comment, error)
	GetCommentPath(ctx context.Context, postID int64, parentID *int64) (string, error)
	GetCommentsLevel(ctx context.Context, postID int64, path string) ([]*model.Comment, error)

	CloseStorage()
}
