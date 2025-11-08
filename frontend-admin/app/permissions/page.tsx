"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/permissions-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreatePermissionDialog } from "./_components/create-permission-dialog";
import { EditPermissionDialog } from "./_components/edit-permission-dialog";
import type { Permission, PermissionPayload } from "@/types/permission";

const permissions = createResource<Permission, PermissionPayload, PermissionPayload>({
	basePath: "/permissions",
});

export default function PermissionsPage() {
	const [items, setItems] = useState<Permission[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await permissions.list();
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
			if (!confirm("Delete this permission?")) return;
			try {
				await permissions.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<Permission>[]>(
		() => [
			{
				accessorKey: "name",
				header: "Name",
				cell: ({ row }) => row.getValue("name") ?? "-",
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const permission = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => setEditId(permission.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(permission.id)}>
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
				<h1 style={{ fontSize: 20, fontWeight: 600 }}>Permissions</h1>
				<Button onClick={() => setCreateOpen(true)}>New</Button>
			</div>
			{error && <p style={{ color: "red" }}>{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading permissions…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreatePermissionDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditPermissionDialog
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
