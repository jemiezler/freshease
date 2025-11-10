import { Button } from "@/components/ui/button";
import { DialogHeader, DialogFooter } from "@/components/ui/dialog";
import { Field, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogTitle } from "@/components/ui/dialog";
import { useState, useEffect } from "react";
import { createResource } from "@/lib/resource";
import { apiClient } from "@/lib/api";
import type { Product } from "@/types/product";
import type { Category, CategoryPayload } from "@/types/catagory";
import type { DialogProps } from "@/types/dialog";

const categories = createResource<Category, CategoryPayload, CategoryPayload>({ basePath: "/categories" });

export function CreateProductDialog({ open, onOpenChange, onSaved }: DialogProps) {
	const [name, setName] = useState("");
	const [sku, setSku] = useState("");
	const [price, setPrice] = useState<string>("");
	const [description, setDescription] = useState("");
	const [imageFile, setImageFile] = useState<File | null>(null);
	const [imagePreview, setImagePreview] = useState<string>("");
	const [unitLabel, setUnitLabel] = useState("kg");
	const [isActive, setIsActive] = useState("active");
	const [categoryId, setCategoryId] = useState<string>("");
	const [quantity, setQuantity] = useState<string>("0");
	const [restockAmount, setRestockAmount] = useState<string>("0");
	const [categoryItems, setCategoryItems] = useState<Category[]>([]);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	// Auto-generate SKU from name
	useEffect(() => {
		if (name && !sku) {
			const generatedSKU = name
				.toUpperCase()
				.replace(/[^A-Z0-9]/g, "-")
				.replace(/-+/g, "-")
				.replace(/^-|-$/g, "")
				.substring(0, 20) + "-" + crypto.randomUUID().substring(0, 8).toUpperCase();
			setSku(generatedSKU);
		}
	}, [name, sku]);

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

		// Validate required fields
		if (!name || !sku || !price || !description || !imageFile || !unitLabel || !isActive || !quantity || !restockAmount) {
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

			const payload = {
				id,
				name,
				sku,
				price: priceNum,
				description,
				unit_label: unitLabel,
				is_active: isActive === "active",
				created_at: now,
				updated_at: now,
				quantity: quantityNum,
				reorder_level: restockAmountNum,
			};

			// Use postWithImage to send both image and product data
			await apiClient.postWithImage<{ data?: Product; message: string }>(
				"/products",
				imageFile,
				payload,
				"POST"
			);

			await onSaved();
			
			// Reset form
			setName("");
			setSku("");
			setPrice("");
			setDescription("");
			setImageFile(null);
			setImagePreview("");
			setUnitLabel("kg");
			setIsActive("active");
			setCategoryId("");
			setQuantity("0");
			setRestockAmount("0");
		} catch (e) {
			let errorMessage = "Failed to create product";
			if (e instanceof Error) {
				errorMessage = e.message;
				// Try to extract more detailed error from response
				if (e.message.includes("failed to upload image") || e.message.includes("payload")) {
					errorMessage = e.message;
				}
			}
			setError(errorMessage);
			console.error("Product creation error:", e);
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
						<FieldLabel htmlFor="product-sku">SKU *</FieldLabel>
						<Input 
							id="product-sku" 
							value={sku} 
							onChange={(e) => setSku(e.target.value)} 
							required 
							placeholder="Auto-generated from name"
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
							required={!imageFile}
						/>
						{imagePreview && (
							<div className="mt-2">
								<img src={imagePreview} alt="Preview" className="max-w-full h-32 object-contain border rounded" />
								<p className="text-xs text-muted-foreground mt-1">Image preview</p>
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