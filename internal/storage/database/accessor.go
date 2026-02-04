package database

import (
	"database/sql"
)

type DatabaseAccessor struct {
	storage *sql.DB
}

func NewDatabaseAccessor(database *sql.DB) *DatabaseAccessor {
	return &DatabaseAccessor{
		storage: database,
	}
}
