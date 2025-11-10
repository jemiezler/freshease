import { Button } from "@/components/ui/button";
import { DialogHeader, DialogFooter } from "@/components/ui/dialog";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { createResource } from "@/lib/resource";
import { apiClient } from "@/lib/api";
import { useState, useEffect } from "react";
import type { Product, ProductPayload } from "@/types/product";
import type { EditDialogProps } from "@/types/dialog";

const products = createResource<Product, ProductPayload, ProductPayload>({ basePath: "/products" });

export function EditProductDialog({ id, onOpenChange, onSaved }: EditDialogProps) {
	const [name, setName] = useState("");
	const [price, setPrice] = useState<string>("");
	const [description, setDescription] = useState("");
	const [, setImageUrl] = useState("");
	const [imageFile, setImageFile] = useState<File | null>(null);
	const [imagePreview, setImagePreview] = useState<string>("");
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
					setImageUrl(p.image_url ?? "");
					setImagePreview(p.image_url ?? "");
				}
			} catch (e) {
				setError(e instanceof Error ? e.message : "Failed to load");
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => { cancelled = true; };
	}, [id]);

	function handleImageChange(e: React.ChangeEvent<HTMLInputElement>) {
		const file = e.target.files?.[0];
		if (!file) return;

		setImageFile(file);
		setError(null);

		// Create preview URL
		const reader = new FileReader();
		reader.onloadend = () => {
			setImagePreview(reader.result as string);
		};
		reader.readAsDataURL(file);
	}

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			// Prepare payload for update (all fields are optional except ID)
			// Always include ID, and include other fields if they have values
			const payload: Record<string, unknown> = {
				id,
			};

			// Add fields only if they have values
			if (name.trim()) {
				payload.name = name.trim();
			}
			if (price && price.trim()) {
				const priceNum = Number(price);
				if (!isNaN(priceNum) && priceNum > 0) {
					payload.price = priceNum;
				}
			}
			if (description.trim()) {
				payload.description = description.trim();
			}

			// If we have a new image file, use multipart/form-data
			// The backend will handle the image upload and set image_url in the DTO
			if (imageFile) {
				await apiClient.postWithImage<{ data?: Product; message: string }>(
					`/products/${id}`,
					imageFile,
					payload,
					"PATCH"
				);
			} else {
				// No new image, just update the product data using regular PATCH
				// Backend requires ID in payload for UpdateProductDTO validation
				// Check if we have any fields to update (besides ID)
				const fieldsToUpdate = Object.keys(payload).filter(key => key !== 'id');
				if (fieldsToUpdate.length > 0) {
					await apiClient.patch<{ data?: Product; message: string }>(
						`/products/${id}`,
						payload
					);
				} else {
					// No fields to update, this is a no-op but not an error
					// Just call onSaved to refresh the list
					console.info("No fields changed, skipping update");
				}
			}
			await onSaved();
		} catch (e) {
			let errorMessage = "Failed to update product";
			if (e instanceof Error) {
				errorMessage = e.message;
				// Try to extract more detailed error from response
				if (e.message.includes("failed to upload image") || e.message.includes("payload") || e.message.includes("validation")) {
					errorMessage = e.message;
				}
			}
			setError(errorMessage);
			console.error("Product update error:", e);
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
					<Field>
						<FieldLabel htmlFor="edit-product-image">Image</FieldLabel>
						<Input
							id="edit-product-image"
							type="file"
							accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
							onChange={handleImageChange}
						/>
						{imagePreview && (
							<div className="mt-2">
								<img src={imagePreview} alt="Preview" className="max-w-full h-32 object-contain border rounded" />
								<p className="text-xs text-muted-foreground mt-1">
									{imageFile ? "New image preview" : "Current image"}
								</p>
							</div>
						)}
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