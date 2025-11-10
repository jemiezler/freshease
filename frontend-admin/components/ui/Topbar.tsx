"use client";

import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import { Button } from "@/components/ui/button";
import { SearchIcon, LogOut, User } from "lucide-react";
import { useAuth } from "@/lib/auth-context";

export function Topbar() {
  const { user, logout } = useAuth();

  return (
    <div className="sticky top-0 z-30 border-b bg-white/95 backdrop-blur supports-backdrop-filter:bg-white/70">
      <div className="mx-auto flex h-14 w-full max-w-[1600px] items-center justify-between px-4 sm:px-6">
        <div className="flex items-center justify-between gap-3 w-full">
          <div className="text-sm font-semibold">Freshease Admin</div>
          <InputGroup className="max-w-[300px]">
            <InputGroupInput placeholder="Search" />
            <InputGroupAddon>
              <SearchIcon />
            </InputGroupAddon>
          </InputGroup>
          <div className="flex items-center gap-3">
            {user && (
              <div className="flex items-center gap-2 text-sm text-zinc-600">
                <User className="h-4 w-4" />
                <span className="hidden sm:inline">{user.email}</span>
              </div>
            )}
            <Button
              onClick={logout}
              variant="ghost"
              size="sm"
              className="gap-2"
            >
              <LogOut className="h-4 w-4" />
              <span className="hidden sm:inline">Logout</span>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
