import { createResource } from "@/lib/resource";
import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import type { Category, CategoryPayload } from "@/types/catagory";
import type { DialogProps } from "@/types/dialog";

const categories = createResource<Category, CategoryPayload, CategoryPayload>({ basePath: "/product_categories" });

export function CreateCategoryDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
	const [slug, setSlug] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: CategoryPayload = {
				name,
				description: description || undefined,
				slug,
			};
			await categories.create(payload);
			await onSaved();
			setName("");
			setDescription("");
			setSlug("");
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to create");
		} finally {
			setSubmitting(false);
		}
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>New Category</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="category-name">Name</FieldLabel>
						<Input id="category-name" value={name} onChange={(e) => setName(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="category-description">Description</FieldLabel>
						<Textarea id="category-description" value={description} onChange={(e) => setDescription(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="category-slug">Slug</FieldLabel>
						<Input id="category-slug" value={slug} onChange={(e) => setSlug(e.target.value)} />
					</Field>
					{error && <p style={{ color: "red" }}>{error}</p>}
					<DialogFooter>
						<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
							Cancel
						</Button>
						<Button type="submit" disabled={submitting} className="flex items-center gap-2">
							{submitting && <Spinner className="size-4" />}
							{submitting ? "Saving…" : "Create"}
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
}