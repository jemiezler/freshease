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
import type { Cart, CartPayload } from "@/types/cart";
import type { EditDialogProps } from "@/types/dialog";

const carts = createResource<Cart, CartPayload, CartPayload>({
	basePath: "/carts",
});

export function EditCartDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [status, setStatus] = useState("");
	const [total, setTotal] = useState<string>("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await carts.get(id);
				const c = res.data as Cart | undefined;
				if (!cancelled && c) {
					setStatus(c.status ?? "");
					setTotal(c.total != null ? String(c.total) : "");
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
			const payload: CartPayload = {
				status: status || undefined,
				total: total ? Number(total) : undefined,
			};
			await carts.update(id, payload);
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
					<DialogTitle>Edit Cart</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading cart…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-cart-status">Status</FieldLabel>
							<Input id="edit-cart-status" value={status} onChange={(e) => setStatus(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-cart-total">Total</FieldLabel>
							<Input id="edit-cart-total" type="number" step="0.01" value={total} onChange={(e) => setTotal(e.target.value)} />
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
