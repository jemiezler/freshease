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
import type { Inventory, InventoryPayload } from "@/types/inventory";
import type { EditDialogProps } from "@/types/dialog";

const inventories = createResource<Inventory, InventoryPayload, InventoryPayload>({
	basePath: "/inventories",
});

export function EditInventoryDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [quantity, setQuantity] = useState<string>("");
	const [restockAmount, setRestockAmount] = useState<string>("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await inventories.get(id);
				const inv = res.data as Inventory | undefined;
				if (!cancelled && inv) {
					setQuantity(inv.quantity != null ? String(inv.quantity) : "");
					setRestockAmount(inv.restock_amount != null ? String(inv.restock_amount) : "");
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
			const payload: InventoryPayload = {
				quantity: quantity ? Number(quantity) : undefined,
				restock_amount: restockAmount ? Number(restockAmount) : undefined,
			};
			await inventories.update(id, payload);
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
					<DialogTitle>Edit Inventory</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading inventory…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-inv-quantity">Quantity</FieldLabel>
							<Input id="edit-inv-quantity" type="number" min="1" value={quantity} onChange={(e) => setQuantity(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-inv-restock-amount">Restock Amount</FieldLabel>
							<Input id="edit-inv-restock-amount" type="number" min="1" value={restockAmount} onChange={(e) => setRestockAmount(e.target.value)} />
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
