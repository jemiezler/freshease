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
import type { Order, OrderPayload } from "@/types/order";
import type { DialogProps } from "@/types/dialog";
import { generateUUID } from "@/lib/utils";

const orders = createResource<Order, OrderPayload, OrderPayload>({
	basePath: "/orders",
});

export function CreateOrderDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [orderNo, setOrderNo] = useState("");
	const [status, setStatus] = useState("");
	const [subtotal, setSubtotal] = useState<string>("");
	const [shippingFee, setShippingFee] = useState<string>("");
	const [discount, setDiscount] = useState<string>("");
	const [total, setTotal] = useState<string>("");
	const [userId, setUserId] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: OrderPayload = {
				id: generateUUID(),
				order_no: orderNo,
				status,
				subtotal: subtotal ? Number(subtotal) : 0,
				shipping_fee: shippingFee ? Number(shippingFee) : 0,
				discount: discount ? Number(discount) : 0,
				total: total ? Number(total) : 0,
				placed_at: new Date().toISOString(),
				user_id: userId,
			};
			await orders.create(payload);
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
					<DialogTitle>New Order</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="order-no">Order Number *</FieldLabel>
						<Input id="order-no" value={orderNo} onChange={(e) => setOrderNo(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="order-status">Status *</FieldLabel>
						<Input id="order-status" value={status} onChange={(e) => setStatus(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="order-subtotal">Subtotal</FieldLabel>
						<Input id="order-subtotal" type="number" step="0.01" min="0" value={subtotal} onChange={(e) => setSubtotal(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="order-shipping-fee">Shipping Fee</FieldLabel>
						<Input id="order-shipping-fee" type="number" step="0.01" min="0" value={shippingFee} onChange={(e) => setShippingFee(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="order-discount">Discount</FieldLabel>
						<Input id="order-discount" type="number" step="0.01" min="0" value={discount} onChange={(e) => setDiscount(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="order-total">Total</FieldLabel>
						<Input id="order-total" type="number" step="0.01" min="0" value={total} onChange={(e) => setTotal(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="order-user-id">User ID *</FieldLabel>
						<Input id="order-user-id" value={userId} onChange={(e) => setUserId(e.target.value)} required />
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

