export type MealPlan = {
	id: string;
	user_id: string;
	week_start: string;
	goal?: string | null;
};

export type MealPlanPayload = {
	id: string;
	user_id: string;
	week_start: string;
	goal?: string | null;
};

