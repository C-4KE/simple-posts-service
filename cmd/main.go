package main

import (
	"errors"
	"flag"
	"log"

	"github.com/C-4KE/simple-posts-service/cmd/dbconnection"
	"github.com/C-4KE/simple-posts-service/cmd/server"
	"github.com/C-4KE/simple-posts-service/internal/storage"
	"github.com/C-4KE/simple-posts-service/internal/storage/database"
	"github.com/C-4KE/simple-posts-service/internal/storage/inmemory"
)

const (
	defaultStorageType = "postgres"
)

func main() {
	storageType := defaultStorageType

	flag.StringVar(&storageType, "storage", "postgress", "Set storage type: 'postgres' ('p') or 'memory' ('m')")
	flag.StringVar(&storageType, "s", "p", "Set storage type: 'postgres' ('p') or 'memory' ('m')")
	flag.Parse()

	switch storageType {
	case "postgress":
	case "memory":
	case "p":
		storageType = "postgres"
	case "m":
		storageType = "memory"
	default:
		log.Printf("Incorrect storage type: %s. %s will be used.", storageType, defaultStorageType)
		storageType = defaultStorageType
	}

	storageAccessor, err := createStorageAccessor(storageType)

	if err != nil {
		log.Fatalf("Error while initializing storage: %s", err)
	}

	server.PostsServer(storageAccessor)
}

func createStorageAccessor(storageType string) (storage.Accessor, error) {
	switch storageType {
	case "postgres":
		return createDatabaseStorage()
	case "memory":
		return createInMemoryStorage()
	default:
		return nil, errors.New("Unsupported storage type: " + storageType)
	}
}

func createInMemoryStorage() (storage.Accessor, error) {
	inMemoryStorage := inmemory.NewInMemoryStorage()
	return inmemory.NewInMemoryAccessor(inMemoryStorage), nil
}

func createDatabaseStorage() (storage.Accessor, error) {
	databaseStorage, err := dbconnection.GetPostgressConnetion()
	return database.NewDatabaseAccessor(databaseStorage), err
}
