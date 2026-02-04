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

const (
	maxCommentTextLength = 2000
)

type inMemoryStorage struct {
	posts          map[int64]*model.Post
	comments       map[int64]*model.Comment
	commentsByPath map[string][]int64
	commentPaths   map[int64]string
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
		return nil, errors.New("Post with ID " + strconv.FormatInt(postID, 10) + " was not found")
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
		return nil, errors.New("Post with ID " + strconv.FormatInt(postID, 10) + " was not found")
	}

	if post.AuthorID != authorID {
		return nil, errors.New("User with ID " + strconv.FormatUint(uint64(authorID.ID()), 10) + " is not the author of the post with ID " + strconv.FormatInt(postID, 10) + ".")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	post.CommentsEnabled = newCommentsEnabled
	return post, nil
}

func (inMemoryAccessor *InMemoryAccessor) AddComment(ctx context.Context, newComment *model.CommentInput) (*model.Comment, error) {
	_, ok := inMemoryAccessor.storage.posts[newComment.PostID]

	if !ok {
		return nil, errors.New("Post with ID " + strconv.FormatInt(newComment.PostID, 10) + " was not found")
	}

	if len(newComment.Text) > maxCommentTextLength {
		return nil, errors.New("Length of the text ext in the new comment is too big (greater than " + strconv.Itoa(maxCommentTextLength) + ")")
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

	inMemoryAccessor.lastCommentID++
	comment.ID = inMemoryAccessor.lastCommentID

	newCommentPath, err := inMemoryAccessor.GetCommentPath(ctx, newComment.PostID, newComment.ParentID)

	if err != nil {
		return nil, err
	}

	_, ok = inMemoryAccessor.storage.commentsByPath[newCommentPath]
	if !ok {
		inMemoryAccessor.storage.commentsByPath[newCommentPath] = make([]int64, 3)
	}

	inMemoryAccessor.storage.commentsByPath[newCommentPath] = append(inMemoryAccessor.storage.commentsByPath[newCommentPath], comment.ID)
	inMemoryAccessor.storage.commentPaths[comment.ID] = newCommentPath
	inMemoryAccessor.storage.comments[comment.ID] = comment

	return comment, nil
}

func (inMemoryAccessor *InMemoryAccessor) GetCommentPath(ctx context.Context, postID int64, parentID *int64) (string, error) {
	_, ok := inMemoryAccessor.storage.posts[postID]

	if !ok {
		return "", errors.New("Post with ID " + strconv.FormatInt(postID, 10) + " was not found")
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()

	default:
	}

	var commentPath string
	if parentID != nil {
		_, ok := inMemoryAccessor.storage.comments[*parentID]
		if !ok {
			return "", errors.New("Parent comment with ID " + strconv.FormatInt(*parentID, 10) + " was not found")
		}

		commentPath = inMemoryAccessor.storage.commentPaths[*parentID] + "." + strconv.FormatInt(*parentID, 10)
	} else {
		commentPath = strconv.FormatInt(postID, 10)
	}

	return commentPath, nil
}

func (inMemoryAccessor *InMemoryAccessor) GetCommentsLevel(ctx context.Context, postID int64, path string) ([]*model.Comment, error) {
	_, ok := inMemoryAccessor.storage.posts[postID]

	if !ok {
		return nil, errors.New("Post with ID " + strconv.FormatInt(postID, 10) + " was not found")
	}

	commentIDs, ok := inMemoryAccessor.storage.commentsByPath[path]

	if !ok {
		return nil, errors.New("Path is incorrect. Path: " + path)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	comments := make([]*model.Comment, len(commentIDs))
	for idx, commentID := range commentIDs {
		comments[idx] = inMemoryAccessor.storage.comments[commentID]
	}

	return comments, nil
}
