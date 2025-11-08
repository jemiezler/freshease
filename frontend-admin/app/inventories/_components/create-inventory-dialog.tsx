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
import type { Inventory, InventoryPayload } from "@/types/inventory";
import type { DialogProps } from "@/types/dialog";

const inventories = createResource<Inventory, InventoryPayload, InventoryPayload>({
	basePath: "/inventories",
});

export function CreateInventoryDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [quantity, setQuantity] = useState<string>("");
	const [restockAmount, setRestockAmount] = useState<string>("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: InventoryPayload = {
				quantity: quantity ? Number(quantity) : undefined,
				restock_amount: restockAmount ? Number(restockAmount) : undefined,
			};
			await inventories.create(payload);
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
					<DialogTitle>New Inventory</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="inv-quantity">Quantity *</FieldLabel>
						<Input id="inv-quantity" type="number" min="1" value={quantity} onChange={(e) => setQuantity(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="inv-restock-amount">Restock Amount *</FieldLabel>
						<Input id="inv-restock-amount" type="number" min="1" value={restockAmount} onChange={(e) => setRestockAmount(e.target.value)} required />
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
