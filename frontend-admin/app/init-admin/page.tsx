"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Field, FieldContent } from "@/components/ui/field";
import { initAdmin } from "@/lib/auth";
import { AlertCircle, CheckCircle, Loader2 } from "lucide-react";

export default function InitAdminPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [name, setName] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const [isInitializing, setIsInitializing] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    // Validate passwords match
    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    // Validate password length
    if (password.length < 8) {
      setError("Password must be at least 8 characters long");
      return;
    }

    setIsInitializing(true);

    try {
      await initAdmin(email, password, name);
      setSuccess(true);
      // Redirect to login after 2 seconds
      setTimeout(() => {
        router.push("/login");
      }, 2000);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to initialize admin");
    } finally {
      setIsInitializing(false);
    }
  };

  if (success) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-zinc-50 px-4">
        <div className="w-full max-w-md space-y-6 rounded-lg border border-green-200 bg-green-50 p-8 shadow-sm">
          <div className="flex items-center justify-center">
            <CheckCircle className="h-12 w-12 text-green-600" />
          </div>
          <div className="text-center">
            <h2 className="text-2xl font-bold text-green-900">Admin Initialized!</h2>
            <p className="mt-2 text-sm text-green-700">
              Admin user has been created successfully. Redirecting to login...
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 px-4">
      <div className="w-full max-w-md space-y-8 rounded-lg border border-zinc-200 bg-white p-8 shadow-sm">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-zinc-900">Initialize Admin</h1>
          <p className="mt-2 text-sm text-zinc-600">
            Create the first admin user for the Freshease admin panel
          </p>
        </div>

        {error && (
          <div className="rounded-lg border border-red-200 bg-red-50 p-4">
            <div className="flex items-center gap-2 text-sm text-red-800">
              <AlertCircle className="h-4 w-4" />
              <span>{error}</span>
            </div>
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <Field>
            <Label htmlFor="name">Name</Label>
            <FieldContent>
              <Input
                id="name"
                type="text"
                placeholder="Admin Name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                required
                disabled={isInitializing}
                minLength={2}
              />
            </FieldContent>
          </Field>

          <Field>
            <Label htmlFor="email">Email</Label>
            <FieldContent>
              <Input
                id="email"
                type="email"
                placeholder="admin@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                disabled={isInitializing}
              />
            </FieldContent>
          </Field>

          <Field>
            <Label htmlFor="password">Password</Label>
            <FieldContent>
              <Input
                id="password"
                type="password"
                placeholder="Minimum 8 characters"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                disabled={isInitializing}
                minLength={8}
              />
            </FieldContent>
          </Field>

          <Field>
            <Label htmlFor="confirmPassword">Confirm Password</Label>
            <FieldContent>
              <Input
                id="confirmPassword"
                type="password"
                placeholder="Confirm your password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
                disabled={isInitializing}
                minLength={8}
              />
            </FieldContent>
          </Field>

          <Button
            type="submit"
            disabled={isInitializing || !email || !password || !name || !confirmPassword}
            className="w-full"
            size="lg"
          >
            {isInitializing ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Initializing...
              </>
            ) : (
              "Initialize Admin"
            )}
          </Button>
        </form>

        <div className="text-center">
          <Button
            variant="ghost"
            onClick={() => router.push("/login")}
            disabled={isInitializing}
            className="text-sm"
          >
            Back to Login
          </Button>
        </div>
      </div>
    </div>
  );
}

