"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon, PlusIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/payments-table";
import { ColumnDef } from "@tanstack/react-table";
import type { Payment, PaymentPayload } from "@/types/payment";

const payments = createResource<Payment, PaymentPayload, PaymentPayload>({
	basePath: "/payments",
});

export default function PaymentsPage() {
	const [items, setItems] = useState<Payment[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await payments.list();
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
			if (!confirm("Delete this payment?")) return;
			try {
				await payments.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<Payment>[]>(
		() => [
			{
				accessorKey: "provider",
				header: "Provider",
				cell: ({ row }) => row.getValue("provider") ?? "-",
			},
			{
				accessorKey: "status",
				header: "Status",
				cell: ({ row }) => {
					const status = row.getValue("status") as string;
					return (
						<span
							className={`rounded-full px-2 py-1 text-xs font-medium ${
								status === "paid"
									? "bg-green-100 text-green-800"
									: status === "pending"
										? "bg-yellow-100 text-yellow-800"
										: status === "failed"
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
				accessorKey: "amount",
				header: "Amount",
				cell: ({ row }) => {
					const amount = row.getValue("amount") as number;
					return <span className="font-semibold">${amount?.toFixed(2) || "0.00"}</span>;
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
					const payment = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => onDelete(payment.id)}>
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
				<h1 className="text-3xl font-bold text-zinc-900">Payments</h1>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading payments…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
		</div>
	);
}

