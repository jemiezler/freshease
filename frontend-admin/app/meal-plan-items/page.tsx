"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { TrashIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/meal-plan-items-table";
import { ColumnDef } from "@tanstack/react-table";
import type { MealPlanItem, MealPlanItemPayload } from "@/types/meal-plan-item";

const mealPlanItems = createResource<MealPlanItem, MealPlanItemPayload, MealPlanItemPayload>({
	basePath: "/meal_plan_items",
});

export default function MealPlanItemsPage() {
	const [items, setItems] = useState<MealPlanItem[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await mealPlanItems.list();
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
			if (!confirm("Delete this meal plan item?")) return;
			try {
				await mealPlanItems.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<MealPlanItem>[]>(
		() => [
			{
				accessorKey: "day",
				header: "Day",
				cell: ({ row }) => {
					const day = row.getValue("day") as string;
					return day ? new Date(day).toLocaleDateString() : "-";
				},
			},
			{
				accessorKey: "slot",
				header: "Slot",
				cell: ({ row }) => row.getValue("slot") ?? "-",
			},
			{
				accessorKey: "meal_plan_id",
				header: "Meal Plan ID",
				cell: ({ row }) => {
					const mealPlanId = row.getValue("meal_plan_id") as string;
					return <span className="font-mono text-xs">{mealPlanId.slice(0, 8)}...</span>;
				},
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
				<h1 className="text-3xl font-bold text-zinc-900">Meal Plan Items</h1>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading meal plan items…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
		</div>
	);
}

