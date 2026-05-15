const BASE = "/api/v1";

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
  if (!res.ok) throw new Error(`Exchange failed: ${res.status}`);
}
