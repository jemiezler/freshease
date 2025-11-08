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
import type { Role, RolePayload } from "@/types/role";
import type { DialogProps } from "@/types/dialog";

const roles = createResource<Role, RolePayload, RolePayload>({
	basePath: "/roles",
});

export function CreateRoleDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: RolePayload = { 
				name: name || undefined,
				description: description || undefined,
			};
			await roles.create(payload);
			await onSaved();
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to create");
		} finally {
			setSubmitting(false);
		}
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>New Role</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="role-name">Name *</FieldLabel>
						<Input id="role-name" value={name} onChange={(e) => setName(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="role-description">Description *</FieldLabel>
						<Textarea id="role-description" value={description} onChange={(e) => setDescription(e.target.value)} required />
					</Field>
					{error && <p style={{ color: "red" }}>{error}</p>}
					<DialogFooter>
						<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
							Cancel
						</Button>
						<Button type="submit" disabled={submitting} className="flex items-center gap-2">
							{submitting && <Spinner className="size-4" />}
							{submitting ? "Saving…" : "Create"}
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
}
