package store

import (
	"context"
	"database/sql"

	"github.com/saleemlawal/lumen/internal/domain"
)

type PlaidRepository struct {
	db *sql.DB
}

func (r *PlaidRepository) UpsertPlaidItem(ctx context.Context, item *domain.PlaidItem, tx *sql.Tx) error {
	query := `
		INSERT INTO plaid_items (access_token, item_id)
		VALUES ($1, $2)
		ON CONFLICT (item_id) DO UPDATE SET
			access_token = EXCLUDED.access_token,
			updated_at = CURRENT_TIMESTAMP
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, item.AccessToken, item.ItemID)
	return err
}

func (r *PlaidRepository) UpdateCursor(ctx context.Context, itemID, cursor string, tx *sql.Tx) error {
	query := `
		UPDATE plaid_items
		SET transactions_cursor = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE item_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, cursor, itemID)
	return err
}

