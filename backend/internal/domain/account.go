package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	ID               string
	PlaidItemID      string
	AccountID        string
	Name             string
	Type             string
	Subtype          string
	CurrentBalance   *decimal.Decimal
	AvailableBalance *decimal.Decimal
	CurrencyCode     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
