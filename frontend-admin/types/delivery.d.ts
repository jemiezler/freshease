export type Delivery = {
	id: string;
	provider: string;
	tracking_no?: string | null;
	status: string;
	eta?: string | null;
	delivered_at?: string | null;
	order_id: string;
};

export type DeliveryPayload = {
	id: string;
	provider: string;
	tracking_no?: string | null;
	status: string;
	eta?: string | null;
	delivered_at?: string | null;
	order_id: string;
};

