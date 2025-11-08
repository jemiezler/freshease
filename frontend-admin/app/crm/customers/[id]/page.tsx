"use client";

import { useCallback, useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { ArrowLeft, Mail, Phone, User, Calendar, MapPin } from "lucide-react";
import Link from "next/link";
import type { User } from "@/types/user";
import type { Cart } from "@/types/cart";

const users = createResource<User, any, any>({
	basePath: "/users",
});

const carts = createResource<Cart, any, any>({
	basePath: "/carts",
});

export default function CustomerDetailPage() {
	const params = useParams();
	const router = useRouter();
	const id = params.id as string;

	const [customer, setCustomer] = useState<User | null>(null);
	const [orders, setOrders] = useState<Cart[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const [customerRes, ordersRes] = await Promise.all([
				users.get(id),
				carts.list().catch(() => ({ data: [] })),
			]);
			setCustomer(customerRes.data ?? null);
			// Filter orders by user if we have user_id in cart
			setOrders(ordersRes.data ?? []);
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to load");
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
				<span>Loading customer details…</span>
			</div>
		);
	}

	if (error || !customer) {
		return (
			<div className="space-y-4">
				<Link href="/crm/customers">
					<Button variant="ghost" className="mb-4">
						<ArrowLeft className="mr-2 size-4" />
						Back to Customers
					</Button>
				</Link>
				<div className="rounded-lg border border-red-200 bg-red-50 p-4">
					<p className="text-red-800">{error || "Customer not found"}</p>
				</div>
			</div>
		);
	}

	return (
		<div className="space-y-6">
			<div className="flex items-center justify-between">
				<div className="flex items-center gap-4">
					<Link href="/crm/customers">
						<Button variant="ghost" size="icon">
							<ArrowLeft className="size-4" />
						</Button>
					</Link>
					<div>
						<h1 className="text-3xl font-bold text-zinc-900">{customer.name || "Customer"}</h1>
						<p className="mt-1 text-sm text-zinc-600">Customer Details</p>
					</div>
				</div>
			</div>

			<div className="grid gap-6 lg:grid-cols-3">
				{/* Customer Information */}
				<div className="lg:col-span-2 space-y-6">
					<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
						<h2 className="mb-4 text-lg font-semibold text-zinc-900">Contact Information</h2>
						<div className="space-y-4">
							<div className="flex items-center gap-3">
								<Mail className="size-5 text-zinc-400" />
								<div>
									<p className="text-sm text-zinc-600">Email</p>
									<p className="font-medium text-zinc-900">{customer.email || "-"}</p>
								</div>
							</div>
							<div className="flex items-center gap-3">
								<Phone className="size-5 text-zinc-400" />
								<div>
									<p className="text-sm text-zinc-600">Phone</p>
									<p className="font-medium text-zinc-900">{customer.phone || "-"}</p>
								</div>
							</div>
							{customer.bio && (
								<div className="flex items-start gap-3">
									<User className="size-5 text-zinc-400" />
									<div>
										<p className="text-sm text-zinc-600">Bio</p>
										<p className="font-medium text-zinc-900">{customer.bio}</p>
									</div>
								</div>
							)}
						</div>
					</div>

					{/* Order History */}
					<div className="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
						<h2 className="mb-4 text-lg font-semibold text-zinc-900">Order History</h2>
						{orders.length === 0 ? (
							<p className="text-sm text-zinc-500">No orders found</p>
						) : (
							<div className="space-y-3">
								{orders.map((order) => (
									<div
										key={order.id}
										className="flex items-center justify-between rounded-lg border border-zinc-200 p-4"
									>
										<div>
											<p className="font-medium text-zinc-900">Order #{order.id.slice(0, 8)}</p>
											<p className="text-sm text-zinc-600">
												{order.status && (
													<span
														className={`rounded-full px-2 py-1 text-xs font-medium ${
															order.status === "completed"
																? "bg-green-100 text-green-800"
																: order.status === "pending"
																	? "bg-yellow-100 text-yellow-800"
																	: "bg-zinc-100 text-zinc-800"
														}`}
													>
														{order.status}
													</span>
												)}
											</p>
										</div>
										<div className="text-right">
											<p className="font-semibold text-zinc-900">
												${order.total?.toLocaleString() || "0.00"}
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
						<h2 className="mb-4 text-lg font-semibold text-zinc-900">Account Status</h2>
						<div className="space-y-3">
							<div>
								<p className="text-sm text-zinc-600">Status</p>
								<p className="mt-1">
									<span
										className={`rounded-full px-3 py-1 text-sm font-medium ${
											customer.status === "active"
												? "bg-green-100 text-green-800"
												: customer.status === "inactive"
													? "bg-red-100 text-red-800"
													: "bg-zinc-100 text-zinc-800"
										}`}
									>
										{customer.status || "unknown"}
									</span>
								</p>
							</div>
							<div>
								<p className="text-sm text-zinc-600">Total Orders</p>
								<p className="mt-1 text-lg font-semibold text-zinc-900">{orders.length}</p>
							</div>
							<div>
								<p className="text-sm text-zinc-600">Total Spent</p>
								<p className="mt-1 text-lg font-semibold text-zinc-900">
									${orders.reduce((sum, o) => sum + (o.total ?? 0), 0).toLocaleString()}
								</p>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
}

