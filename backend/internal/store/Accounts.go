package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/plaid/plaid-go/v42/plaid"
	"github.com/saleemlawal/lumen/internal/domain"
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

func (r *AccountRepository) GetAll(ctx context.Context, itemID *uuid.UUID) ([]domain.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	query := `
		SELECT id, plaid_item_id, account_id, name, type, subtype,
			current_balance, available_balance, currency_code,
			created_at, updated_at
		FROM accounts
		WHERE ($1::uuid IS NULL OR plaid_item_id = $1)
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []domain.Account{}
	for rows.Next() {
		var account domain.Account
		err := rows.Scan(
			&account.ID,
			&account.PlaidItemID,
			&account.AccountID,
			&account.Name,
			&account.Type,
			&account.Subtype,
			&account.CurrentBalance,
			&account.AvailableBalance,
			&account.CurrencyCode,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
