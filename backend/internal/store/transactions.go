package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/plaid/plaid-go/v42/plaid"
)

type TransactionRepository struct {
	db *sql.DB
}

func (r *TransactionRepository) UpsertTransactions(ctx context.Context, itemID string, transactions []plaid.Transaction, tx *sql.Tx) error {
	if len(transactions) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	row := tx.QueryRowContext(ctx, `SELECT id FROM plaid_items WHERE item_id = $1`, itemID)
	var plaidItemID uuid.UUID
	if err := row.Scan(&plaidItemID); err != nil {
		return err
	}

	const numCols = 8
	placeholders := buildBatchPlaceholders(len(transactions), numCols)
	args := make([]any, 0, len(transactions)*numCols)

	for _, tx := range transactions {
		args = append(args,
			plaidItemID,
			tx.AccountId,
			tx.TransactionId,
			tx.Name,
			tx.Amount,
			tx.Date,
			pq.Array(tx.Category),
			tx.Pending,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO transactions (plaid_item_id, account_id, plaid_transaction_id, name, amount, date, category, pending)
		VALUES %s
		ON CONFLICT (plaid_transaction_id) DO UPDATE SET
			name = EXCLUDED.name,
			amount = EXCLUDED.amount,
			date = EXCLUDED.date,
			category = EXCLUDED.category,
			pending = EXCLUDED.pending,
			updated_at = CURRENT_TIMESTAMP
	`, strings.Join(placeholders, ","))

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

func (r *TransactionRepository) DeleteTransactions(ctx context.Context, removedTransactions []plaid.RemovedTransaction, tx *sql.Tx) error {
	if len(removedTransactions) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	ids := make([]string, len(removedTransactions))
	for i, tx := range removedTransactions {
		ids[i] = tx.TransactionId
	}

	_, err := tx.ExecContext(ctx,
		`DELETE FROM transactions WHERE plaid_transaction_id = ANY($1)`,
		pq.Array(ids),
	)
	return err
}
