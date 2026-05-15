CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plaid_item_id UUID NOT NULL REFERENCES plaid_items(id),
    account_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    subtype TEXT NOT NULL,
    current_balance DECIMAL(12, 2) NOT NULL,
    available_balance DECIMAL(12, 2) NOT NULL,
    currency_code TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);