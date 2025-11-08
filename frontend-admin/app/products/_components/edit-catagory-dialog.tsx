import { useState, useEffect } from "react";
import { createResource } from "@/lib/resource";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import type { Category, CategoryPayload } from "@/types/catagory";
import type { EditDialogProps } from "@/types/dialog";

const categories = createResource<Category, CategoryPayload, CategoryPayload>({ basePath: "/product_categories" });

export function EditCategoryDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
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
					setDescription(c.description ?? "");
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
			const payload: CategoryPayload = {
				name,
				description: description || undefined,
				slug,
			};
			await categories.update(id, payload);
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
					<DialogTitle>Edit Category</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading category…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-category-name">Name</FieldLabel>
							<Input id="edit-category-name" value={name} onChange={(e) => setName(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-category-description">Description</FieldLabel>
							<Textarea
								id="edit-category-description"
								value={description}
								onChange={(e) => setDescription(e.target.value)}
							/>
						</Field>
						{error && <p style={{ color: "red" }}>{error}</p>}
						<DialogFooter>
							<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
								Cancel
							</Button>
							<Button type="submit" disabled={submitting} className="flex items-center gap-2">
								{submitting && <Spinner className="size-4" />}
								{submitting ? "Saving…" : "Save"}
							</Button>
						</DialogFooter>
					</form>
				)}
			</DialogContent>
		</Dialog>
	);
}
