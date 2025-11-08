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
import type { CartItem, CartItemPayload } from "@/types/cart-item";
import type { DialogProps } from "@/types/dialog";

const cartItems = createResource<CartItem, CartItemPayload, CartItemPayload>({
	basePath: "/cart_items",
});

export function CreateCartItemDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
	const [cart, setCart] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: CartItemPayload = {
				name: name || undefined,
				description: description || undefined,
				cart: cart || undefined,
			};
			await cartItems.create(payload);
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
					<DialogTitle>New Cart Item</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="ci-name">Name *</FieldLabel>
						<Input id="ci-name" value={name} onChange={(e) => setName(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="ci-description">Description *</FieldLabel>
						<Textarea id="ci-description" value={description} onChange={(e) => setDescription(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="ci-cart">Cart ID *</FieldLabel>
						<Input id="ci-cart" value={cart} onChange={(e) => setCart(e.target.value)} required />
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
