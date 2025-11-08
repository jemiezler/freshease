"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { TrashIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/bundle-items-table";
import { ColumnDef } from "@tanstack/react-table";
import type { BundleItem, BundleItemPayload } from "@/types/bundle-item";

const bundleItems = createResource<BundleItem, BundleItemPayload, BundleItemPayload>({
	basePath: "/bundle_items",
});

export default function BundleItemsPage() {
	const [items, setItems] = useState<BundleItem[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await bundleItems.list();
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
			if (!confirm("Delete this bundle item?")) return;
			try {
				await bundleItems.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<BundleItem>[]>(
		() => [
			{
				accessorKey: "qty",
				header: "Quantity",
				cell: ({ row }) => row.getValue("qty") ?? "-",
			},
			{
				accessorKey: "bundle_id",
				header: "Bundle ID",
				cell: ({ row }) => {
					const bundleId = row.getValue("bundle_id") as string;
					return <span className="font-mono text-xs">{bundleId.slice(0, 8)}...</span>;
				},
			},
			{
				accessorKey: "product_id",
				header: "Product ID",
				cell: ({ row }) => {
					const productId = row.getValue("product_id") as string;
					return <span className="font-mono text-xs">{productId.slice(0, 8)}...</span>;
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
			<div className="mb-6">
				<h1 className="text-3xl font-bold text-zinc-900">Bundle Items</h1>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading bundle items…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
		</div>
	);
}

