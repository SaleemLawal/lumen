package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/plaid/plaid-go/v42/plaid"
	"github.com/saleemlawal/lumen/internal/domain"
)

var QUERY_TIMEOUT_DURATION = 5 * time.Second

// Holds interface for repository implementations
type Storage struct {
	Plaid interface {
		UpsertPlaidItem(context.Context, *domain.PlaidItem) error
	}

	Accounts interface {
		UpsertAccounts(context.Context, string, []plaid.AccountBase) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Plaid:    &PlaidRepository{db},
		Accounts: &AccountRepository{db},
	}
}
