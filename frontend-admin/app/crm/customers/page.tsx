"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon, Eye } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "@/app/users/_components/users-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateUserDialog } from "@/app/users/_components/create-user-dialog";
import { EditUserDialog } from "@/app/users/_components/edit-user-dialog";
import type { User, UserPayload } from "@/types/user";
import Link from "next/link";

const users = createResource<User, UserPayload, UserPayload>({
	basePath: "/users",
	updateMethod: "PUT",
});

export default function CustomersPage() {
	const [items, setItems] = useState<User[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await users.list();
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
			if (!confirm("Delete this customer?")) return;
			try {
				await users.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<User>[]>(
		() => [
			{
				accessorKey: "email",
				header: "Email",
				cell: ({ row }) => row.getValue("email") ?? "-",
			},
			{
				accessorKey: "name",
				header: "Name",
				cell: ({ row }) => row.getValue("name") ?? "-",
			},
			{
				accessorKey: "phone",
				header: "Phone",
				cell: ({ row }) => row.getValue("phone") ?? "-",
			},
			{
				accessorKey: "status",
				header: "Status",
				cell: ({ row }) => {
					const status = row.getValue("status") as string;
					return (
						<span
							className={`rounded-full px-2 py-1 text-xs font-medium ${
								status === "active"
									? "bg-green-100 text-green-800"
									: status === "inactive"
										? "bg-red-100 text-red-800"
										: "bg-zinc-100 text-zinc-800"
							}`}
						>
							{status ?? "unknown"}
						</span>
					);
				},
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const user = row.original;
					return (
						<div className="flex gap-2">
							<Link href={`/crm/customers/${user.id}`}>
								<Button size="icon" variant="ghost">
									<Eye className="size-4" />
								</Button>
							</Link>
							<Button size="icon" variant="ghost" onClick={() => setEditId(user.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(user.id)}>
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
					<h1 className="text-3xl font-bold text-zinc-900">Customers</h1>
					<p className="mt-1 text-sm text-zinc-600">Manage customer accounts and information</p>
				</div>
				<Button onClick={() => setCreateOpen(true)}>New Customer</Button>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading customers…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreateUserDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditUserDialog
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

