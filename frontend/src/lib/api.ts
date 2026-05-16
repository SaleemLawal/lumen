const BASE = "/api/v1";

export interface AccountSummary {
  account_id: string;
  name: string;
  type: string;
  subtype: string;
}

export interface LinkedItem {
  id: string;
  institution_id: string;
  accounts: AccountSummary[];
}

export async function fetchLinkToken(): Promise<string> {
  const res = await fetch(`${BASE}/plaid/link-token`);
  if (!res.ok) throw new Error(`Failed to fetch link token: ${res.status}`);
  return res.json();
}

export async function exchangePublicToken(publicToken: string): Promise<void> {
  const res = await fetch(`${BASE}/plaid/exchange-public-token`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ public_token: publicToken }),
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error ?? `Exchange failed: ${res.status}`);
  }
}

export async function fetchLinkedItems(): Promise<LinkedItem[]> {
  const res = await fetch(`${BASE}/plaid/items`);
  if (!res.ok) throw new Error(`Failed to fetch linked items: ${res.status}`);
  return res.json();
}

export async function fetchUpdateLinkToken(itemId: string): Promise<string> {
  const res = await fetch(`${BASE}/plaid/items/${itemId}/link-token`);
  if (!res.ok) throw new Error(`Failed to fetch update link token: ${res.status}`);
  return res.json();
}

export async function syncItemAccounts(itemId: string): Promise<void> {
  const res = await fetch(`${BASE}/plaid/items/${itemId}/sync-accounts`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  });
  if (!res.ok) throw new Error(`Sync failed: ${res.status}`);
}
