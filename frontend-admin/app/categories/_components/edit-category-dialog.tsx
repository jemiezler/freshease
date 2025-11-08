"use client";

import { useState, useEffect } from "react";
import { createResource } from "@/lib/resource";
import { Input } from "@/components/ui/input";
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
import type { Category, CategoryPayload } from "@/types/category";
import type { EditDialogProps } from "@/types/dialog";

const categories = createResource<Category, CategoryPayload, CategoryPayload>({
	basePath: "/categories",
});

export function EditCategoryDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [name, setName] = useState("");
	const [slug, setSlug] = useState("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await categories.get(id);
				const c = res.data as Category | undefined;
				if (!cancelled && c) {
					setName(c.name ?? "");
					setSlug(c.slug ?? "");
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
			const payload: Partial<CategoryPayload> = {
				name,
				slug: slug || name.toLowerCase().replace(/\s+/g, "-"),
				updated_at: new Date().toISOString(),
			};
			await categories.update(id, payload as CategoryPayload);
			await onSaved();
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to update");
		} finally {
			setSubmitting(false);
		}
	}

	if (loading) {
		return (
			<Dialog open onOpenChange={onOpenChange}>
				<DialogContent>
					<div className="flex items-center justify-center p-8">
						<Spinner className="size-6" />
					</div>
				</DialogContent>
			</Dialog>
		);
	}

	return (
		<Dialog open onOpenChange={onOpenChange}>
			<DialogContent style={{ maxWidth: "600px" }}>
				<DialogHeader>
					<DialogTitle>Edit Category</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="category-name">Name</FieldLabel>
						<Input id="category-name" value={name} onChange={(e) => setName(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="category-slug">Slug</FieldLabel>
						<Input 
							id="category-slug" 
							value={slug} 
							onChange={(e) => setSlug(e.target.value)} 
							required
						/>
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
			</DialogContent>
		</Dialog>
	);
}

