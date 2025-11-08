"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon, PlusIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/notifications-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateNotificationDialog } from "./_components/create-notification-dialog";
import { EditNotificationDialog } from "./_components/edit-notification-dialog";
import type { Notification, NotificationPayload } from "@/types/notification";

const notifications = createResource<Notification, NotificationPayload, NotificationPayload>({
	basePath: "/notifications",
});

export default function NotificationsPage() {
	const [items, setItems] = useState<Notification[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await notifications.list();
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
			if (!confirm("Delete this notification?")) return;
			try {
				await notifications.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<Notification>[]>(
		() => [
			{
				accessorKey: "title",
				header: "Title",
				cell: ({ row }) => row.getValue("title") ?? "-",
			},
			{
				accessorKey: "channel",
				header: "Channel",
				cell: ({ row }) => row.getValue("channel") ?? "-",
			},
			{
				accessorKey: "status",
				header: "Status",
				cell: ({ row }) => {
					const status = row.getValue("status") as string;
					return (
						<span
							className={`rounded-full px-2 py-1 text-xs font-medium ${
								status === "sent"
									? "bg-green-100 text-green-800"
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
				accessorKey: "created_at",
				header: "Created At",
				cell: ({ row }) => {
					const date = row.getValue("created_at") as string;
					return date ? new Date(date).toLocaleDateString() : "-";
				},
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const notification = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => setEditId(notification.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(notification.id)}>
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
				<h1 className="text-3xl font-bold text-zinc-900">Notifications</h1>
				<Button onClick={() => setCreateOpen(true)}>
					<PlusIcon className="size-4 mr-2" />
					New Notification
				</Button>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading notifications…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreateNotificationDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditNotificationDialog
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

