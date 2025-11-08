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
import type { Delivery, DeliveryPayload } from "@/types/delivery";
import type { EditDialogProps } from "@/types/dialog";

const deliveries = createResource<Delivery, DeliveryPayload, DeliveryPayload>({
	basePath: "/deliveries",
});

export function EditDeliveryDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [provider, setProvider] = useState("");
	const [trackingNo, setTrackingNo] = useState("");
	const [status, setStatus] = useState("");
	const [eta, setEta] = useState("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await deliveries.get(id);
				const d = res.data as Delivery | undefined;
				if (!cancelled && d) {
					setProvider(d.provider ?? "");
					setTrackingNo(d.tracking_no ?? "");
					setStatus(d.status ?? "");
					setEta(d.eta ? new Date(d.eta).toISOString().slice(0, 16) : "");
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
			const payload: Partial<DeliveryPayload> = {
				provider,
				tracking_no: trackingNo || null,
				status,
				eta: eta || null,
			};
			await deliveries.update(id, payload as DeliveryPayload);
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
					<DialogTitle>Edit Delivery</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading delivery…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-delivery-provider">Provider</FieldLabel>
							<Input id="edit-delivery-provider" value={provider} onChange={(e) => setProvider(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-delivery-tracking">Tracking Number</FieldLabel>
							<Input id="edit-delivery-tracking" value={trackingNo} onChange={(e) => setTrackingNo(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-delivery-status">Status</FieldLabel>
							<Input id="edit-delivery-status" value={status} onChange={(e) => setStatus(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-delivery-eta">ETA</FieldLabel>
							<Input id="edit-delivery-eta" type="datetime-local" value={eta} onChange={(e) => setEta(e.target.value)} />
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
