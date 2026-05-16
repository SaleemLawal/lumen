package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/plaid/plaid-go/v42/plaid"
)

type AccountRepository struct {
	db *sql.DB
}

func (r *AccountRepository) UpsertAccounts(ctx context.Context, itemID string, accounts []plaid.AccountBase, tx *sql.Tx) error {
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	row := tx.QueryRowContext(ctx, `SELECT id FROM plaid_items WHERE item_id = $1`, itemID)
	var plaidItemID uuid.UUID
	if err := row.Scan(&plaidItemID); err != nil {
		return err
	}

	placeholders := buildBatchPlaceholders(len(accounts), 8)
	args := make([]any, 0, len(accounts)*8)

	for _, account := range accounts {
		args = append(args,
			plaidItemID,
			account.AccountId,
			account.Name,
			account.Type,
			account.Subtype.Get(),
			account.Balances.Current.Get(),
			account.Balances.Available.Get(),
			account.Balances.IsoCurrencyCode.Get(),
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO accounts (plaid_item_id, account_id, name, type, subtype, current_balance, available_balance, currency_code)
		VALUES %s
		ON CONFLICT (account_id) DO UPDATE SET
			name = EXCLUDED.name,
			type = EXCLUDED.type,
			subtype = EXCLUDED.subtype,
			current_balance = EXCLUDED.current_balance,
			available_balance = EXCLUDED.available_balance,
			currency_code = EXCLUDED.currency_code,
			updated_at = CURRENT_TIMESTAMP
	`, strings.Join(placeholders, ","))

	_, err := tx.ExecContext(ctx, query, args...)
	return err
}
