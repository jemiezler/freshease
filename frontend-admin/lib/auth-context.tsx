"use client";

import { createContext, useContext, useEffect, useState, type ReactNode } from "react";
import { useRouter, usePathname } from "next/navigation";
import type { User } from "./auth";
import {
  setAdminToken,
  removeAdminToken,
  getCurrentUser,
  isAuthenticated as checkIsAuthenticated,
} from "./auth";

interface AuthContextType {
  user: User | null;
  loading: boolean;
  isAuthenticated: boolean;
  login: (token: string) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();
  const pathname = usePathname();

  const isAuthenticated = checkIsAuthenticated();

  const login = async (token: string) => {
    setAdminToken(token);
    // Immediately refresh user after setting token
    await refreshUser();
  };

  const logout = () => {
    removeAdminToken();
    setUser(null);
    router.push("/login");
  };

  const refreshUser = async () => {
    // Check authentication state
    const hasToken = checkIsAuthenticated();
    if (!hasToken) {
      setUser(null);
      setLoading(false);
      return;
    }

    setLoading(true);
    try {
      const userData = await getCurrentUser();
      setUser(userData);
    } catch {
      // Token might be invalid, clear it
      removeAdminToken();
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // Check authentication on mount and route change
    // OAuth callbacks are handled by the dedicated callback page
    // Skip auth check for init-admin page
    if (pathname === "/init-admin") {
      setLoading(false);
      return;
    }
    
    if (isAuthenticated) {
      refreshUser();
    } else {
      setLoading(false);
    }
  }, [pathname, isAuthenticated]);

  // Protect routes - redirect to login if not authenticated
  useEffect(() => {
    // Don't redirect from init-admin or login pages
    if (pathname === "/login" || pathname === "/init-admin") {
      return;
    }
    
    if (!loading && !isAuthenticated) {
      router.push("/login");
    }
  }, [loading, isAuthenticated, pathname, router]);

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        isAuthenticated,
        login,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}

