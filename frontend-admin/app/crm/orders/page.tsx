"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon, Eye } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "@/app/users/_components/users-table";
import { ColumnDef } from "@tanstack/react-table";
import type { Cart, CartPayload } from "@/types/cart";
import Link from "next/link";

const carts = createResource<Cart, CartPayload, CartPayload>({
	basePath: "/carts",
});

export default function OrdersPage() {
	const [items, setItems] = useState<Cart[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await carts.list();
			setItems(res.data ?? []);
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to load");
		} finally {
			setLoading(false);
		}
	}, []);

	useEffect(() => {
		void load();
	}, [load]);

	const onDelete = useCallback(
		async (id: string) => {
			if (!confirm("Delete this order?")) return;
			try {
				await carts.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<Cart>[]>(
		() => [
			{
				accessorKey: "id",
				header: "Order ID",
				cell: ({ row }) => {
					const id = row.getValue("id") as string;
					return <span className="font-mono text-xs">{id.slice(0, 8)}...</span>;
				},
			},
			{
				accessorKey: "status",
				header: "Status",
				cell: ({ row }) => {
					const status = row.getValue("status") as string;
					return (
						<span
							className={`rounded-full px-2 py-1 text-xs font-medium ${
								status === "completed"
									? "bg-green-100 text-green-800"
									: status === "paid"
										? "bg-blue-100 text-blue-800"
										: status === "pending"
											? "bg-yellow-100 text-yellow-800"
											: status === "cancelled"
												? "bg-red-100 text-red-800"
												: "bg-zinc-100 text-zinc-800"
							}`}
						>
							{status || "unknown"}
						</span>
					);
				},
			},
			{
				accessorKey: "total",
				header: "Total",
				cell: ({ row }) => {
					const total = row.getValue("total") as number;
					return <span className="font-semibold">${total?.toLocaleString() || "0.00"}</span>;
				},
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const cart = row.original;
					return (
						<div className="flex gap-2">
							<Link href={`/crm/orders/${cart.id}`}>
								<Button size="icon" variant="ghost">
									<Eye className="size-4" />
								</Button>
							</Link>
							<Button size="icon" variant="ghost" onClick={() => setEditId(cart.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(cart.id)}>
								<TrashIcon className="size-4 text-red-500" />
							</Button>
						</div>
					);
				},
			},
		],
		[onDelete]
	);

	return (
		<div>
			<div className="mb-6 flex items-center justify-between">
				<div>
					<h1 className="text-3xl font-bold text-zinc-900">Orders</h1>
					<p className="mt-1 text-sm text-zinc-600">Manage customer orders and transactions</p>
				</div>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading orders…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
		</div>
	);
}

