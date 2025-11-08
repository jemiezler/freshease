export type Order = {
	id: string;
	order_no: string;
	status: string;
	subtotal: number;
	shipping_fee: number;
	discount: number;
	total: number;
	placed_at?: string | null;
	updated_at: string;
	user_id: string;
	shipping_address_id?: string | null;
	billing_address_id?: string | null;
};

export type OrderPayload = {
	id: string;
	order_no: string;
	status: string;
	subtotal: number;
	shipping_fee: number;
	discount: number;
	total: number;
	placed_at?: string | null;
	user_id: string;
	shipping_address_id?: string | null;
	billing_address_id?: string | null;
};

