export type Bundle = {
	id: string;
	name: string;
	description?: string | null;
	price: number;
	is_active: boolean;
};

export type BundlePayload = {
	id: string;
	name: string;
	description?: string | null;
	price: number;
	is_active: boolean;
};

