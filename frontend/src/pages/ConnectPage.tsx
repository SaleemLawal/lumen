import { useCallback, useEffect, useState } from "react";
import { usePlaidLink } from "react-plaid-link";
import { Button } from "@/components/ui/button";
import { exchangePublicToken, fetchLinkToken } from "@/src/lib/api";

export default function ConnectPage() {
  const [linkToken, setLinkToken] = useState<string | null>(null);
  const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");
  const [message, setMessage] = useState<string>("");

  useEffect(() => {
    fetchLinkToken()
      .then(setLinkToken)
      .catch((e) => {
        setStatus("error");
        setMessage(e.message);
      });
  }, []);

  const onSuccess = useCallback(async (publicToken: string) => {
    setStatus("loading");
    try {
      await exchangePublicToken(publicToken);
      setStatus("success");
      setMessage("Account linked successfully.");
    } catch (e) {
      setStatus("error");
      setMessage(e instanceof Error ? e.message : "Exchange failed.");
    }
  }, []);

  const { open, ready } = usePlaidLink({
    token: linkToken ?? "",
    onSuccess,
  });

  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-6 bg-background text-foreground">
      <div className="flex flex-col items-center gap-3 text-center">
        <h1 className="text-2xl font-semibold tracking-tight">Connect a bank account</h1>
        <p className="text-sm text-muted-foreground">
          Link your institution via Plaid to test the backend integration.
        </p>
      </div>

      <Button
        size="lg"
        disabled={!ready || status === "loading"}
        onClick={() => open()}
      >
        {status === "loading" ? "Linking…" : "Open Plaid Link"}
      </Button>

      {status === "success" && (
        <p className="text-sm font-medium text-green-600">{message}</p>
      )}
      {status === "error" && (
        <p className="text-sm font-medium text-destructive">{message}</p>
      )}
    </div>
  );
}
