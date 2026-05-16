ALTER TABLE plaid_items ADD COLUMN institution_id TEXT;

CREATE UNIQUE INDEX plaid_items_institution_id_key
    ON plaid_items (institution_id)
    WHERE institution_id IS NOT NULL;
