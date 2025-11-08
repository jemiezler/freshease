"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon, PlusIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/order-items-table";
import { ColumnDef } from "@tanstack/react-table";
import type { OrderItem, OrderItemPayload } from "@/types/order-item";

const orderItems = createResource<OrderItem, OrderItemPayload, OrderItemPayload>({
	basePath: "/order_items",
});

export default function OrderItemsPage() {
	const [items, setItems] = useState<OrderItem[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await orderItems.list();
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
			if (!confirm("Delete this order item?")) return;
			try {
				await orderItems.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<OrderItem>[]>(
		() => [
			{
				accessorKey: "qty",
				header: "Quantity",
				cell: ({ row }) => row.getValue("qty") ?? "-",
			},
			{
				accessorKey: "unit_price",
				header: "Unit Price",
				cell: ({ row }) => {
					const price = row.getValue("unit_price") as number;
					return `$${price?.toFixed(2) || "0.00"}`;
				},
			},
			{
				accessorKey: "line_total",
				header: "Line Total",
				cell: ({ row }) => {
					const total = row.getValue("line_total") as number;
					return `$${total?.toFixed(2) || "0.00"}`;
				},
			},
			{
				accessorKey: "order_id",
				header: "Order ID",
				cell: ({ row }) => {
					const orderId = row.getValue("order_id") as string;
					return <span className="font-mono text-xs">{orderId.slice(0, 8)}...</span>;
				},
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const item = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => onDelete(item.id)}>
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
				<h1 className="text-3xl font-bold text-zinc-900">Order Items</h1>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading order items…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
		</div>
	);
}

