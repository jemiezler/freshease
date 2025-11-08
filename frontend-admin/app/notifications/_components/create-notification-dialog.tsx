"use client";

import { useState } from "react";
import { createResource } from "@/lib/resource";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Field, FieldLabel } from "@/components/ui/field";
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import type { Notification, NotificationPayload } from "@/types/notification";
import type { DialogProps } from "@/types/dialog";
import { generateUUID } from "@/lib/utils";

const notifications = createResource<Notification, NotificationPayload, NotificationPayload>({
	basePath: "/notifications",
});

export function CreateNotificationDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [title, setTitle] = useState("");
	const [body, setBody] = useState("");
	const [channel, setChannel] = useState("");
	const [status, setStatus] = useState("");
	const [userId, setUserId] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: NotificationPayload = {
				id: generateUUID(),
				title,
				body: body || null,
				channel,
				status,
				user_id: userId,
				created_at: new Date().toISOString(),
			};
			await notifications.create(payload);
			await onSaved();
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to create");
		} finally {
			setSubmitting(false);
		}
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent style={{ maxWidth: "600px" }}>
				<DialogHeader>
					<DialogTitle>New Notification</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="notification-title">Title *</FieldLabel>
						<Input id="notification-title" value={title} onChange={(e) => setTitle(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="notification-body">Body</FieldLabel>
						<Textarea id="notification-body" value={body} onChange={(e) => setBody(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="notification-channel">Channel *</FieldLabel>
						<Input id="notification-channel" value={channel} onChange={(e) => setChannel(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="notification-status">Status *</FieldLabel>
						<Input id="notification-status" value={status} onChange={(e) => setStatus(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="notification-user-id">User ID *</FieldLabel>
						<Input id="notification-user-id" value={userId} onChange={(e) => setUserId(e.target.value)} required />
					</Field>
					{error && <p style={{ color: "red" }}>{error}</p>}
					<DialogFooter>
						<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
							Cancel
						</Button>
						<Button type="submit" disabled={submitting} className="flex items-center gap-2">
							{submitting && <Spinner className="size-4" />}
							{submitting ? "Creating…" : "Create"}
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
}

