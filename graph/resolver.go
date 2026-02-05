package graph

import "github.com/C-4KE/simple-posts-service/internal/storage"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	storageAccessor *storage.Accessor
}

func NewResolver(accessor *storage.Accessor) *Resolver {
	return &Resolver{
		storageAccessor: accessor,
	}
}
