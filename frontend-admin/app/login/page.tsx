"use client";

import { useEffect, useState, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Field, FieldContent } from "@/components/ui/field";
import { useAuth } from "@/lib/auth-context";
import { startOAuth, loginWithPassword } from "@/lib/auth";
import { Chrome, AlertCircle, Mail, Loader2 } from "lucide-react";

function LoginForm() {
  const { isAuthenticated, loading, login } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const error = searchParams.get("error");

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loginError, setLoginError] = useState("");
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const [showPasswordLogin, setShowPasswordLogin] = useState(false);

  useEffect(() => {
    // Redirect to home if already authenticated
    if (!loading && isAuthenticated) {
      router.push("/");
    }
  }, [isAuthenticated, loading, router]);

  const handleGoogleLogin = () => {
    startOAuth("google");
  };

  const handlePasswordLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoginError("");
    setIsLoggingIn(true);

    try {
      const response = await loginWithPassword(email, password);
      await login(response.data.accessToken);
      // Wait a moment for state to update
      await new Promise((resolve) => setTimeout(resolve, 100));
      router.push("/");
    } catch (err) {
      setLoginError(err instanceof Error ? err.message : "Login failed");
      setIsLoggingIn(false);
    }
  };

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <div className="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-zinc-900 border-r-transparent"></div>
          <p className="text-sm text-zinc-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 px-4">
      <div className="w-full max-w-md space-y-8 rounded-lg border border-zinc-200 bg-white p-8 shadow-sm">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-zinc-900">Freshease Admin</h1>
          <p className="mt-2 text-sm text-zinc-600">Sign in to access the admin panel</p>
        </div>

        {(error === "auth_failed" || loginError) && (
          <div className="rounded-lg border border-red-200 bg-red-50 p-4">
            <div className="flex items-center gap-2 text-sm text-red-800">
              <AlertCircle className="h-4 w-4" />
              <span>{loginError || "Authentication failed. Please try again."}</span>
            </div>
          </div>
        )}

        <div className="space-y-6">
          {!showPasswordLogin ? (
            <>
              {/* Google OAuth Button */}
              <Button
                onClick={handleGoogleLogin}
                className="w-full"
                size="lg"
                variant="outline"
              >
                <Chrome className="mr-2 h-5 w-5" />
                Continue with Google
              </Button>

              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t border-zinc-200" />
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-white px-2 text-zinc-500">Or</span>
                </div>
              </div>

              {/* Email/Password Login Button */}
              <Button
                onClick={() => setShowPasswordLogin(true)}
                className="w-full"
                size="lg"
                variant="outline"
              >
                <Mail className="mr-2 h-5 w-5" />
                Continue with Email
              </Button>
            </>
          ) : (
            <form onSubmit={handlePasswordLogin} className="space-y-4">
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
                    disabled={isLoggingIn}
                  />
                </FieldContent>
              </Field>

              <Field>
                <Label htmlFor="password">Password</Label>
                <FieldContent>
                  <Input
                    id="password"
                    type="password"
                    placeholder="Enter your password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    disabled={isLoggingIn}
                    minLength={8}
                  />
                </FieldContent>
              </Field>

              <div className="flex gap-3">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => {
                    setShowPasswordLogin(false);
                    setLoginError("");
                    setEmail("");
                    setPassword("");
                  }}
                  disabled={isLoggingIn}
                  className="flex-1"
                >
                  Back
                </Button>
                <Button
                  type="submit"
                  disabled={isLoggingIn || !email || !password}
                  className="flex-1"
                >
                  {isLoggingIn ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Signing in...
                    </>
                  ) : (
                    "Sign in"
                  )}
                </Button>
              </div>
            </form>
          )}
        </div>

        <div className="text-center text-xs text-zinc-500">
          By continuing, you agree to our Terms of Service and Privacy Policy
        </div>

        <div className="pt-4 border-t border-zinc-200">
          <p className="text-center text-xs text-zinc-500 mb-2">
            First time setting up?
          </p>
          <Button
            variant="link"
            onClick={() => router.push("/init-admin")}
            className="w-full text-sm"
          >
            Initialize Admin User
          </Button>
        </div>
      </div>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen items-center justify-center">
          <div className="text-center">
            <div className="mb-4 inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-zinc-900 border-r-transparent"></div>
            <p className="text-sm text-zinc-600">Loading...</p>
          </div>
        </div>
      }
    >
      <LoginForm />
    </Suspense>
  );
}

