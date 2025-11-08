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
import type { Role, RolePayload } from "@/types/role";
import type { EditDialogProps } from "@/types/dialog";

const roles = createResource<Role, RolePayload, RolePayload>({
	basePath: "/roles",
});

export function EditRoleDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await roles.get(id);
				const r = res.data as Role | undefined;
				if (!cancelled && r) {
					setName(r.name ?? "");
					setDescription(r.description ?? "");
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
			const payload: RolePayload = { 
				name: name || undefined,
				description: description || undefined,
			};
			await roles.update(id, payload);
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
					<DialogTitle>Edit Role</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading role…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-role-name">Name</FieldLabel>
							<Input id="edit-role-name" value={name} onChange={(e) => setName(e.target.value)} required />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-role-description">Description</FieldLabel>
							<Textarea id="edit-role-description" value={description} onChange={(e) => setDescription(e.target.value)} required />
						</Field>
						{error && <p style={{ color: "red" }}>{error}</p>}
						<DialogFooter>
							<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
								Cancel
							</Button>
							<Button type="submit" disabled={submitting} className="flex items-center gap-2">
								{submitting && <Spinner className="size-4" />}
								{submitting ? "Saving…" : "Save"}
							</Button>
						</DialogFooter>
					</form>
				)}
			</DialogContent>
		</Dialog>
	);
}
