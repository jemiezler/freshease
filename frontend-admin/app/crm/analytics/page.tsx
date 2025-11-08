"use client";

import { useCallback, useEffect, useState } from "react";
import { createResource } from "@/lib/resource";
import { Spinner } from "@/components/ui/spinner";
import { TrendingUp, DollarSign, ShoppingCart, Users } from "lucide-react";
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

export default function AnalyticsPage() {
	const [loading, setLoading] = useState(true);
	const [analytics, setAnalytics] = useState({
		totalCustomers: 0,
		totalOrders: 0,
		totalRevenue: 0,
		averageOrderValue: 0,
		completedOrders: 0,
		pendingOrders: 0,
		totalProducts: 0,
	});

	const loadAnalytics = useCallback(async () => {
		setLoading(true);
		try {
			const [customersRes, ordersRes, productsRes] = await Promise.all([
				users.list().catch((error) => {
					console.error("Failed to load customers:", error);
					return { data: [] };
				}),
				carts.list().catch((error) => {
					console.error("Failed to load carts:", error);
					return { data: [] };
				}),
				products.list().catch((error) => {
					console.error("Failed to load products:", error);
					return { data: [] };
				}),
			]);

			const customers = customersRes.data ?? [];
			const orders = ordersRes.data ?? [];
			const productsList = productsRes.data ?? [];

			const completedOrders = orders.filter(
				(cart: Cart) => cart.status === "completed" || cart.status === "paid"
			);
			const pendingOrders = orders.filter(
				(cart: Cart) => cart.status === "pending" || cart.status === "active"
			);

			const revenue = completedOrders.reduce((sum: number, cart: Cart) => sum + (cart.total ?? 0), 0);
			const avgOrderValue =
				completedOrders.length > 0 ? revenue / completedOrders.length : 0;

			setAnalytics({
				totalCustomers: customers.length,
				totalOrders: orders.length,
				totalRevenue: revenue,
				averageOrderValue: avgOrderValue,
				completedOrders: completedOrders.length,
				pendingOrders: pendingOrders.length,
				totalProducts: productsList.length,
			});
		} catch (e) {
			console.error("Failed to load analytics:", e);
		} finally {
			setLoading(false);
		}
	}, []);

	useEffect(() => {
		void loadAnalytics();
	}, [loadAnalytics]);

	if (loading) {
		return (
			<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
				<Spinner className="size-6" />
				<span>Loading analytics…</span>
			</div>
		);
	}

	return (
		<div className="space-y-8">
			<div>
				<h1 className="text-3xl font-bold text-zinc-900">Analytics & Reports</h1>
				<p className="mt-2 text-sm text-zinc-600">Business insights and performance metrics</p>
			</div>

			{/* Key Metrics Grid */}
			<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<div className="flex items-center justify-between">
						<div>
							<p className="text-sm font-medium text-zinc-600">Total Revenue</p>
							<p className="mt-2 text-3xl font-bold text-zinc-900">
								${analytics.totalRevenue.toLocaleString()}
							</p>
						</div>
						<div className="rounded-lg bg-green-100 p-3">
							<DollarSign className="size-6 text-green-700" />
						</div>
					</div>
				</div>

				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<div className="flex items-center justify-between">
						<div>
							<p className="text-sm font-medium text-zinc-600">Total Orders</p>
							<p className="mt-2 text-3xl font-bold text-zinc-900">{analytics.totalOrders}</p>
							<p className="mt-1 text-xs text-zinc-500">
								{analytics.completedOrders} completed, {analytics.pendingOrders} pending
							</p>
						</div>
						<div className="rounded-lg bg-blue-100 p-3">
							<ShoppingCart className="size-6 text-blue-700" />
						</div>
					</div>
				</div>

				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<div className="flex items-center justify-between">
						<div>
							<p className="text-sm font-medium text-zinc-600">Avg Order Value</p>
							<p className="mt-2 text-3xl font-bold text-zinc-900">
								${Math.round(analytics.averageOrderValue).toLocaleString()}
							</p>
						</div>
						<div className="rounded-lg bg-purple-100 p-3">
							<TrendingUp className="size-6 text-purple-700" />
						</div>
					</div>
				</div>

				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<div className="flex items-center justify-between">
						<div>
							<p className="text-sm font-medium text-zinc-600">Total Customers</p>
							<p className="mt-2 text-3xl font-bold text-zinc-900">{analytics.totalCustomers}</p>
						</div>
						<div className="rounded-lg bg-orange-100 p-3">
							<Users className="size-6 text-orange-700" />
						</div>
					</div>
				</div>
			</div>

			{/* Performance Metrics */}
			<div className="grid gap-6 lg:grid-cols-2">
				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<h2 className="mb-4 text-lg font-semibold text-zinc-900">Order Performance</h2>
					<div className="space-y-4">
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Completion Rate</span>
							<span className="font-semibold text-zinc-900">
								{analytics.totalOrders > 0
									? `${Math.round((analytics.completedOrders / analytics.totalOrders) * 100)}%`
									: "0%"}
							</span>
						</div>
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Pending Orders</span>
							<span className="font-semibold text-zinc-900">{analytics.pendingOrders}</span>
						</div>
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Completed Orders</span>
							<span className="font-semibold text-zinc-900">{analytics.completedOrders}</span>
						</div>
					</div>
				</div>

				<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
					<h2 className="mb-4 text-lg font-semibold text-zinc-900">Business Metrics</h2>
					<div className="space-y-4">
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Total Products</span>
							<span className="font-semibold text-zinc-900">{analytics.totalProducts}</span>
						</div>
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Revenue per Customer</span>
							<span className="font-semibold text-zinc-900">
								$
								{analytics.totalCustomers > 0
									? Math.round(analytics.totalRevenue / analytics.totalCustomers).toLocaleString()
									: "0"}
							</span>
						</div>
						<div className="flex items-center justify-between">
							<span className="text-sm text-zinc-600">Orders per Customer</span>
							<span className="font-semibold text-zinc-900">
								{analytics.totalCustomers > 0
									? (analytics.totalOrders / analytics.totalCustomers).toFixed(1)
									: "0"}
							</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}

