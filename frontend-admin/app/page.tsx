"use client";

import { useCallback, useEffect, useState } from "react";
import { createResource } from "@/lib/resource";
import { Spinner } from "@/components/ui/spinner";
import {
	Users,
	Box,
	Store,
	ShoppingCart,
	Boxes,
	Shield,
	Key,
	MapPin,
	ListChecks,
	TrendingUp,
	ArrowRight,
} from "lucide-react";
import Link from "next/link";
import type { User, UserPayload } from "@/types/user";
import type { Product, ProductPayload } from "@/types/product";
import type { Vendor, VendorPayload } from "@/types/vendor";
import type { Cart, CartPayload } from "@/types/cart";
import type { Inventory, InventoryPayload } from "@/types/inventory";
import type { Role, RolePayload } from "@/types/role";
import type { Permission, PermissionPayload } from "@/types/permission";
import type { Address, AddressPayload } from "@/types/address";
import type { CartItem, CartItemPayload } from "@/types/cart-item";

const users = createResource<User, UserPayload, UserPayload>({
	basePath: "/users",
});

const products = createResource<Product, ProductPayload, ProductPayload>({
	basePath: "/products",
});

const vendors = createResource<Vendor, VendorPayload, VendorPayload>({
	basePath: "/vendors",
});

const carts = createResource<Cart, CartPayload, CartPayload>({
	basePath: "/carts",
});

const inventories = createResource<Inventory, InventoryPayload, InventoryPayload>({
	basePath: "/inventories",
});

const roles = createResource<Role, RolePayload, RolePayload>({
	basePath: "/roles",
});

const permissions = createResource<Permission, PermissionPayload, PermissionPayload>({
	basePath: "/permissions",
});

const addresses = createResource<Address, AddressPayload, AddressPayload>({
	basePath: "/addresses",
});

const cartItems = createResource<CartItem, CartItemPayload, CartItemPayload>({
	basePath: "/cart_items",
});

type StatCardProps = {
	title: string;
	value: number | string;
	icon: React.ComponentType<{ className?: string }>;
	href: string;
	loading?: boolean;
	trend?: string;
};

function StatCard({ title, value, icon: Icon, href, loading, trend }: StatCardProps) {
	return (
		<Link
			href={href}
			className="group relative overflow-hidden rounded-lg border border-zinc-200 bg-white p-6 shadow-sm transition-all hover:border-zinc-300 hover:shadow-md"
		>
			<div className="flex items-start justify-between">
				<div className="flex-1">
					<p className="text-sm font-medium text-zinc-600">{title}</p>
					{loading ? (
						<div className="mt-2 flex items-center gap-2">
							<Spinner className="size-4" />
							<span className="text-xs text-zinc-400">Loading...</span>
						</div>
					) : (
						<p className="mt-2 text-3xl font-bold text-zinc-900">{value}</p>
					)}
					{trend && (
						<div className="mt-2 flex items-center gap-1 text-xs text-zinc-500">
							<TrendingUp className="size-3" />
							<span>{trend}</span>
						</div>
					)}
				</div>
				<div className="rounded-lg bg-zinc-100 p-3 group-hover:bg-zinc-200 transition-colors">
					<Icon className="size-6 text-zinc-700" />
				</div>
			</div>
			<div className="mt-4 flex items-center text-xs font-medium text-zinc-600 group-hover:text-zinc-900 transition-colors">
				View all
				<ArrowRight className="ml-1 size-3" />
			</div>
		</Link>
	);
}

