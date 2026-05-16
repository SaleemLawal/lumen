package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/saleemlawal/lumen/internal/domain"
)

type PlaidRepository struct {
	db *sql.DB
}

func (r *PlaidRepository) UpsertPlaidItem(ctx context.Context, item *domain.PlaidItem, tx *sql.Tx) error {
	query := `
		INSERT INTO plaid_items (access_token, item_id, institution_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (item_id) DO UPDATE SET
			access_token = EXCLUDED.access_token,
			institution_id = EXCLUDED.institution_id,
			updated_at = CURRENT_TIMESTAMP
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, item.AccessToken, item.ItemID, item.InstitutionID)
	return err
}

func (r *PlaidRepository) InstitutionLinked(ctx context.Context, institutionID string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM plaid_items WHERE institution_id = $1)`,
		institutionID,
	).Scan(&exists)
	return exists, err
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

func (r *PlaidRepository) GetAllItems(ctx context.Context) ([]domain.PlaidItemSummary, error) {
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			pi.id,
			pi.institution_id,
			a.account_id,
			a.name,
			a.type,
			a.subtype
		FROM plaid_items pi
		LEFT JOIN accounts a ON a.plaid_item_id = pi.id
		ORDER BY pi.created_at, a.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemMap := make(map[string]*domain.PlaidItemSummary)
	var order []string

	for rows.Next() {
		var itemID, institutionID string
		var accountID, name, typ, subtype sql.NullString

		if err := rows.Scan(&itemID, &institutionID, &accountID, &name, &typ, &subtype); err != nil {
			return nil, err
		}

		if _, seen := itemMap[itemID]; !seen {
			itemMap[itemID] = &domain.PlaidItemSummary{
				ID:            itemID,
				InstitutionID: institutionID,
				Accounts:      []domain.AccountSummary{},
			}
			order = append(order, itemID)
		}

		if accountID.Valid {
			itemMap[itemID].Accounts = append(itemMap[itemID].Accounts, domain.AccountSummary{
				AccountID: accountID.String,
				Name:      name.String,
				Type:      typ.String,
				Subtype:   subtype.String,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	items := make([]domain.PlaidItemSummary, 0, len(order))
	for _, id := range order {
		items = append(items, *itemMap[id])
	}
	return items, nil
}

func (r *PlaidRepository) GetItemByID(ctx context.Context, id uuid.UUID) (*domain.PlaidItem, error) {
	ctx, cancel := context.WithTimeout(ctx, QUERY_TIMEOUT_DURATION)
	defer cancel()

	var item domain.PlaidItem
	err := r.db.QueryRowContext(ctx,
		`SELECT item_id, access_token, institution_id, transactions_cursor FROM plaid_items WHERE id = $1`,
		id,
	).Scan(&item.ItemID, &item.AccessToken, &item.InstitutionID, &item.Cursor)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
