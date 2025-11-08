export type Recipe = {
	id: string;
	name: string;
	instructions?: string | null;
	kcal: number;
};

export type RecipePayload = {
	id: string;
	name: string;
	instructions?: string | null;
	kcal: number;
};