export default function Home() {
	const [stats, setStats] = useState({
		users: { count: 0, loading: true },
		products: { count: 0, loading: true },
		vendors: { count: 0, loading: true },
		carts: { count: 0, loading: true },
		inventories: { count: 0, loading: true },
		roles: { count: 0, loading: true },
		permissions: { count: 0, loading: true },
		addresses: { count: 0, loading: true },
		cartItems: { count: 0, loading: true },
	});

	const loadStats = useCallback(async () => {
		const loadCount = async <T,>(
			resource: ReturnType<typeof createResource<T, any, any>>,
			key: keyof typeof stats
		) => {
			try {
				const res = await resource.list();
				setStats((prev) => ({
					...prev,
					[key]: { count: res.data?.length ?? 0, loading: false },
				}));
			} catch (e) {
				setStats((prev) => ({
					...prev,
					[key]: { count: 0, loading: false },
				}));
			}
		};

		await Promise.all([
			loadCount(users, "users"),
			loadCount(products, "products"),
			loadCount(vendors, "vendors"),
			loadCount(carts, "carts"),
			loadCount(inventories, "inventories"),
			loadCount(roles, "roles"),
			loadCount(permissions, "permissions"),
			loadCount(addresses, "addresses"),
			loadCount(cartItems, "cartItems"),
		]);
	}, []);

	useEffect(() => {
		void loadStats();
	}, [loadStats]);

	return (
		<div className="space-y-8">
			<div>
				<h1 className="text-3xl font-bold text-zinc-900">Dashboard</h1>
				<p className="mt-2 text-sm text-zinc-600">
					Welcome to Freshease Admin. Here's an overview of your system.
				</p>
			</div>

			{/* Statistics Grid */}
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
				<StatCard
					title="Users"
					value={stats.users.count}
					icon={Users}
					href="/users"
					loading={stats.users.loading}
				/>
				<StatCard
					title="Products"
					value={stats.products.count}
					icon={Box}
					href="/products"
					loading={stats.products.loading}
				/>
				<StatCard
					title="Vendors"
					value={stats.vendors.count}
					icon={Store}
					href="/vendors"
					loading={stats.vendors.loading}
				/>
				<StatCard
					title="Carts"
					value={stats.carts.count}
					icon={ShoppingCart}
					href="/carts"
					loading={stats.carts.loading}
				/>
				<StatCard
					title="Inventories"
					value={stats.inventories.count}
					icon={Boxes}
					href="/inventories"
					loading={stats.inventories.loading}
				/>
				<StatCard
					title="Roles"
					value={stats.roles.count}
					icon={Shield}
					href="/roles"
					loading={stats.roles.loading}
				/>
				<StatCard
					title="Permissions"
					value={stats.permissions.count}
					icon={Key}
					href="/permissions"
					loading={stats.permissions.loading}
				/>
				<StatCard
					title="Addresses"
					value={stats.addresses.count}
					icon={MapPin}
					href="/addresses"
					loading={stats.addresses.loading}
				/>
				<StatCard
					title="Cart Items"
					value={stats.cartItems.count}
					icon={ListChecks}
					href="/cart-items"
					loading={stats.cartItems.loading}
				/>
			</div>

			{/* Quick Actions */}
			<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
				<h2 className="text-lg font-semibold text-zinc-900 mb-4">Quick Actions</h2>
				<div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
					<Link
						href="/products"
						className="flex items-center gap-3 rounded-lg border border-zinc-200 p-4 transition-colors hover:border-zinc-300 hover:bg-zinc-50"
					>
						<div className="rounded-lg bg-blue-100 p-2">
							<Box className="size-5 text-blue-700" />
						</div>
						<div>
							<p className="font-medium text-zinc-900">Add Product</p>
							<p className="text-xs text-zinc-500">Create a new product</p>
						</div>
					</Link>
					<Link
						href="/vendors"
						className="flex items-center gap-3 rounded-lg border border-zinc-200 p-4 transition-colors hover:border-zinc-300 hover:bg-zinc-50"
					>
						<div className="rounded-lg bg-green-100 p-2">
							<Store className="size-5 text-green-700" />
						</div>
						<div>
							<p className="font-medium text-zinc-900">Add Vendor</p>
							<p className="text-xs text-zinc-500">Register a new vendor</p>
						</div>
					</Link>
					<Link
						href="/users"
						className="flex items-center gap-3 rounded-lg border border-zinc-200 p-4 transition-colors hover:border-zinc-300 hover:bg-zinc-50"
					>
						<div className="rounded-lg bg-purple-100 p-2">
							<Users className="size-5 text-purple-700" />
						</div>
						<div>
							<p className="font-medium text-zinc-900">Add User</p>
							<p className="text-xs text-zinc-500">Create a new user account</p>
						</div>
					</Link>
				</div>
			</div>

			{/* System Overview */}
			<div className="grid gap-6 lg:grid-cols-2">
				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<h2 className="text-lg font-semibold text-zinc-900 mb-4">System Overview</h2>
					<div className="space-y-3">
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Total Entities</span>
							<span className="font-semibold text-zinc-900">
								{Object.values(stats).reduce((acc, stat) => acc + (stat.loading ? 0 : stat.count), 0)}
							</span>
						</div>
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Active Products</span>
							<span className="font-semibold text-zinc-900">
								{stats.products.loading ? "..." : stats.products.count}
							</span>
						</div>
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Registered Users</span>
							<span className="font-semibold text-zinc-900">
								{stats.users.loading ? "..." : stats.users.count}
							</span>
						</div>
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Vendors</span>
							<span className="font-semibold text-zinc-900">
								{stats.vendors.loading ? "..." : stats.vendors.count}
							</span>
						</div>
					</div>
				</div>

				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<h2 className="text-lg font-semibold text-zinc-900 mb-4">Recent Activity</h2>
					<div className="space-y-3">
						<div className="flex items-center gap-3 text-sm">
							<div className="h-2 w-2 rounded-full bg-zinc-400" />
							<span className="text-zinc-600">System is running normally</span>
						</div>
						<div className="flex items-center gap-3 text-sm">
							<div className="h-2 w-2 rounded-full bg-green-500" />
							<span className="text-zinc-600">All services operational</span>
						</div>
						<div className="flex items-center gap-3 text-sm">
							<div className="h-2 w-2 rounded-full bg-blue-500" />
							<span className="text-zinc-600">Ready to manage your inventory</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}
