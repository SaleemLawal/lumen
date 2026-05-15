package store

import "database/sql"

type Storage struct {
	// Holds interface for repository implementations
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{}
}
