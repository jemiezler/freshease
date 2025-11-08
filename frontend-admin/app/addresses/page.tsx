"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/addresses-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateAddressDialog } from "./_components/create-address-dialog";
import { EditAddressDialog } from "./_components/edit-address-dialog";
import type { Address, AddressPayload } from "@/types/address";

const addresses = createResource<Address, AddressPayload, AddressPayload>({
	basePath: "/addresses",
});

export default function AddressesPage() {
	const [items, setItems] = useState<Address[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await addresses.list();
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
			if (!confirm("Delete this address?")) return;
			try {
				await addresses.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<Address>[]>(
		() => [
			{
				accessorKey: "line1",
				header: "Line 1",
				cell: ({ row }) => row.getValue("line1") ?? "-",
			},
			{
				accessorKey: "city",
				header: "City",
				cell: ({ row }) => row.getValue("city") ?? "-",
			},
			{
				accessorKey: "province",
				header: "Province",
				cell: ({ row }) => row.getValue("province") ?? "-",
			},
			{
				accessorKey: "country",
				header: "Country",
				cell: ({ row }) => row.getValue("country") ?? "-",
			},
			{
				accessorKey: "zip",
				header: "Zip",
				cell: ({ row }) => row.getValue("zip") ?? "-",
			},
			{
				accessorKey: "is_default",
				header: "Default",
				cell: ({ row }) => (row.getValue("is_default") ? "Yes" : "No"),
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const address = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => setEditId(address.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(address.id)}>
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
				<h1 style={{ fontSize: 20, fontWeight: 600 }}>Addresses</h1>
				<Button onClick={() => setCreateOpen(true)}>New</Button>
			</div>
			{error && <p style={{ color: "red" }}>{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading addresses…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreateAddressDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditAddressDialog
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
