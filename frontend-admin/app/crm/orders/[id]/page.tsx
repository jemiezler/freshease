"use client";

import { useCallback, useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import type { Cart, CartPayload } from "@/types/cart";
import type { CartItem, CartItemPayload } from "@/types/cart-item";

const carts = createResource<Cart, CartPayload, CartPayload>({
	basePath: "/carts",
});

const cartItems = createResource<CartItem, CartItemPayload, CartItemPayload>({
	basePath: "/cart_items",
});

export default function OrderDetailPage() {
	const params = useParams();
	const id = params.id as string;

	const [order, setOrder] = useState<Cart | null>(null);
	const [items, setItems] = useState<CartItem[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const [orderRes, itemsRes] = await Promise.all([
				carts.get(id),
				cartItems.list().catch((error) => {
					console.error("Failed to load cart items:", error);
					return { data: [] };
				}),
			]);
			setOrder(orderRes.data ?? null);
			// Filter items by cart_id if available
			const filteredItems = (itemsRes.data ?? []).filter(
				(item: CartItem) => item.cart === id
			);
			setItems(filteredItems);
		} catch (error) {
			console.error("Failed to load order details:", error);
			setError(error instanceof Error ? error.message : "Failed to load");
		} finally {
			setLoading(false);
		}
	}, [id]);

	useEffect(() => {
		if (id) {
			void load();
		}
	}, [id, load]);

	if (loading) {
		return (
			<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
				<Spinner className="size-6" />
				<span>Loading order details…</span>
			</div>
		);
	}

	if (error || !order) {
		return (
			<div className="space-y-4">
				<Link href="/crm/orders">
					<Button variant="ghost" className="mb-4">
						<ArrowLeft className="mr-2 size-4" />
						Back to Orders
					</Button>
				</Link>
				<div className="rounded-lg border border-red-200 bg-red-50 p-4">
					<p className="text-red-800">{error || "Order not found"}</p>
				</div>
			</div>
		);
	}

	return (
		<div className="space-y-6">
			<div className="flex items-center justify-between">
				<div className="flex items-center gap-4">
					<Link href="/crm/orders">
						<Button variant="ghost" size="icon">
							<ArrowLeft className="size-4" />
						</Button>
					</Link>
					<div>
						<h1 className="text-3xl font-bold text-zinc-900">Order #{order.id.slice(0, 8)}</h1>
						<p className="mt-1 text-sm text-zinc-600">Order Details</p>
					</div>
				</div>
			</div>

			<div className="grid gap-6 lg:grid-cols-3">
				{/* Order Information */}
				<div className="lg:col-span-2 space-y-6">
					<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
						<h2 className="mb-4 text-lg font-semibold text-zinc-900">Order Items</h2>
						{items.length === 0 ? (
							<p className="text-sm text-zinc-500">No items in this order</p>
						) : (
							<div className="space-y-3">
								{items.map((item) => (
									<div
										key={item.id}
										className="flex items-center justify-between rounded-lg border border-zinc-200 p-4"
									>
										<div>
											<p className="font-medium text-zinc-900">Item #{item.id.slice(0, 8)}</p>
											<p className="text-sm text-zinc-600">
												Quantity: {item.qty?.toLocaleString() || "1"}
												</p>
										</div>
										<div className="text-right">
											<p className="font-semibold text-zinc-900">
												${item.unit_price?.toLocaleString() || "0.00"}
											</p>
										</div>
									</div>
								))}
							</div>
						)}
					</div>
				</div>

				{/* Sidebar */}
				<div className="space-y-6">
					<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
						<h2 className="mb-4 text-lg font-semibold text-zinc-900">Order Summary</h2>
						<div className="space-y-3">
							<div>
								<p className="text-sm text-zinc-600">Status</p>
								<p className="mt-1">
									<span
										className={`rounded-full px-3 py-1 text-sm font-medium ${
											order.status === "completed"
												? "bg-green-100 text-green-800"
												: order.status === "paid"
													? "bg-blue-100 text-blue-800"
													: order.status === "pending"
														? "bg-yellow-100 text-yellow-800"
														: order.status === "cancelled"
															? "bg-red-100 text-red-800"
															: "bg-zinc-100 text-zinc-800"
										}`}
									>
										{order.status || "unknown"}
									</span>
								</p>
							</div>
							<div>
								<p className="text-sm text-zinc-600">Total Amount</p>
								<p className="mt-1 text-2xl font-bold text-zinc-900">
									${order.total?.toLocaleString() || "0.00"}
								</p>
							</div>
							<div>
								<p className="text-sm text-zinc-600">Items Count</p>
								<p className="mt-1 text-lg font-semibold text-zinc-900">{items.length}</p>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}

