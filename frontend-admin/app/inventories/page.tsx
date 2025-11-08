"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/inventories-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateInventoryDialog } from "./_components/create-inventory-dialog";
import { EditInventoryDialog } from "./_components/edit-inventory-dialog";
import type { Inventory, InventoryPayload } from "@/types/inventory";

const inventories = createResource<Inventory, InventoryPayload, InventoryPayload>({
	basePath: "/inventories",
});

export default function InventoriesPage() {
	const [items, setItems] = useState<Inventory[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await inventories.list();
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
			if (!confirm("Delete this inventory?")) return;
			try {
				await inventories.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<Inventory>[]>(
		() => [
			{
				accessorKey: "quantity",
				header: "Quantity",
				cell: ({ row }) => row.getValue("quantity") ?? "-",
			},
			{
				accessorKey: "restock_amount",
				header: "Restock Amount",
				cell: ({ row }) => row.getValue("restock_amount") ?? "-",
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const inventory = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => setEditId(inventory.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(inventory.id)}>
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
			<div
				style={{
					display: "flex",
					justifyContent: "space-between",
					alignItems: "center",
					marginBottom: 12,
				}}
			>
				<h1 style={{ fontSize: 20, fontWeight: 600 }}>Inventories</h1>
				<Button onClick={() => setCreateOpen(true)}>New</Button>
			</div>
			{error && <p style={{ color: "red" }}>{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading inventories…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreateInventoryDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditInventoryDialog
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
