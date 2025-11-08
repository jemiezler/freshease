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
import type { Review, ReviewPayload } from "@/types/review";
import type { EditDialogProps } from "@/types/dialog";

const reviews = createResource<Review, ReviewPayload, ReviewPayload>({
	basePath: "/reviews",
});

export function EditReviewDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [rating, setRating] = useState<string>("");
	const [comment, setComment] = useState("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await reviews.get(id);
				const r = res.data as Review | undefined;
				if (!cancelled && r) {
					setRating(r.rating != null ? String(r.rating) : "");
					setComment(r.comment ?? "");
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
			const payload: Partial<ReviewPayload> = {
				rating: rating ? Number(rating) : undefined,
				comment: comment || null,
			};
			await reviews.update(id, payload);
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
					<DialogTitle>Edit Review</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading review…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-review-rating">Rating (1-5)</FieldLabel>
							<Input id="edit-review-rating" type="number" min="1" max="5" value={rating} onChange={(e) => setRating(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-review-comment">Comment</FieldLabel>
							<Textarea id="edit-review-comment" value={comment} onChange={(e) => setComment(e.target.value)} />
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

