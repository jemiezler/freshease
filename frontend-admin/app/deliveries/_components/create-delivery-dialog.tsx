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
import type { Delivery, DeliveryPayload } from "@/types/delivery";
import type { DialogProps } from "@/types/dialog";
import { generateUUID } from "@/lib/utils";

const deliveries = createResource<Delivery, DeliveryPayload, DeliveryPayload>({
	basePath: "/deliveries",
});

export function CreateDeliveryDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [provider, setProvider] = useState("");
	const [trackingNo, setTrackingNo] = useState("");
	const [status, setStatus] = useState("");
	const [eta, setEta] = useState("");
	const [orderId, setOrderId] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: DeliveryPayload = {
				id: generateUUID(),
				provider,
				tracking_no: trackingNo || null,
				status,
				eta: eta || null,
				delivered_at: null,
				order_id: orderId,
			};
			await deliveries.create(payload);
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
					<DialogTitle>New Delivery</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="delivery-provider">Provider *</FieldLabel>
						<Input id="delivery-provider" value={provider} onChange={(e) => setProvider(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="delivery-tracking">Tracking Number</FieldLabel>
						<Input id="delivery-tracking" value={trackingNo} onChange={(e) => setTrackingNo(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="delivery-status">Status *</FieldLabel>
						<Input id="delivery-status" value={status} onChange={(e) => setStatus(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="delivery-eta">ETA</FieldLabel>
						<Input id="delivery-eta" type="datetime-local" value={eta} onChange={(e) => setEta(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="delivery-order-id">Order ID *</FieldLabel>
						<Input id="delivery-order-id" value={orderId} onChange={(e) => setOrderId(e.target.value)} required />
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
