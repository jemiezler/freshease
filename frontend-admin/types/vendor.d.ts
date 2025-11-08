export type Vendor = { 
	id: string; 
	name?: string; 
	email?: string; 
	phone?: string; 
	address?: string; 
	city?: string; 
	state?: string; 
	country?: string; 
	postal_code?: string; 
	website?: string; 
	logo_url?: string; 
	description?: string; 
	is_active?: string;
};

export type VendorPayload = { 
	name?: string; 
	email?: string; 
	phone?: string; 
	address?: string; 
	city?: string; 
	state?: string; 
	country?: string; 
	postal_code?: string; 
	website?: string; 
	logo_url?: string; 
	description?: string; 
	is_active?: string;
};