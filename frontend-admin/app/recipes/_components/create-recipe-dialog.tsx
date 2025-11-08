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
import type { Recipe, RecipePayload } from "@/types/recipe";
import type { DialogProps } from "@/types/dialog";
import { generateUUID } from "@/lib/utils";

const recipes = createResource<Recipe, RecipePayload, RecipePayload>({
	basePath: "/recipes",
});

export function CreateRecipeDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [name, setName] = useState("");
	const [instructions, setInstructions] = useState("");
	const [kcal, setKcal] = useState<string>("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: RecipePayload = {
				id: generateUUID(),
				name,
				instructions: instructions || null,
				kcal: kcal ? Number(kcal) : 0,
			};
			await recipes.create(payload);
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
					<DialogTitle>New Recipe</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="recipe-name">Name *</FieldLabel>
						<Input id="recipe-name" value={name} onChange={(e) => setName(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="recipe-instructions">Instructions</FieldLabel>
						<Textarea id="recipe-instructions" value={instructions} onChange={(e) => setInstructions(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="recipe-kcal">Calories (kcal)</FieldLabel>
						<Input id="recipe-kcal" type="number" min="0" value={kcal} onChange={(e) => setKcal(e.target.value)} />
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

