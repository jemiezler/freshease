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
import type { Recipe, RecipePayload } from "@/types/recipe";
import type { EditDialogProps } from "@/types/dialog";

const recipes = createResource<Recipe, RecipePayload, RecipePayload>({
	basePath: "/recipes",
});

export function EditRecipeDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [name, setName] = useState("");
	const [instructions, setInstructions] = useState("");
	const [kcal, setKcal] = useState<string>("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await recipes.get(id);
				const r = res.data as Recipe | undefined;
				if (!cancelled && r) {
					setName(r.name ?? "");
					setInstructions(r.instructions ?? "");
					setKcal(r.kcal != null ? String(r.kcal) : "");
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
			const payload: Partial<RecipePayload> = {
				name,
				instructions: instructions || null,
				kcal: kcal ? Number(kcal) : undefined,
			};
			await recipes.update(id, payload as RecipePayload);
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
					<DialogTitle>Edit Recipe</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading recipe…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-recipe-name">Name</FieldLabel>
							<Input id="edit-recipe-name" value={name} onChange={(e) => setName(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-recipe-instructions">Instructions</FieldLabel>
							<Textarea id="edit-recipe-instructions" value={instructions} onChange={(e) => setInstructions(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-recipe-kcal">Calories (kcal)</FieldLabel>
							<Input id="edit-recipe-kcal" type="number" min="0" value={kcal} onChange={(e) => setKcal(e.target.value)} />
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

