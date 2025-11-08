"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/cart-items-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateCartItemDialog } from "./_components/create-cart-item-dialog";
import { EditCartItemDialog } from "./_components/edit-cart-item-dialog";
import type { CartItem, CartItemPayload } from "@/types/cart-item";

const cartItems = createResource<CartItem, CartItemPayload, CartItemPayload>({
	basePath: "/cart_items",
});

export default function CartItemsPage() {
	const [items, setItems] = useState<CartItem[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await cartItems.list();
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
			if (!confirm("Delete this item?")) return;
			try {
				await cartItems.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<CartItem>[]>(
		() => [
			{
				accessorKey: "name",
				header: "Name",
				cell: ({ row }) => row.getValue("name") ?? "-",
			},
			{
				accessorKey: "description",
				header: "Description",
				cell: ({ row }) => {
					const desc = row.getValue("description") as string | undefined;
					return desc ? (desc.length > 50 ? desc.substring(0, 50) + "..." : desc) : "-";
				},
			},
			{
				accessorKey: "cart",
				header: "Cart ID",
				cell: ({ row }) => row.getValue("cart") ?? "-",
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const cartItem = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => setEditId(cartItem.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(cartItem.id)}>
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
				<h1 style={{ fontSize: 20, fontWeight: 600 }}>Cart Items</h1>
				<Button onClick={() => setCreateOpen(true)}>New</Button>
			</div>
			{error && <p style={{ color: "red" }}>{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading cart items…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreateCartItemDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditCartItemDialog
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
