"use client";

import { useState, useEffect } from "react";
import { createResource } from "@/lib/resource";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Field, FieldLabel } from "@/components/ui/field";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import type { MealPlan, MealPlanPayload } from "@/types/meal-plan";
import type { EditDialogProps } from "@/types/dialog";

const mealPlans = createResource<MealPlan, MealPlanPayload, MealPlanPayload>({
	basePath: "/meal_plans",
});

export function EditMealPlanDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [weekStart, setWeekStart] = useState("");
	const [goal, setGoal] = useState("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await mealPlans.get(id);
				const mp = res.data as MealPlan | undefined;
				if (!cancelled && mp) {
					setWeekStart(mp.week_start ? new Date(mp.week_start).toISOString().split("T")[0] : "");
					setGoal(mp.goal ?? "");
				}
			} catch (e) {
				setError(e instanceof Error ? e.message : "Failed to load");
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [id]);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: Partial<MealPlanPayload> = {
				week_start: weekStart,
				goal: goal || null,
			};
			await mealPlans.update(id, payload as MealPlanPayload);
			await onSaved();
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to update");
		} finally {
			setSubmitting(false);
		}
	}

	return (
		<Dialog open onOpenChange={onOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Edit Meal Plan</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading meal plan…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-meal-plan-week-start">Week Start</FieldLabel>
							<Input id="edit-meal-plan-week-start" type="date" value={weekStart} onChange={(e) => setWeekStart(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-meal-plan-goal">Goal</FieldLabel>
							<Textarea id="edit-meal-plan-goal" value={goal} onChange={(e) => setGoal(e.target.value)} />
						</Field>
						{error && <p style={{ color: "red" }}>{error}</p>}
						<DialogFooter>
							<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
								Cancel
							</Button>
							<Button type="submit" disabled={submitting} className="flex items-center gap-2">
								{submitting && <Spinner className="size-4" />}
								{submitting ? "Updating…" : "Update"}
							</Button>
						</DialogFooter>
					</form>
				)}
			</DialogContent>
		</Dialog>
	);
}

