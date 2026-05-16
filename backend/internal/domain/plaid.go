package domain

// PlaidItem is persisted metadata for a linked Plaid Item.
type PlaidItem struct {
	ItemID        string
	AccessToken   string
	InstitutionID string
}

type PlaidItemSummary struct {
	ID            string
	InstitutionID string
	Accounts      []AccountSummary
}

type AccountSummary struct {
	AccountID string
	Name      string
	Type      string
	Subtype   string
}
