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
import type { CartItem, CartItemPayload } from "@/types/cart-item";
import type { EditDialogProps } from "@/types/dialog";

const cartItems = createResource<CartItem, CartItemPayload, CartItemPayload>({
	basePath: "/cart_items",
});

export function EditCartItemDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [qty, setQty] = useState<number>(0);
	const [unitPrice, setUnitPrice] = useState<number>(0);
	const [lineTotal, setLineTotal] = useState<number>(0);
	const [cart, setCart] = useState("");
	const [product, setProduct] = useState("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await cartItems.get(id);
				const ci = res.data as CartItem | undefined;
				if (!cancelled && ci) {
					setQty(ci.qty ?? 0);
					setUnitPrice(ci.unit_price ?? 0);
					setLineTotal(ci.line_total ?? 0);
					setCart(ci.cart ?? "");
					setProduct(ci.product ?? "");
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
			const payload: CartItemPayload = {
				qty: qty ? Number(qty) : undefined,
				unit_price: unitPrice ? Number(unitPrice) : undefined,
				line_total: lineTotal ? Number(lineTotal) : undefined,
				cart: cart,
				product: product ?? undefined,
			};
			await cartItems.update(id, payload);
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
					<DialogTitle>Edit Cart Item</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading cart item…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-ci-name">Quantity</FieldLabel>
							<Input id="edit-ci-name" value={qty} onChange={(e) => setQty(Number(e.target.value))} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-ci-description">Unit Price</FieldLabel>
							<Textarea id="edit-ci-description" value={unitPrice} onChange={(e) => setUnitPrice(Number(e.target.value))} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-ci-cart">Cart ID</FieldLabel>
							<Input id="edit-ci-cart" value={cart} onChange={(e) => setCart(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-ci-product">Product</FieldLabel>
							<Input id="edit-ci-product" value={product} onChange={(e) => setProduct(e.target.value)} required />
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
