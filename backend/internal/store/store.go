package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/saleemlawal/lumen/internal/domain"
)

var QUERY_TIMEOUT_DURATION = 5 * time.Second

// Holds interface for repository implementations
type Storage struct {
	Plaid interface {
		UpsertPlaidItem(context.Context, *domain.PlaidItem) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Plaid: &PlaidRepository{db},
	}
}
