"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { TrashIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/recipe-items-table";
import { ColumnDef } from "@tanstack/react-table";
import type { RecipeItem, RecipeItemPayload } from "@/types/recipe-item";

const recipeItems = createResource<RecipeItem, RecipeItemPayload, RecipeItemPayload>({
	basePath: "/recipe_items",
});

export default function RecipeItemsPage() {
	const [items, setItems] = useState<RecipeItem[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await recipeItems.list();
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
			if (!confirm("Delete this recipe item?")) return;
			try {
				await recipeItems.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<RecipeItem>[]>(
		() => [
			{
				accessorKey: "amount",
				header: "Amount",
				cell: ({ row }) => row.getValue("amount") ?? "-",
			},
			{
				accessorKey: "unit",
				header: "Unit",
				cell: ({ row }) => row.getValue("unit") ?? "-",
			},
			{
				accessorKey: "recipe_id",
				header: "Recipe ID",
				cell: ({ row }) => {
					const recipeId = row.getValue("recipe_id") as string;
					return <span className="font-mono text-xs">{recipeId.slice(0, 8)}...</span>;
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
				<h1 className="text-3xl font-bold text-zinc-900">Recipe Items</h1>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading recipe items…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
		</div>
	);
}

