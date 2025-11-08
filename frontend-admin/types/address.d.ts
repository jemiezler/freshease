export type Address = { 
	id: string; 
	line1?: string; 
	line2?: string; 
	city?: string; 
	province?: string; 
	country?: string; 
	zip?: string; 
	is_default?: boolean;
};

export type AddressPayload = { 
	line1?: string; 
	line2?: string; 
	city?: string; 
	province?: string; 
	country?: string; 
	zip?: string; 
	is_default?: boolean;
};

