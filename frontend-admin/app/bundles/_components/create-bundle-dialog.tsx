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
import type { Bundle, BundlePayload } from "@/types/bundle";
import type { DialogProps } from "@/types/dialog";
import { generateUUID } from "@/lib/utils";

const bundles = createResource<Bundle, BundlePayload, BundlePayload>({
	basePath: "/bundles",
});

export function CreateBundleDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
	const [price, setPrice] = useState<string>("");
	const [isActive, setIsActive] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: BundlePayload = {
				id: generateUUID(),
				name,
				description: description || null,
				price: Number(price),
				is_active: isActive,
			};
			await bundles.create(payload);
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
					<DialogTitle>New Bundle</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="bundle-name">Name *</FieldLabel>
						<Input id="bundle-name" value={name} onChange={(e) => setName(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="bundle-description">Description</FieldLabel>
						<Textarea id="bundle-description" value={description} onChange={(e) => setDescription(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="bundle-price">Price *</FieldLabel>
						<Input id="bundle-price" type="number" step="0.01" min="0" value={price} onChange={(e) => setPrice(e.target.value)} required />
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
							{submitting ? "Creating…" : "Create"}
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
}

