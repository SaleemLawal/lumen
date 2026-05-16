package domain

// PlaidItem is persisted metadata for a linked Plaid Item.
type PlaidItem struct {
	ItemID        string
	AccessToken   string
	InstitutionID string
}
