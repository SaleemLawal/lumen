import { useCallback, useEffect, useState } from "react";
import { usePlaidLink } from "react-plaid-link";
import { Button } from "@/components/ui/button";
import {
  exchangePublicToken,
  fetchLinkToken,
  fetchLinkedItems,
  fetchUpdateLinkToken,
  syncItemAccounts,
  type LinkedItem,
} from "@/src/lib/api";

export default function ConnectPage() {
  const [linkToken, setLinkToken] = useState<string | null>(null);
  const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
  const [message, setMessage] = useState<string>("");
  const [linkedItems, setLinkedItems] = useState<LinkedItem[]>([]);

  const loadItems = useCallback(() => {
    fetchLinkedItems()
      .then(setLinkedItems)
      .catch(() => {});
  }, []);

  useEffect(() => {
    fetchLinkToken()
      .then(setLinkToken)
      .catch((e) => {
        setStatus("error");
        setMessage(e.message);
      });
    loadItems();
  }, [loadItems]);

  const onSuccess = useCallback(
    async (publicToken: string) => {
      setStatus("loading");
      try {
        await exchangePublicToken(publicToken);
        setStatus("success");
        setMessage("Account linked successfully.");
        loadItems();
      } catch (e) {
        setStatus("error");
        setMessage(e instanceof Error ? e.message : "Exchange failed.");
      }
    },
    [loadItems],
  );

  const { open, ready } = usePlaidLink({
    token: linkToken ?? "",
    onSuccess,
  });

  return (
    <div className="mx-auto flex max-w-lg flex-col gap-10 py-10">
      {/* Link new institution */}
      <section className="flex flex-col items-center gap-4 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">Connect a bank account</h1>
        <p className="text-sm text-muted-foreground">
          Link your institution via Plaid to import accounts and transactions.
        </p>
        <Button size="lg" disabled={!ready || status === "loading"} onClick={() => open()}>
          {status === "loading" ? "Linking…" : "Open Plaid Link"}
        </Button>
        {status === "success" && (
          <p className="text-sm font-medium text-green-600">{message}</p>
        )}
        {status === "error" && (
          <p className="text-sm font-medium text-destructive">{message}</p>
        )}
      </section>

      {/* Manage existing connections */}
      {linkedItems.length > 0 && (
        <section className="flex flex-col gap-3">
          <h2 className="text-base font-semibold tracking-tight">Connected institutions</h2>
          <ul className="flex flex-col gap-2">
            {linkedItems.map((item) => (
              <LinkedItemRow key={item.id} item={item} onSynced={loadItems} />
            ))}
          </ul>
        </section>
      )}
    </div>
  );
}

function LinkedItemRow({ item, onSynced }: { item: LinkedItem; onSynced: () => void }) {
  const [updateToken, setUpdateToken] = useState<string | null>(null);
  const [syncing, setSyncing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const openUpdateLink = useCallback(async () => {
    setError(null);
    try {
      const token = await fetchUpdateLinkToken(item.id);
      setUpdateToken(token);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to open edit flow.");
    }
  }, [item.id]);

  const onSuccess = useCallback(async () => {
    setSyncing(true);
    setError(null);
    try {
      await syncItemAccounts(item.id);
      onSynced();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Sync failed.");
    } finally {
      setSyncing(false);
      setUpdateToken(null);
    }
  }, [item.id, onSynced]);

  const { open, ready } = usePlaidLink({
    token: updateToken ?? "",
    onSuccess,
  });

  useEffect(() => {
    if (updateToken && ready) {
      open();
    }
  }, [updateToken, ready, open]);

  return (
    <li className="flex items-center justify-between rounded-lg border border-border px-4 py-3">
      <div className="flex flex-col gap-0.5">
        <span className="text-sm font-medium">{item.institution_id}</span>
        <span className="text-xs text-muted-foreground">
          {item.accounts.length} account{item.accounts.length !== 1 ? "s" : ""} —{" "}
          {item.accounts.map((a) => a.name).join(", ")}
        </span>
        {error && <span className="text-xs text-destructive">{error}</span>}
      </div>
      <Button
        variant="outline"
        size="sm"
        disabled={syncing || !!updateToken}
        onClick={openUpdateLink}
      >
        {syncing ? "Syncing…" : "Edit accounts"}
      </Button>
    </li>
  );
}
