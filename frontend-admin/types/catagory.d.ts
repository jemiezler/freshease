type Category = { id: string; name: string; description?: string; slug: string };
type CategoryPayload = { name: string; description?: string; slug: string };
type Product = {
	id: string;
	name: string;
	description?: string;
	price?: number;
	category_id?: string;
};
type ProductPayload = { name: string; description?: string; price?: number };
