"use client";

import { useState, useEffect } from "react";
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
import type { EditDialogProps } from "@/types/dialog";

const notifications = createResource<Notification, NotificationPayload, NotificationPayload>({
	basePath: "/notifications",
});

export function EditNotificationDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [title, setTitle] = useState("");
	const [body, setBody] = useState("");
	const [channel, setChannel] = useState("");
	const [status, setStatus] = useState("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await notifications.get(id);
				const n = res.data as Notification | undefined;
				if (!cancelled && n) {
					setTitle(n.title ?? "");
					setBody(n.body ?? "");
					setChannel(n.channel ?? "");
					setStatus(n.status ?? "");
				}
			} catch (e) {
				setError(e instanceof Error ? e.message : "Failed to load");
			} finally {
				if (!cancelled) setLoading(false);
			}
		})();
		return () => {
			cancelled = true;
		};
	}, [id]);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: Partial<NotificationPayload> = {
				title,
				body: body || null,
				channel,
				status,
			};
			await notifications.update(id, payload as NotificationPayload);
			await onSaved();
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to update");
		} finally {
			setSubmitting(false);
		}
	}

	return (
		<Dialog open onOpenChange={onOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Edit Notification</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading notification…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-notification-title">Title</FieldLabel>
							<Input id="edit-notification-title" value={title} onChange={(e) => setTitle(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-notification-body">Body</FieldLabel>
							<Textarea id="edit-notification-body" value={body} onChange={(e) => setBody(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-notification-channel">Channel</FieldLabel>
							<Input id="edit-notification-channel" value={channel} onChange={(e) => setChannel(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-notification-status">Status</FieldLabel>
							<Input id="edit-notification-status" value={status} onChange={(e) => setStatus(e.target.value)} required />
						</Field>
						{error && <p style={{ color: "red" }}>{error}</p>}
						<DialogFooter>
							<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
								Cancel
							</Button>
							<Button type="submit" disabled={submitting} className="flex items-center gap-2">
								{submitting && <Spinner className="size-4" />}
								{submitting ? "Updating…" : "Update"}
							</Button>
						</DialogFooter>
					</form>
				)}
			</DialogContent>
		</Dialog>
	);
}

