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
import type { Order, OrderPayload } from "@/types/order";
import type { EditDialogProps } from "@/types/dialog";

const orders = createResource<Order, OrderPayload, OrderPayload>({
	basePath: "/orders",
});

export function EditOrderDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [orderNo, setOrderNo] = useState("");
	const [status, setStatus] = useState("");
	const [subtotal, setSubtotal] = useState<string>("");
	const [shippingFee, setShippingFee] = useState<string>("");
	const [discount, setDiscount] = useState<string>("");
	const [total, setTotal] = useState<string>("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await orders.get(id);
				const o = res.data as Order | undefined;
				if (!cancelled && o) {
					setOrderNo(o.order_no ?? "");
					setStatus(o.status ?? "");
					setSubtotal(o.subtotal != null ? String(o.subtotal) : "");
					setShippingFee(o.shipping_fee != null ? String(o.shipping_fee) : "");
					setDiscount(o.discount != null ? String(o.discount) : "");
					setTotal(o.total != null ? String(o.total) : "");
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
			const payload: Partial<OrderPayload> = {
				order_no: orderNo,
				status,
				subtotal: subtotal ? Number(subtotal) : undefined,
				shipping_fee: shippingFee ? Number(shippingFee) : undefined,
				discount: discount ? Number(discount) : undefined,
				total: total ? Number(total) : undefined,
			};
			await orders.update(id, payload);
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
					<DialogTitle>Edit Order</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading order…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-order-no">Order Number</FieldLabel>
							<Input id="edit-order-no" value={orderNo} onChange={(e) => setOrderNo(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-order-status">Status</FieldLabel>
							<Input id="edit-order-status" value={status} onChange={(e) => setStatus(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-order-subtotal">Subtotal</FieldLabel>
							<Input id="edit-order-subtotal" type="number" step="0.01" min="0" value={subtotal} onChange={(e) => setSubtotal(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-order-shipping-fee">Shipping Fee</FieldLabel>
							<Input id="edit-order-shipping-fee" type="number" step="0.01" min="0" value={shippingFee} onChange={(e) => setShippingFee(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-order-discount">Discount</FieldLabel>
							<Input id="edit-order-discount" type="number" step="0.01" min="0" value={discount} onChange={(e) => setDiscount(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-order-total">Total</FieldLabel>
							<Input id="edit-order-total" type="number" step="0.01" min="0" value={total} onChange={(e) => setTotal(e.target.value)} />
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

