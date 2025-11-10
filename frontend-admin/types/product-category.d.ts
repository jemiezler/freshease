export type ProductCategory = {
	id: string;
	product_id: string;
	category_id: string;
};

export type ProductCategoryPayload = {
	id: string;
	product_id: string;
	category_id: string;
};

export type UpdateProductCategoryPayload = {
	id: string;
	product_id?: string;
	category_id?: string;
};

