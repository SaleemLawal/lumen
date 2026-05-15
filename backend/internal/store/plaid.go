package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/saleemlawal/lumen/internal/domain"
)

type PlaidRepository struct {
	db *sql.DB
}

func (r *PlaidRepository) Create(ctx context.Context, item *domain.PlaidItem) error {
	query := `
		INSERT INTO plaid_items (access_token, item_id)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, item.AccessToken, item.ItemID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("Failed to create Plaid Item")
	}

	return nil
}
