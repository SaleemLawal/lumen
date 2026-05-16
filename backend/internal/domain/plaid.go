package domain

// PlaidItem is persisted metadata for a linked Plaid Item.
type PlaidItem struct {
	ItemID        string `json:"item_id" example:"1234567890"`
	AccessToken   string `json:"access_token" example:"1234567890"`
	InstitutionID string `json:"institution_id" example:"1234567890"`
	Cursor        *string `json:"cursor" example:"s1234567890"`
}

type PlaidItemSummary struct {
	ID            string           `json:"id" example:"1234567890"`
	InstitutionID string           `json:"institution_id" example:"1234567890"`
	Accounts      []AccountSummary `json:"accounts"`
}

type AccountSummary struct {
	AccountID string `json:"account_id" example:"1234567890"`
	Name      string `json:"name" example:"Checking Account"`
	Type      string `json:"type" example:"checking"`
	Subtype   string `json:"subtype" example:"checking"`
}
