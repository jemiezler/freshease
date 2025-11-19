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
import type { Category, CategoryPayload } from "@/types/catagory";
import type { ProductCategory } from "@/types/product-category";
import type { EditDialogProps } from "@/types/dialog";

const products = createResource<Product, ProductPayload, ProductPayload>({ basePath: "/products" });
const categories = createResource<Category, CategoryPayload, CategoryPayload>({ basePath: "/categories" });
const productCategories = createResource<ProductCategory, ProductCategory, ProductCategory>({ basePath: "/product_categories" });

export function EditProductDialog({ id, onOpenChange, onSaved }: EditDialogProps) {
	const [name, setName] = useState("");
	const [price, setPrice] = useState<string>("");
	const [description, setDescription] = useState("");
	const [, setImageUrl] = useState("");
	const [imageFile, setImageFile] = useState<File | null>(null);
	const [imagePreview, setImagePreview] = useState<string>("");
	const [selectedCategoryIds, setSelectedCategoryIds] = useState<string[]>([]);
	const [categoryItems, setCategoryItems] = useState<Category[]>([]);
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				// Load product, categories, and product categories in parallel
				const [productRes, categoriesRes, productCategoriesRes] = await Promise.all([
					products.get(id),
					categories.list(),
					productCategories.list(),
				]);

				const p = productRes.data as Product | undefined;
				if (!cancelled && p) {
					setName(p.name ?? "");
					setPrice(p.price != null ? String(p.price) : "");
					setDescription(p.description ?? "");
					setImageUrl(p.image_url ?? "");
					setImagePreview(p.image_url ?? "");
				}

				// Set available categories
				if (!cancelled) {
					setCategoryItems(categoriesRes.data ?? []);
				}

				// Set selected categories for this product
				if (!cancelled) {
					const productCats = (productCategoriesRes.data ?? []).filter(
						(pc) => pc.product_id === id
					);
					setSelectedCategoryIds(productCats.map((pc) => pc.category_id));
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
			// Update product categories if they changed
			// First, get current product categories
			const currentProductCatsRes = await productCategories.list();
			const currentProductCats = (currentProductCatsRes.data ?? []).filter(
				(pc) => pc.product_id === id
			);
			const currentCategoryIds = currentProductCats.map((pc) => pc.category_id);

			// Find categories to add and remove
			const toAdd = selectedCategoryIds.filter((cid) => !currentCategoryIds.includes(cid));
			const toRemove = currentCategoryIds.filter((cid) => !selectedCategoryIds.includes(cid));

			// Add new product categories
			for (const categoryId of toAdd) {
				try {
					await productCategories.create({
						id: crypto.randomUUID(),
						product_id: id,
						category_id: categoryId,
					});
				} catch (e) {
					console.error(`Failed to add category ${categoryId}:`, e);
				}
			}

			// Remove old product categories
			for (const categoryId of toRemove) {
				const pcToDelete = currentProductCats.find((pc) => pc.category_id === categoryId);
				if (pcToDelete) {
					try {
						await productCategories.delete(pcToDelete.id);
					} catch (e) {
						console.error(`Failed to remove category ${categoryId}:`, e);
					}
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
								<img src={`${imagePreview}`} alt="Preview" className="max-w-full h-32 object-contain border rounded" />
								<p className="text-xs text-muted-foreground mt-1">
									{imageFile ? "New image preview" : "Current image"}
								</p>
							</div>
						)}
					</Field>
					<Field>
						<FieldLabel htmlFor="edit-product-categories">Categories</FieldLabel>
						<div className="space-y-2 max-h-48 overflow-y-auto border rounded-md p-3">
							{categoryItems.length === 0 ? (
								<p className="text-sm text-muted-foreground">No categories available</p>
							) : (
								categoryItems.map((cat) => (
									<label
										key={cat.id}
										className="flex items-center gap-2 cursor-pointer hover:bg-muted/50 p-2 rounded"
									>
										<input
											type="checkbox"
											checked={selectedCategoryIds.includes(cat.id)}
											onChange={(e) => {
												if (e.target.checked) {
													setSelectedCategoryIds([...selectedCategoryIds, cat.id]);
												} else {
													setSelectedCategoryIds(selectedCategoryIds.filter((id) => id !== cat.id));
												}
											}}
											className="h-4 w-4 rounded border-gray-300"
										/>
										<span className="text-sm">{cat.name}</span>
									</label>
								))
							)}
						</div>
						{selectedCategoryIds.length > 0 && (
							<p className="text-xs text-muted-foreground mt-1">
								{selectedCategoryIds.length} categor{selectedCategoryIds.length === 1 ? "y" : "ies"} selected
							</p>
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