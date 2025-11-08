"use client";

import { useCallback, useEffect, useState } from "react";
import { createResource } from "@/lib/resource";
import { Spinner } from "@/components/ui/spinner";
import {
	Users,
	ShoppingCart,
	DollarSign,
	TrendingUp,
	Package,
	ArrowRight,
	Activity,
} from "lucide-react";
import Link from "next/link";
import type { User, UserPayload } from "@/types/user";
import type { Cart, CartPayload } from "@/types/cart";
import type { Product, ProductPayload } from "@/types/product";

const users = createResource<User, UserPayload, UserPayload>({
	basePath: "/users",
});

const carts = createResource<Cart, CartPayload, CartPayload>({
	basePath: "/carts",
});

const products = createResource<Product, ProductPayload, ProductPayload>({
	basePath: "/products",
});

type StatCardProps = {
	title: string;
	value: number | string;
	icon: React.ComponentType<{ className?: string }>;
	href?: string;
	loading?: boolean;
	trend?: string;
	subtitle?: string;
};

function StatCard({ title, value, icon: Icon, href, loading, trend, subtitle }: StatCardProps) {
	const content = (
		<div className="group relative overflow-hidden rounded-lg border border-zinc-200 bg-white p-6 shadow-sm transition-all hover:border-zinc-300 hover:shadow-md">
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
					{subtitle && (
						<p className="mt-1 text-xs text-zinc-500">{subtitle}</p>
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
			{href && (
				<div className="mt-4 flex items-center text-xs font-medium text-zinc-600 group-hover:text-zinc-900 transition-colors">
					View all
					<ArrowRight className="ml-1 size-3" />
				</div>
			)}
		</div>
	);

	if (href) {
		return <Link href={href}>{content}</Link>;
	}

	return content;
}

export default function CRMPage() {
	const [stats, setStats] = useState({
		customers: { count: 0, loading: true },
		orders: { count: 0, loading: true, active: 0 },
		revenue: { total: 0, loading: true },
		products: { count: 0, loading: true },
	});

	const loadStats = useCallback(async () => {
		try {
			const [customersRes, ordersRes, productsRes] = await Promise.all([
				users.list().catch(() => ({ data: [] })),
				carts.list().catch(() => ({ data: [] })),
				products.list().catch(() => ({ data: [] })),
			]);

			const customers = customersRes.data ?? [];
			const orders = ordersRes.data ?? [];
			const productsList = productsRes.data ?? [];

			// Calculate revenue from carts (assuming carts with status "completed" or "paid")
			const revenue = orders
				.filter((cart: Cart) => cart.status === "completed" || cart.status === "paid")
				.reduce((sum: number, cart: Cart) => sum + (cart.total ?? 0), 0);

			const activeOrders = orders.filter(
				(cart: Cart) => cart.status === "pending" || cart.status === "active"
			).length;

			setStats({
				customers: { count: customers.length, loading: false },
				orders: { count: orders.length, active: activeOrders, loading: false },
				revenue: { total: revenue, loading: false },
				products: { count: productsList.length, loading: false },
			});
		} catch (e) {
			console.error("Failed to load stats:", e);
			setStats({
				customers: { count: 0, loading: false },
				orders: { count: 0, active: 0, loading: false },
				revenue: { total: 0, loading: false },
				products: { count: 0, loading: false },
			});
		}
	}, []);

	useEffect(() => {
		void loadStats();
	}, [loadStats]);

	return (
		<div className="space-y-8">
			<div>
				<h1 className="text-3xl font-bold text-zinc-900">CRM Dashboard</h1>
				<p className="mt-2 text-sm text-zinc-600">
					Customer relationship management and business analytics
				</p>
			</div>

			{/* Key Metrics */}
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
				<StatCard
					title="Total Customers"
					value={stats.customers.count}
					icon={Users}
					href="/crm/customers"
					loading={stats.customers.loading}
				/>
				<StatCard
					title="Total Orders"
					value={stats.orders.count}
					icon={ShoppingCart}
					href="/crm/orders"
					loading={stats.orders.loading}
					subtitle={`${stats.orders.active} active`}
				/>
				<StatCard
					title="Total Revenue"
					value={`$${stats.revenue.total.toLocaleString()}`}
					icon={DollarSign}
					loading={stats.revenue.loading}
				/>
				<StatCard
					title="Products"
					value={stats.products.count}
					icon={Package}
					href="/products"
					loading={stats.products.loading}
				/>
			</div>

			{/* Quick Actions */}
			<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
				<h2 className="text-lg font-semibold text-zinc-900 mb-4">Quick Actions</h2>
				<div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
					<Link
						href="/crm/customers"
						className="flex items-center gap-3 rounded-lg border border-zinc-200 p-4 transition-colors hover:border-zinc-300 hover:bg-zinc-50"
					>
						<div className="rounded-lg bg-blue-100 p-2">
							<Users className="size-5 text-blue-700" />
						</div>
						<div>
							<p className="font-medium text-zinc-900">Manage Customers</p>
							<p className="text-xs text-zinc-500">View and manage customer accounts</p>
						</div>
					</Link>
					<Link
						href="/crm/orders"
						className="flex items-center gap-3 rounded-lg border border-zinc-200 p-4 transition-colors hover:border-zinc-300 hover:bg-zinc-50"
					>
						<div className="rounded-lg bg-green-100 p-2">
							<ShoppingCart className="size-5 text-green-700" />
						</div>
						<div>
							<p className="font-medium text-zinc-900">View Orders</p>
							<p className="text-xs text-zinc-500">Track and manage customer orders</p>
						</div>
					</Link>
					<Link
						href="/crm/analytics"
						className="flex items-center gap-3 rounded-lg border border-zinc-200 p-4 transition-colors hover:border-zinc-300 hover:bg-zinc-50"
					>
						<div className="rounded-lg bg-purple-100 p-2">
							<Activity className="size-5 text-purple-700" />
						</div>
						<div>
							<p className="font-medium text-zinc-900">Analytics</p>
							<p className="text-xs text-zinc-500">View business insights and reports</p>
						</div>
					</Link>
				</div>
			</div>

			{/* Recent Activity Summary */}
			<div className="grid gap-6 lg:grid-cols-2">
				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<h2 className="text-lg font-semibold text-zinc-900 mb-4">Business Overview</h2>
					<div className="space-y-3">
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Customer Growth</span>
							<span className="font-semibold text-zinc-900">
								{stats.customers.loading ? "..." : `${stats.customers.count} customers`}
							</span>
						</div>
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Order Completion Rate</span>
							<span className="font-semibold text-zinc-900">
								{stats.orders.loading
									? "..."
									: stats.orders.count > 0
										? `${Math.round(((stats.orders.count - stats.orders.active) / stats.orders.count) * 100)}%`
										: "0%"}
							</span>
						</div>
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Average Order Value</span>
							<span className="font-semibold text-zinc-900">
								{stats.orders.loading || stats.orders.count === 0
									? "..."
									: `$${Math.round(stats.revenue.total / (stats.orders.count - stats.orders.active) || 0)}`}
							</span>
						</div>
					</div>
				</div>

				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<h2 className="text-lg font-semibold text-zinc-900 mb-4">System Status</h2>
					<div className="space-y-3">
						<div className="flex items-center gap-3 text-sm">
							<div className="h-2 w-2 rounded-full bg-green-500" />
							<span className="text-zinc-600">CRM system operational</span>
						</div>
						<div className="flex items-center gap-3 text-sm">
							<div className="h-2 w-2 rounded-full bg-blue-500" />
							<span className="text-zinc-600">All services connected</span>
						</div>
						<div className="flex items-center gap-3 text-sm">
							<div className="h-2 w-2 rounded-full bg-zinc-400" />
							<span className="text-zinc-600">Ready for customer management</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}

