package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/saleemlawal/lumen/internal/domain"
)

var (
	QUERY_TIMEOUT_DURATION = 5 * time.Second
	ErrRecordNotCreated    = errors.New("Record not created")
)

// Holds interface for repository implementations
type Storage struct {
	Plaid interface {
		Create(context.Context, *domain.PlaidItem) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Plaid: &PlaidRepository{db},
	}
}
