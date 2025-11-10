"use client";

import { useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { exchangeOAuthCode, setAdminToken, getCurrentUser } from "@/lib/auth";

export default function AuthCallbackPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const code = searchParams.get("code");
    const state = searchParams.get("state");
    
    // Determine provider from URL or default to google
    const provider: "google" | "line" = (searchParams.get("provider") as "google" | "line") || "google";

    if (!code || !state) {
      setError("Missing authorization code or state");
      setLoading(false);
      setTimeout(() => router.push("/login?error=missing_params"), 2000);
      return;
    }

    // Exchange code for token
    exchangeOAuthCode(provider, code, state)
      .then(async (response) => {
        setAdminToken(response.data.accessToken);
        // Verify token by getting user info
        await getCurrentUser();
        // Redirect to home
        router.push("/");
      })
      .catch((err) => {
        console.error("OAuth exchange failed:", err);
        setError(err.message || "Authentication failed");
        setLoading(false);
        setTimeout(() => router.push("/login?error=auth_failed"), 2000);
      });
  }, [searchParams, router]);

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <div className="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-zinc-900 border-r-transparent"></div>
          <p className="text-sm text-zinc-600">Completing authentication...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <p className="text-red-600 mb-4">Error: {error}</p>
          <p className="text-sm text-zinc-600">Redirecting to login...</p>
        </div>
      </div>
    );
  }

  return null;
}

