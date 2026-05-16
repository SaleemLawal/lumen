package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID               string           `json:"id" example:"1234567890"`
	PlaidItemID      string           `json:"plaid_item_id" example:"1234567890"`
	AccountID        string           `json:"account_id" example:"1234567890"`
	Name             string           `json:"name" example:"Checking Account"`
	Type             string           `json:"type" example:"checking"`
	Subtype          string           `json:"subtype" example:"checking"`
	CurrentBalance   *decimal.Decimal `json:"current_balance" example:"1000.00"`
	AvailableBalance *decimal.Decimal `json:"available_balance" example:"1000.00"`
	CurrencyCode     string           `json:"currency_code" example:"USD"`
	CreatedAt        time.Time        `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt        time.Time        `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}
