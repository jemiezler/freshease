"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutDashboard, Box, Tag, Boxes, Store, Users, Shield, Key, ShoppingCart, ListChecks, MapPin, Bell, Star, Gift, Calendar, ChefHat, Truck, BarChart3, Receipt, CreditCard } from "lucide-react";
import type { LucideIcon } from "lucide-react";

type NavLink = {
	href: string;
	label: string;
	icon: LucideIcon;
	children?: Array<{ href: string; label: string }>;
};

const links: NavLink[] = [
	{ href: "/", label: "Dashboard", icon: LayoutDashboard },
	{ 
		href: "/crm", 
		label: "CRM", 
		icon: BarChart3,
		children: [
			{ href: "/crm", label: "CRM Dashboard" },
			{ href: "/crm/customers", label: "Customers" },
			{ href: "/crm/orders", label: "Orders" },
			{ href: "/crm/analytics", label: "Analytics" },
		]
	},
	{ href: "/products", label: "Products", icon: Box },
	{ href: "/categories", label: "Categories", icon: Tag },
	{ href: "/inventories", label: "Inventories", icon: Boxes },
	{ href: "/vendors", label: "Vendors", icon: Store },
	{ href: "/bundles", label: "Bundles", icon: Gift },
	{ href: "/users", label: "Users", icon: Users },
	{ href: "/roles", label: "Roles", icon: Shield },
	{ href: "/permissions", label: "Permissions", icon: Key },
	{ href: "/carts", label: "Carts", icon: ShoppingCart },
	{ href: "/cart-items", label: "Cart Items", icon: ListChecks },
	{ href: "/orders", label: "Orders", icon: Receipt },
	{ href: "/order-items", label: "Order Items", icon: ListChecks },
	{ href: "/payments", label: "Payments", icon: CreditCard },
	{ href: "/addresses", label: "Addresses", icon: MapPin },
	{ href: "/deliveries", label: "Deliveries", icon: Truck },
	{ href: "/notifications", label: "Notifications", icon: Bell },
	{ href: "/reviews", label: "Reviews", icon: Star },
	{ href: "/meal-plans", label: "Meal Plans", icon: Calendar },
	{ href: "/recipes", label: "Recipes", icon: ChefHat },
];

export function Sidebar() {
	const pathname = usePathname();
	return (
		<aside className="fixed inset-y-0 left-0 z-40 w-[220px] border-r bg-white/95 backdrop-blur supports-backdrop-filter:bg-white/70 overflow-y-auto">
			<div className="px-3 py-3">
				<div className="mb-3 flex items-center gap-2 px-2 text-sm font-semibold text-zinc-800">
					<div className="h-6 w-6 rounded bg-zinc-900" />
					Freshease Admin
				</div>
				<div className="space-y-6">
					<nav className="grid gap-1">
						{links.map((l) => {
							const isActive = pathname === l.href || (l.children && pathname?.startsWith(l.href));
							const Icon = l.icon;
							
							if (l.children) {
								const hasActiveChild = l.children.some(child => pathname === child.href || pathname?.startsWith(child.href + "/"));
								return (
									<div key={l.href} className="space-y-1">
										<Link
											href={l.href}
											className={[
												"flex items-center gap-2 rounded-md px-2 py-2 text-sm",
												isActive ? "bg-zinc-900 text-white" : "text-zinc-700 hover:bg-zinc-100",
											].join(" ")}
										>
											<Icon className="h-4 w-4" />
											<span className="truncate">{l.label}</span>
										</Link>
										{hasActiveChild && (
											<div className="ml-4 space-y-1 border-l-2 border-zinc-200 pl-2">
												{l.children.map((child) => {
													const childActive = pathname === child.href || pathname?.startsWith(child.href + "/");
													return (
														<Link
															key={child.href}
															href={child.href}
															className={[
																"flex items-center gap-2 rounded-md px-2 py-1.5 text-xs",
																childActive ? "bg-zinc-100 text-zinc-900 font-medium" : "text-zinc-600 hover:bg-zinc-50",
															].join(" ")}
														>
															<span className="truncate">{child.label}</span>
														</Link>
													);
												})}
											</div>
										)}
									</div>
								);
							}
							
							return (
								<Link
									key={l.href}
									href={l.href}
									className={[
										"flex items-center gap-2 rounded-md px-2 py-2 text-sm",
										isActive ? "bg-zinc-900 text-white" : "text-zinc-700 hover:bg-zinc-100",
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
