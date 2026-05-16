package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/plaid/plaid-go/v42/plaid"
	"github.com/saleemlawal/lumen/internal/domain"
)

var QUERY_TIMEOUT_DURATION = 5 * time.Second

type Storage struct {
	db *sql.DB

	Plaid interface {
		UpsertPlaidItem(ctx context.Context, item *domain.PlaidItem, tx *sql.Tx) error
		UpdateCursor(ctx context.Context, itemID, cursor string, tx *sql.Tx) error
		InstitutionLinked(ctx context.Context, institutionID string) (bool, error)
		GetAllItems(ctx context.Context) ([]domain.PlaidItemSummary, error)
		GetItemByID(ctx context.Context, id uuid.UUID) (*domain.PlaidItem, error)
	}

	Accounts interface {
		UpsertAccounts(ctx context.Context, itemID string, accounts []plaid.AccountBase, tx *sql.Tx) error
		GetAll(ctx context.Context, itemID *uuid.UUID) ([]domain.Account, error)
	}

	Transactions interface {
		UpsertTransactions(ctx context.Context, itemID string, transactions []plaid.Transaction, tx *sql.Tx) error
		DeleteTransactions(ctx context.Context, removedTransactions []plaid.RemovedTransaction, tx *sql.Tx) error
		GetAll(ctx context.Context, itemID *uuid.UUID, accountID *string) ([]domain.Transaction, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db:           db,
		Plaid:        &PlaidRepository{db},
		Accounts:     &AccountRepository{db},
		Transactions: &TransactionRepository{db},
	}
}

func (s *Storage) StoreLinkSync(
	ctx context.Context,
	item *domain.PlaidItem,
	accounts []plaid.AccountBase,
	added []plaid.Transaction,
	removed []plaid.RemovedTransaction,
	nextCursor string,
) error {
	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Plaid.UpsertPlaidItem(ctx, item, tx); err != nil {
			return err
		}
		if err := s.Accounts.UpsertAccounts(ctx, item.ItemID, accounts, tx); err != nil {
			return err
		}
		if err := s.Transactions.UpsertTransactions(ctx, item.ItemID, added, tx); err != nil {
			return err
		}
		if err := s.Transactions.DeleteTransactions(ctx, removed, tx); err != nil {
			return err
		}
		return s.Plaid.UpdateCursor(ctx, item.ItemID, nextCursor, tx)
	})
}

func (s *Storage) SyncItemTransactions(
	ctx context.Context,
	itemID string,
	added []plaid.Transaction,
	modified []plaid.Transaction,
	removed []plaid.RemovedTransaction,
	nextCursor string,
) error {
	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Transactions.UpsertTransactions(ctx, itemID, added, tx); err != nil {
			return err
		}
		if err := s.Transactions.UpsertTransactions(ctx, itemID, modified, tx); err != nil {
			return err
		}
		if err := s.Transactions.DeleteTransactions(ctx, removed, tx); err != nil {
			return err
		}
		return s.Plaid.UpdateCursor(ctx, itemID, nextCursor, tx)
	})
}

func (s *Storage) SyncItemAccounts(ctx context.Context, itemID string, accounts []plaid.AccountBase) error {
	return WithTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.Accounts.UpsertAccounts(ctx, itemID, accounts, tx)
	})
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}
