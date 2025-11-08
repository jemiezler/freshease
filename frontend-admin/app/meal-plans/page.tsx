"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon, PlusIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/meal-plans-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateMealPlanDialog } from "./_components/create-meal-plan-dialog";
import { EditMealPlanDialog } from "./_components/edit-meal-plan-dialog";
import type { MealPlan, MealPlanPayload } from "@/types/meal-plan";

const mealPlans = createResource<MealPlan, MealPlanPayload, MealPlanPayload>({
	basePath: "/meal_plans",
});

export default function MealPlansPage() {
	const [items, setItems] = useState<MealPlan[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await mealPlans.list();
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
			if (!confirm("Delete this meal plan?")) return;
			try {
				await mealPlans.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<MealPlan>[]>(
		() => [
			{
				accessorKey: "week_start",
				header: "Week Start",
				cell: ({ row }) => {
					const date = row.getValue("week_start") as string;
					return date ? new Date(date).toLocaleDateString() : "-";
				},
			},
			{
				accessorKey: "goal",
				header: "Goal",
				cell: ({ row }) => row.getValue("goal") ?? "-",
			},
			{
				accessorKey: "user_id",
				header: "User ID",
				cell: ({ row }) => {
					const userId = row.getValue("user_id") as string;
					return <span className="font-mono text-xs">{userId.slice(0, 8)}...</span>;
				},
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const mealPlan = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => setEditId(mealPlan.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(mealPlan.id)}>
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
				<h1 className="text-3xl font-bold text-zinc-900">Meal Plans</h1>
				<Button onClick={() => setCreateOpen(true)}>
					<PlusIcon className="size-4 mr-2" />
					New Meal Plan
				</Button>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading meal plans…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreateMealPlanDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditMealPlanDialog
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

