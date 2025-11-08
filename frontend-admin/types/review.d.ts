export type Review = {
	id: string;
	user_id: string;
	product_id: string;
	rating: number;
	comment?: string | null;
	created_at: string;
};

export type ReviewPayload = {
	id: string;
	user_id: string;
	product_id: string;
	rating: number;
	comment?: string | null;
	created_at: string;
};

