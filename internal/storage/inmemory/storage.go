package inmemory

import "github.com/C-4KE/simple-posts-service/graph/model"

type InMemoryStorage struct {
	posts          map[int64]*model.Post
	comments       map[int64]*model.Comment
	commentsByPath map[string][]int64
	commentPaths   map[int64]string
}
