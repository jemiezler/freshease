export type Product = {
	id: string;
	name: string;
	price: number;
	description: string;
	image_url: string;
	unit_label: string;
	is_active: string;
	created_at: string;
	updated_at: string;
	deleted_at?: string | null;
	category_id?: string;
};

export type ProductPayload = {
	id: string;
	name: string;
	price: number;
	description: string;
	image_url: string;
	unit_label: string;
	is_active: string;
	created_at: string;
	updated_at: string;
	quantity: number;
	restock_amount: number;
};

