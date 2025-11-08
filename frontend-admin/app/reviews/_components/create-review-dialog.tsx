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
import type { Review, ReviewPayload } from "@/types/review";
import type { DialogProps } from "@/types/dialog";
import { generateUUID } from "@/lib/utils";

const reviews = createResource<Review, ReviewPayload, ReviewPayload>({
	basePath: "/reviews",
});

export function CreateReviewDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [rating, setRating] = useState<string>("");
	const [comment, setComment] = useState("");
	const [userId, setUserId] = useState("");
	const [productId, setProductId] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: ReviewPayload = {
				id: generateUUID(),
				rating: rating ? Number(rating) : 1,
				comment: comment || null,
				user_id: userId,
				product_id: productId,
				created_at: new Date().toISOString(),
			};
			await reviews.create(payload);
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
					<DialogTitle>New Review</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="review-rating">Rating (1-5) *</FieldLabel>
						<Input id="review-rating" type="number" min="1" max="5" value={rating} onChange={(e) => setRating(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="review-comment">Comment</FieldLabel>
						<Textarea id="review-comment" value={comment} onChange={(e) => setComment(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="review-user-id">User ID *</FieldLabel>
						<Input id="review-user-id" value={userId} onChange={(e) => setUserId(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="review-product-id">Product ID *</FieldLabel>
						<Input id="review-product-id" value={productId} onChange={(e) => setProductId(e.target.value)} required />
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

