export type Notification = {
	id: string;
	user_id: string;
	title: string;
	body?: string | null;
	channel: string;
	status: string;
	created_at: string;
};

export type NotificationPayload = {
	id: string;
	user_id: string;
	title: string;
	body?: string | null;
	channel: string;
	status: string;
	created_at: string;
};

