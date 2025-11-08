"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon, PlusIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/deliveries-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateDeliveryDialog } from "./_components/create-delivery-dialog";
import { EditDeliveryDialog } from "./_components/edit-delivery-dialog";
import type { Delivery, DeliveryPayload } from "@/types/delivery";

const deliveries = createResource<Delivery, DeliveryPayload, DeliveryPayload>({
	basePath: "/deliveries",
});

export default function DeliveriesPage() {
	const [items, setItems] = useState<Delivery[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await deliveries.list();
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
			if (!confirm("Delete this delivery?")) return;
			try {
				await deliveries.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<Delivery>[]>(
		() => [
			{
				accessorKey: "provider",
				header: "Provider",
				cell: ({ row }) => row.getValue("provider") ?? "-",
			},
			{
				accessorKey: "tracking_no",
				header: "Tracking Number",
				cell: ({ row }) => row.getValue("tracking_no") ?? "-",
			},
			{
				accessorKey: "status",
				header: "Status",
				cell: ({ row }) => {
					const status = row.getValue("status") as string;
					return (
						<span
							className={`rounded-full px-2 py-1 text-xs font-medium ${
								status === "delivered"
									? "bg-green-100 text-green-800"
									: status === "in_transit"
										? "bg-blue-100 text-blue-800"
										: status === "pending"
											? "bg-yellow-100 text-yellow-800"
											: "bg-zinc-100 text-zinc-800"
							}`}
						>
							{status || "unknown"}
						</span>
					);
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
					const delivery = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => setEditId(delivery.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(delivery.id)}>
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
				<h1 className="text-3xl font-bold text-zinc-900">Deliveries</h1>
				<Button onClick={() => setCreateOpen(true)}>
					<PlusIcon className="size-4 mr-2" />
					New Delivery
				</Button>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading deliveries…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreateDeliveryDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditDeliveryDialog
					id={editId}
					onOpenChange={(open) => {
						if (!open) setEditId(null);
					}}
					onSaved={async () => {
						setEditId(null);
						await load();
					}}
				/>
			)}
		</div>
	);
}

