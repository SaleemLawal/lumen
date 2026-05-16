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
		UpsertPlaidItem(ctx context.Context, item *domain.PlaidItem) error
		UpdateCursor(ctx context.Context, itemID, cursor string) error
	}

	Accounts interface {
		UpsertAccounts(ctx context.Context, itemID string, accounts []plaid.AccountBase) error
	}

	Transactions interface {
		UpsertTransactions(ctx context.Context, itemID string, transactions []plaid.Transaction) error
		DeleteTransactions(ctx context.Context, removedTransactions []plaid.RemovedTransaction) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Plaid:        &PlaidRepository{db},
		Accounts:     &AccountRepository{db},
		Transactions: &TransactionRepository{db},
	}
}
