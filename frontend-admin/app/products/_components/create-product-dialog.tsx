import { Button } from "@/components/ui/button";
import { DialogHeader, DialogFooter } from "@/components/ui/dialog";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { useState, useEffect } from "react";
import { createResource } from "@/lib/resource";
import type { Product, ProductPayload } from "@/types/product";
import type { Category, CategoryPayload } from "@/types/catagory";
import type { DialogProps } from "@/types/dialog";

const products = createResource<Product, ProductPayload, ProductPayload>({ basePath: "/products" });
const categories = createResource<Category, CategoryPayload, CategoryPayload>({ basePath: "/product_categories" });

export function CreateProductDialog({ open, onOpenChange, onSaved }: DialogProps) {
	const [name, setName] = useState("");
	const [price, setPrice] = useState<string>("");
	const [description, setDescription] = useState("");
	const [imageUrl, setImageUrl] = useState("");
	const [unitLabel, setUnitLabel] = useState("kg");
	const [isActive, setIsActive] = useState("active");
	const [categoryId, setCategoryId] = useState<string>("");
	const [quantity, setQuantity] = useState<string>("0");
	const [restockAmount, setRestockAmount] = useState<string>("0");
	const [uploadingImage, setUploadingImage] = useState(false);
	const [categoryItems, setCategoryItems] = useState<Category[]>([]);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		if (open) {
			// Load categories when dialog opens
			categories.list().then((res) => {
				setCategoryItems(res.data ?? []);
			}).catch(() => {
				// Ignore errors
			});
		}
	}, [open]);

	async function handleImageChange(e: React.ChangeEvent<HTMLInputElement>) {
		const file = e.target.files?.[0];
		if (!file) return;

		setUploadingImage(true);
		setError(null);

		try {
			const formData = new FormData();
			formData.append("file", file);
			formData.append("folder", "products");

			const token = localStorage.getItem("admin_token");
			const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";
			const headers: HeadersInit = {};
			if (token) {
				headers["Authorization"] = `Bearer ${token}`;
			}

			const response = await fetch(`${baseUrl}/uploads/images`, {
				method: "POST",
				headers,
				body: formData,
			});

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: "Upload failed" }));
				throw new Error(errorData.message || "Failed to upload image");
			}

			const data = await response.json();
			setImageUrl(data.url || data.object_name);
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

		// Validate required fields
		if (!name || !price || !description || !imageUrl || !unitLabel || !isActive || !quantity || !restockAmount) {
			setError("Please fill in all required fields");
			setSubmitting(false);
			return;
		}

		const priceNum = Number(price);
		if (isNaN(priceNum) || priceNum <= 0) {
			setError("Price must be a positive number");
			setSubmitting(false);
			return;
		}

		const quantityNum = Number(quantity);
		if (isNaN(quantityNum) || quantityNum <= 0) {
			setError("Quantity must be a positive number");
			setSubmitting(false);
			return;
		}

		const restockAmountNum = Number(restockAmount);
		if (isNaN(restockAmountNum) || restockAmountNum <= 0) {
			setError("Restock amount must be a positive number");
			setSubmitting(false);
			return;
		}

		try {
			// Generate UUID and timestamps
			const id = crypto.randomUUID();
			const now = new Date().toISOString();

			const payload: ProductPayload = {
				id,
				name,
				price: priceNum,
				description,
				image_url: imageUrl,
				unit_label: unitLabel,
				is_active: isActive,
				created_at: now,
				updated_at: now,
				quantity: quantityNum,
				restock_amount: restockAmountNum,
			};

			await products.create(payload);
			await onSaved();
			
			// Reset form
			setName("");
			setPrice("");
			setDescription("");
			setImageUrl("");
			setUnitLabel("kg");
			setIsActive("active");
			setCategoryId("");
			setQuantity("0");
			setRestockAmount("0");
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to create product");
		} finally {
			setSubmitting(false);
		}
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent style={{ maxWidth: "600px", maxHeight: "90vh", overflowY: "auto" }}>
				<DialogHeader>
					<DialogTitle>New Product</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="product-name">Name *</FieldLabel>
						<Input 
							id="product-name" 
							value={name} 
							onChange={(e) => setName(e.target.value)} 
							required 
							minLength={2}
							maxLength={60}
						/>
					</Field>
					<Field>
						<FieldLabel htmlFor="product-price">Price *</FieldLabel>
						<Input 
							id="product-price" 
							type="number" 
							step="0.01" 
							min="0.01"
							value={price} 
							onChange={(e) => setPrice(e.target.value)} 
							required 
						/>
					</Field>
					<Field>
						<FieldLabel htmlFor="product-description">Description *</FieldLabel>
						<Textarea 
							id="product-description" 
							value={description} 
							onChange={(e) => setDescription(e.target.value)} 
							required 
						/>
					</Field>
					<Field>
						<FieldLabel htmlFor="product-image">Image *</FieldLabel>
						<Input 
							id="product-image" 
							type="file" 
							accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
							onChange={handleImageChange}
							disabled={uploadingImage}
							required={!imageUrl}
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
								<p className="text-xs text-muted-foreground mt-1">Image uploaded successfully</p>
							</div>
						)}
					</Field>
					<Field>
						<FieldLabel htmlFor="product-unit-label">Unit Label *</FieldLabel>
						<Input 
							id="product-unit-label" 
							value={unitLabel} 
							onChange={(e) => setUnitLabel(e.target.value)} 
							placeholder="e.g., kg, lb, piece"
							required 
						/>
					</Field>
					<Field>
						<FieldLabel htmlFor="product-is-active">Status *</FieldLabel>
						<select
							id="product-is-active"
							value={isActive}
							onChange={(e) => setIsActive(e.target.value)}
							required
							className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							<option value="active">Active</option>
							<option value="inactive">Inactive</option>
						</select>
					</Field>
					<Field>
						<FieldLabel htmlFor="product-category">Category</FieldLabel>
						<select
							id="product-category"
							value={categoryId}
							onChange={(e) => setCategoryId(e.target.value)}
							className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							<option value="">None</option>
							{categoryItems.map((cat) => (
								<option key={cat.id} value={cat.id}>
									{cat.name}
								</option>
							))}
						</select>
					</Field>
					<Field>
						<FieldLabel htmlFor="product-quantity">Initial Quantity *</FieldLabel>
						<Input 
							id="product-quantity" 
							type="number" 
							min="1"
							step="1"
							value={quantity} 
							onChange={(e) => setQuantity(e.target.value)} 
							placeholder="e.g., 100"
							required 
						/>
					</Field>
					<Field>
						<FieldLabel htmlFor="product-restock-amount">Restock Amount *</FieldLabel>
						<Input 
							id="product-restock-amount" 
							type="number" 
							min="1"
							step="1"
							value={restockAmount} 
							onChange={(e) => setRestockAmount(e.target.value)} 
							placeholder="e.g., 50"
							required 
						/>
						<p className="text-xs text-muted-foreground mt-1">Amount to restock when inventory is low</p>
					</Field>
					{error && <p style={{ color: "red" }}>{error}</p>}
					<DialogFooter>
						<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
							Cancel
						</Button>
						<Button type="submit" disabled={submitting || uploadingImage} className="flex items-center gap-2">
							{submitting && <Spinner className="size-4" />}
							{submitting ? "Saving…" : "Create"}
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
}