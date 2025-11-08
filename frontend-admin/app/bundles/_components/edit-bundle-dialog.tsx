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
import type { Bundle, BundlePayload } from "@/types/bundle";
import type { EditDialogProps } from "@/types/dialog";

const bundles = createResource<Bundle, BundlePayload, BundlePayload>({
	basePath: "/bundles",
});

export function EditBundleDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
	const [price, setPrice] = useState<string>("");
	const [isActive, setIsActive] = useState(true);
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await bundles.get(id);
				const b = res.data as Bundle | undefined;
				if (!cancelled && b) {
					setName(b.name ?? "");
					setDescription(b.description ?? "");
					setPrice(b.price != null ? String(b.price) : "");
					setIsActive(b.is_active ?? true);
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
			const payload: Partial<BundlePayload> = {
				name,
				description: description || null,
				price: Number(price),
				is_active: isActive,
			};
			await bundles.update(id, payload);
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
					<DialogTitle>Edit Bundle</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading bundle…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-bundle-name">Name</FieldLabel>
							<Input id="edit-bundle-name" value={name} onChange={(e) => setName(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-bundle-description">Description</FieldLabel>
							<Textarea id="edit-bundle-description" value={description} onChange={(e) => setDescription(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-bundle-price">Price</FieldLabel>
							<Input id="edit-bundle-price" type="number" step="0.01" min="0" value={price} onChange={(e) => setPrice(e.target.value)} required />
						</Field>
						<Field>
							<label className="flex items-center gap-2">
								<input type="checkbox" checked={isActive} onChange={(e) => setIsActive(e.target.checked)} />
								<span>Active</span>
							</label>
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

