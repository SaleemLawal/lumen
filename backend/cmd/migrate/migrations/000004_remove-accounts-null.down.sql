UPDATE accounts SET current_balance = 0 WHERE current_balance IS NULL;
UPDATE accounts SET available_balance = 0 WHERE available_balance IS NULL;
UPDATE accounts SET currency_code = '' WHERE currency_code IS NULL;
UPDATE accounts SET subtype = '' WHERE subtype IS NULL;

ALTER TABLE accounts
ALTER COLUMN current_balance SET NOT NULL;
ALTER TABLE accounts
ALTER COLUMN available_balance SET NOT NULL;
ALTER TABLE accounts
ALTER COLUMN currency_code SET NOT NULL;
ALTER TABLE accounts
ALTER COLUMN subtype SET NOT NULL;