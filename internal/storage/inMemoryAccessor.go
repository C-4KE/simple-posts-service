package storage

import (
	"context"
	"errors"
	"maps"
	"slices"
	"strconv"
	"time"

	"github.com/C-4KE/simple-posts-service/graph/model"
	"github.com/google/uuid"
)

type inMemoryStorage struct {
	posts        map[int64]*model.Post
	comments     map[int64]*model.Comment
	commentPaths map[int64]string
}

type InMemoryAccessor struct {
	storage       *inMemoryStorage
	lastPostID    int64
	lastCommentID int64
}

func NewInMemoryAccessor(storage *inMemoryStorage) *InMemoryAccessor {
	return &InMemoryAccessor{
		storage:       storage,
		lastPostID:    -1,
		lastCommentID: -1,
	}
}

func (inMemoryAccessor *InMemoryAccessor) AddPost(ctx context.Context, newPost *model.PostInput) (*model.Post, error) {
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

	inMemoryAccessor.lastPostID++
	post.ID = inMemoryAccessor.lastPostID

	inMemoryAccessor.storage.posts[post.ID] = post

	return post, nil
}

func (inMemoryAccessor *InMemoryAccessor) GetPost(ctx context.Context, postID int64) (*model.Post, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	post, ok := inMemoryAccessor.storage.posts[postID]

	if ok {
		return post, nil
	} else {
		return nil, errors.New("Post with ID " + strconv.Itoa(int(postID)) + " was not found")
	}
}

func (InMemoryAccessor *InMemoryAccessor) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	return slices.Collect(maps.Values(InMemoryAccessor.storage.posts)), nil
}

func (inMemoryAccessor *InMemoryAccessor) UpdateCommentsEnabled(ctx context.Context, postID int64, authorID uuid.UUID, newCommentsEnabled bool) (*model.Post, error) {
	post, ok := inMemoryAccessor.storage.posts[postID]

	if !ok {
		return nil, errors.New("Post with ID " + strconv.Itoa(int(postID)) + " was not found")
	}

	if post.AuthorID != authorID {
		return nil, errors.New("User with ID " + strconv.Itoa(int(authorID.ID())) + " is not the author of the post with ID " + strconv.Itoa(int(postID)) + ".")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	post.CommentsEnabled = newCommentsEnabled
	return post, nil
}
