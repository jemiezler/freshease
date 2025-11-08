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
	const [imageUrl, setImageUrl] = useState("");
	const [imageFile, setImageFile] = useState<File | null>(null);
	const [uploadingImage, setUploadingImage] = useState(false);
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
				}
			} catch (e) {
				setError(e instanceof Error ? e.message : "Failed to load");
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => { cancelled = true; };
	}, [id]);

	async function handleImageChange(e: React.ChangeEvent<HTMLInputElement>) {
		const file = e.target.files?.[0];
		if (!file) return;

		setUploadingImage(true);
		setError(null);

		try {
			const data = await apiClient.uploadImage(file, "products");
			setImageUrl(data.url);
			setImageFile(file);
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to upload image");
		} finally {
			setUploadingImage(false);
		}
	}

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: Partial<ProductPayload> = {
				name: name || undefined,
				price: price ? Number(price) : undefined,
				description: description || undefined,
				image_url: imageUrl || undefined,
			};

			// If we have a new image file, use multipart/form-data
			if (imageFile) {
				await apiClient.postWithImage<{ data?: Product; message: string }>(
					`/products/${id}`,
					imageFile,
					payload,
					"PATCH"
				);
			} else {
				await products.update(id, payload);
			}
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
					<Field>
						<FieldLabel htmlFor="edit-product-image">Image</FieldLabel>
						<Input
							id="edit-product-image"
							type="file"
							accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
							onChange={handleImageChange}
							disabled={uploadingImage}
						/>
						{uploadingImage && (
							<div className="flex items-center gap-2 text-sm text-muted-foreground mt-2">
								<Spinner className="size-4" />
								<span>Uploading image...</span>
							</div>
						)}
						{imageUrl && !uploadingImage && (
							<div className="mt-2">
								<img src={imageUrl} alt="Preview" className="max-w-full h-32 object-contain border rounded" />
								<p className="text-xs text-muted-foreground mt-1">Current image</p>
							</div>
						)}
					</Field>
					{error && <p style={{ color: "red" }}>{error}</p>}
						<DialogFooter>
							<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
								Cancel
							</Button>
						<Button type="submit" disabled={submitting || uploadingImage} className="flex items-center gap-2">
							{(submitting || uploadingImage) && <Spinner className="size-4" />}
							{submitting || uploadingImage ? "Saving…" : "Save"}
						</Button>
						</DialogFooter>
					</form>
				)}
			</DialogContent>
		</Dialog>
	);
}