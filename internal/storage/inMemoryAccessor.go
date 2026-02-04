package storage

import (
	"github.com/C-4KE/simple-posts-service/graph/model"
)

type inMemoryStorage struct {
	posts        map[int64]*model.Post
	comments     map[int64]*model.Comment
	commentPaths map[int64]string
}

type InMemoryAccessor struct {
	storage *inMemoryStorage
}

func NewInMemoryAccessor(storage *inMemoryStorage) *InMemoryAccessor {
	return &InMemoryAccessor{
		storage: storage,
	}
}
