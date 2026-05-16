package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID                 string          `json:"id"`
	PlaidItemID        string          `json:"plaid_item_id"`
	AccountID          string          `json:"account_id"`
	PlaidTransactionID string          `json:"plaid_transaction_id"`
	Name               string          `json:"name"`
	Amount             decimal.Decimal `json:"amount"`
	Date               string          `json:"date"`
	Category           []string        `json:"category"`
	Pending            bool            `json:"pending"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}
