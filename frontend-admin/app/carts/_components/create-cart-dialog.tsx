"use client";

import { useState } from "react";
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
import type { DialogProps } from "@/types/dialog";

const carts = createResource<Cart, CartPayload, CartPayload>({
	basePath: "/carts",
});

export function CreateCartDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [status, setStatus] = useState("");
	const [total, setTotal] = useState<string>("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: CartPayload = {
				status: status || undefined,
				total: total ? Number(total) : undefined,
			};
			await carts.create(payload);
			await onSaved();
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
					<DialogTitle>New Cart</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="cart-status">Status</FieldLabel>
						<Input id="cart-status" value={status} onChange={(e) => setStatus(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="cart-total">Total</FieldLabel>
						<Input id="cart-total" type="number" step="0.01" value={total} onChange={(e) => setTotal(e.target.value)} />
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
