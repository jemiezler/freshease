export type Payment = {
	id: string;
	provider: string;
	provider_ref?: string | null;
	status: string;
	amount: number;
	paid_at?: string | null;
	order_id: string;
};

export type PaymentPayload = {
	id: string;
	provider: string;
	provider_ref?: string | null;
	status: string;
	amount: number;
	paid_at?: string | null;
	order_id: string;
};

