package inmemory

import (
	"github.com/C-4KE/simple-posts-service/graph/model"
	"github.com/C-4KE/simple-posts-service/internal/helpers"
)

type InMemoryStorage struct {
	posts          *helpers.SafeMap[int64, *model.Post]
	comments       *helpers.SafeMap[int64, *model.Comment]
	commentsByPath *helpers.SafeMap[string, []int64]
	commentPaths   *helpers.SafeMap[int64, string]
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		posts:          helpers.NewSafeMap(make(map[int64]*model.Post)),
		comments:       helpers.NewSafeMap(make(map[int64]*model.Comment)),
		commentsByPath: helpers.NewSafeMap(make(map[string][]int64)),
		commentPaths:   helpers.NewSafeMap(make(map[int64]string)),
	}
}
