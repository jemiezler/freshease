"use client";

import { useState } from "react";
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
import type { DialogProps } from "@/types/dialog";
import { generateUUID } from "@/lib/utils";

const mealPlans = createResource<MealPlan, MealPlanPayload, MealPlanPayload>({
	basePath: "/meal_plans",
});

export function CreateMealPlanDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [weekStart, setWeekStart] = useState("");
	const [goal, setGoal] = useState("");
	const [userId, setUserId] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: MealPlanPayload = {
				id: generateUUID(),
				week_start: weekStart,
				goal: goal || null,
				user_id: userId,
			};
			await mealPlans.create(payload);
			await onSaved();
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to create");
		} finally {
			setSubmitting(false);
		}
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent style={{ maxWidth: "600px" }}>
				<DialogHeader>
					<DialogTitle>New Meal Plan</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="meal-plan-week-start">Week Start *</FieldLabel>
						<Input id="meal-plan-week-start" type="date" value={weekStart} onChange={(e) => setWeekStart(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="meal-plan-goal">Goal</FieldLabel>
						<Textarea id="meal-plan-goal" value={goal} onChange={(e) => setGoal(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="meal-plan-user-id">User ID *</FieldLabel>
						<Input id="meal-plan-user-id" value={userId} onChange={(e) => setUserId(e.target.value)} required />
					</Field>
					{error && <p style={{ color: "red" }}>{error}</p>}
					<DialogFooter>
						<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
							Cancel
						</Button>
						<Button type="submit" disabled={submitting} className="flex items-center gap-2">
							{submitting && <Spinner className="size-4" />}
							{submitting ? "Creating…" : "Create"}
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
}

