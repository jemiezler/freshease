import { Button } from "@/components/ui/button";
import { DialogHeader, DialogFooter } from "@/components/ui/dialog";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { createResource } from "@/lib/resource";
import { useState, useEffect } from "react";
import type { Product, ProductPayload } from "@/types/product";
import type { EditDialogProps } from "@/types/dialog";

const products = createResource<Product, ProductPayload, ProductPayload>({ basePath: "/products" });

export function EditProductDialog({ id, onOpenChange, onSaved }: EditDialogProps) {
	const [name, setName] = useState("");
	const [price, setPrice] = useState<string>("");
	const [description, setDescription] = useState("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await products.get(id);
				const p = res.data as Product | undefined;
				if (!cancelled && p) {
					setName(p.name ?? "");
					setPrice(p.price != null ? String(p.price) : "");
					setDescription(p.description ?? "");
				}
			} catch (e) {
				setError(e instanceof Error ? e.message : "Failed to load");
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => { cancelled = true; };
	}, [id]);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: ProductPayload = { name, price: price ? Number(price) : undefined, description: description || undefined };
			await products.update(id, payload);
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
					<DialogTitle>Edit Product</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading product…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-product-name">Name</FieldLabel>
							<Input id="edit-product-name" value={name} onChange={(e) => setName(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-product-price">Price</FieldLabel>
							<Input id="edit-product-price" type="number" step="0.01" value={price} onChange={(e) => setPrice(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-product-description">Description</FieldLabel>
							<Textarea id="edit-product-description" value={description} onChange={(e) => setDescription(e.target.value)} />
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