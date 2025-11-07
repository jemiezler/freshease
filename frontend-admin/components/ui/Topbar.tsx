"use client";

import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import { SearchIcon } from "lucide-react";

export function Topbar() {
  return (
    <div className="sticky top-0 z-30 border-b bg-white/95 backdrop-blur supports-backdrop-filter:bg-white/70">
      <div className="mx-auto flex h-14 w-full max-w-[1600px] items-center justify-between px-4 sm:px-6">
        <div className="flex items-center justify-between gap-3 w-full">
          <div className="text-sm font-semibold">Freshease</div>
          <InputGroup className="max-w-[300px]">
            <InputGroupInput placeholder="Search" />
            <InputGroupAddon>
              <SearchIcon />
            </InputGroupAddon>
          </InputGroup>
        </div>
      </div>
    </div>
  );
}
