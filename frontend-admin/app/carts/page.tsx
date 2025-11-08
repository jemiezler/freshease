"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/carts-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateCartDialog } from "./_components/create-cart-dialog";
import { EditCartDialog } from "./_components/edit-cart-dialog";
import type { Cart, CartPayload } from "@/types/cart";

const carts = createResource<Cart, CartPayload, CartPayload>({
	basePath: "/carts",
});

export default function CartsPage() {
	const [items, setItems] = useState<Cart[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

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
			if (!confirm("Delete this cart?")) return;
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
				accessorKey: "status",
				header: "Status",
				cell: ({ row }) => row.getValue("status") ?? "-",
			},
			{
				accessorKey: "total",
				header: "Total",
				cell: ({ row }) => {
					const total = row.getValue("total") as number | undefined;
					return total != null ? total.toFixed(2) : "-";
				},
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const cart = row.original;
					return (
						<div className="flex gap-2">
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
			<div
				style={{
					display: "flex",
					justifyContent: "space-between",
					alignItems: "center",
					marginBottom: 12,
				}}
			>
				<h1 style={{ fontSize: 20, fontWeight: 600 }}>Carts</h1>
				<Button onClick={() => setCreateOpen(true)}>New</Button>
			</div>
			{error && <p style={{ color: "red" }}>{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading carts…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreateCartDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditCartDialog
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
