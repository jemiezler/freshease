"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutDashboard, Box, Tag, Boxes, Store, Users, Shield, Key, ShoppingCart, ListChecks, MapPin } from "lucide-react";

const links = [
	{ href: "/", label: "Dashboard", icon: LayoutDashboard },
	{ href: "/products", label: "Products", icon: Box },
	{ href: "/categories", label: "Categories", icon: Tag },
	{ href: "/inventories", label: "Inventories", icon: Boxes },
	{ href: "/vendors", label: "Vendors", icon: Store },
	{ href: "/users", label: "Users", icon: Users },
	{ href: "/roles", label: "Roles", icon: Shield },
	{ href: "/permissions", label: "Permissions", icon: Key },
	{ href: "/carts", label: "Carts", icon: ShoppingCart },
	{ href: "/cart-items", label: "Cart Items", icon: ListChecks },
	{ href: "/addresses", label: "Addresses", icon: MapPin },
];

export function Sidebar() {
	const pathname = usePathname();
	return (
		<aside className="fixed inset-y-0 left-0 z-40 w-[220px] border-r bg-white/95 backdrop-blur supports-[backdrop-filter]:bg-white/70">
			<div className="px-3 py-3">
				<div className="mb-3 flex items-center gap-2 px-2 text-sm font-semibold text-zinc-800">
					<div className="h-6 w-6 rounded bg-zinc-900" />
					Freshease Admin
				</div>
				<div className="space-y-6">
					<nav className="grid gap-1">
						{links.map((l) => {
							const active = pathname === l.href;
							const Icon = l.icon;
							return (
								<Link
									key={l.href}
									href={l.href}
									className={[
										"flex items-center gap-2 rounded-md px-2 py-2 text-sm",
										active ? "bg-zinc-900 text-white" : "text-zinc-700 hover:bg-zinc-100",
									].join(" ")}
								>
									<Icon className="h-4 w-4" />
									<span className="truncate">{l.label}</span>
								</Link>
							);
						})}
					</nav>
				</div>
			</div>
		</aside>
	);
}
