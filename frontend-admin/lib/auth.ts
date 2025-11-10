import { apiClient } from "./api";

export interface User {
  id: string;
  email: string;
}

export interface AuthResponse {
  data: {
    accessToken: string;
    user?: {
      id: string;
      email: string;
      name: string;
      role?: string;
    };
  };
  message: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface InitAdminRequest {
  email: string;
  password: string;
  name: string;
}

/**
 * Get the current admin token from localStorage
 */
export function getAdminToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("admin_token");
}

/**
 * Set the admin token in localStorage
 */
export function setAdminToken(token: string): void {
  if (typeof window === "undefined") return;
  localStorage.setItem("admin_token", token);
}

/**
 * Remove the admin token from localStorage
 */
export function removeAdminToken(): void {
  if (typeof window === "undefined") return;
  localStorage.removeItem("admin_token");
}

/**
 * Get current user information
 */
export async function getCurrentUser(): Promise<User> {
  const response = await apiClient.get<{ data: User }>("/whoami");
  return response.data;
}

/**
 * Start OAuth flow by redirecting to provider
 * For web admin, we need to handle the callback differently
 */
export function startOAuth(provider: "google" | "line"): void {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";
  // Redirect to backend OAuth start endpoint
  // The backend will redirect to Google, then back to the callback
  // We'll handle the callback in our callback page
  window.location.href = `${baseUrl}/auth/${provider}/start`;
}

/**
 * Exchange OAuth code for access token
 */
export async function exchangeOAuthCode(
  provider: "google" | "line",
  code: string,
  state: string
): Promise<AuthResponse> {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";
  const response = await fetch(`${baseUrl}/auth/${provider}/exchange`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ code, state }),
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: "Authentication failed" }));
    throw new Error(error.message || "Authentication failed");
  }

  return (await response.json()) as AuthResponse;
}

/**
 * Check if user is authenticated
 */
export function isAuthenticated(): boolean {
  return getAdminToken() !== null;
}

/**
 * Login with email and password
 */
export async function loginWithPassword(email: string, password: string): Promise<AuthResponse> {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";
  const response = await fetch(`${baseUrl}/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password }),
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: "Login failed" }));
    throw new Error(error.message || "Login failed");
  }

  return (await response.json()) as AuthResponse;
}

/**
 * Initialize admin user
 */
export async function initAdmin(email: string, password: string, name: string): Promise<{ data: { user: User }; message: string }> {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";
  const response = await fetch(`${baseUrl}/auth/init-admin`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email, password, name }),
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: "Failed to initialize admin" }));
    throw new Error(error.message || "Failed to initialize admin");
  }

  return (await response.json()) as { data: { user: User }; message: string };
}

/**
 * Logout user
 */
export function logout(): void {
  removeAdminToken();
  if (typeof window !== "undefined") {
    window.location.href = "/login";
  }
}

