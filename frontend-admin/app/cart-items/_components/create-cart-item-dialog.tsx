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
		const [qty, setQty] = useState<string>("");
	const [unitPrice, setUnitPrice] = useState<string>("");
	const [lineTotal, ] = useState<string>("");
	const [cart, setCart] = useState("");
	const [product, ] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: CartItemPayload = {
				qty: qty ? Number(qty) : undefined,
				unit_price: unitPrice ? Number(unitPrice) : undefined,
				line_total: lineTotal ? Number(lineTotal) : undefined,
				cart: cart || "",
				product: product ?? undefined,
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
						<FieldLabel htmlFor="ci-name">Quantity *</FieldLabel>
						<Input id="ci-name" value={qty} onChange={(e) => setQty(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="ci-description">Unit Price *</FieldLabel>
						<Textarea id="ci-description" value={unitPrice} onChange={(e) => setUnitPrice(e.target.value)} required />
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
