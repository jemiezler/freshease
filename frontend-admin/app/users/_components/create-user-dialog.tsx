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
import type { User, UserPayload } from "@/types/user";
import type { DialogProps } from "@/types/dialog";

const users = createResource<User, UserPayload, UserPayload>({
	basePath: "/users",
	updateMethod: "PUT",
});

export function CreateUserDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [email, setEmail] = useState("");
	const [name, setName] = useState("");
	const [password, setPassword] = useState("");
	const [phone, setPhone] = useState("");
	const [bio, setBio] = useState("");
	const [avatar, setAvatar] = useState("");
	const [cover, setCover] = useState("");
	const [dateOfBirth, setDateOfBirth] = useState("");
	const [sex, setSex] = useState("");
	const [goal, setGoal] = useState("");
	const [heightCm, setHeightCm] = useState<string>("");
	const [weightKg, setWeightKg] = useState<string>("");
	const [status, setStatus] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: UserPayload = {
				email: email || undefined,
				name: name || undefined,
				password: password || undefined,
				phone: phone || undefined,
				bio: bio || undefined,
				avatar: avatar || undefined,
				cover: cover || undefined,
				date_of_birth: dateOfBirth || undefined,
				sex: sex || undefined,
				goal: goal || undefined,
				height_cm: heightCm ? Number(heightCm) : undefined,
				weight_kg: weightKg ? Number(weightKg) : undefined,
				status: status || undefined,
			};
			await users.create(payload);
			await onSaved();
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to create");
		} finally {
			setSubmitting(false);
		}
	}

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent style={{ maxWidth: "600px", maxHeight: "90vh", overflowY: "auto" }}>
				<DialogHeader>
					<DialogTitle>New User</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="user-email">Email *</FieldLabel>
						<Input id="user-email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-name">Name *</FieldLabel>
						<Input id="user-name" value={name} onChange={(e) => setName(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-password">Password *</FieldLabel>
						<Input id="user-password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} required minLength={8} maxLength={100} />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-phone">Phone</FieldLabel>
						<Input id="user-phone" value={phone} onChange={(e) => setPhone(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-bio">Bio</FieldLabel>
						<Textarea id="user-bio" value={bio} onChange={(e) => setBio(e.target.value)} maxLength={500} />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-avatar">Avatar URL</FieldLabel>
						<Input id="user-avatar" type="url" value={avatar} onChange={(e) => setAvatar(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-cover">Cover URL</FieldLabel>
						<Input id="user-cover" type="url" value={cover} onChange={(e) => setCover(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-date-of-birth">Date of Birth</FieldLabel>
						<Input id="user-date-of-birth" type="date" value={dateOfBirth} onChange={(e) => setDateOfBirth(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-sex">Sex</FieldLabel>
						<select
							id="user-sex"
							value={sex}
							onChange={(e) => setSex(e.target.value)}
							className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							<option value="">Select...</option>
							<option value="male">Male</option>
							<option value="female">Female</option>
							<option value="other">Other</option>
						</select>
					</Field>
					<Field>
						<FieldLabel htmlFor="user-goal">Goal</FieldLabel>
						<select
							id="user-goal"
							value={goal}
							onChange={(e) => setGoal(e.target.value)}
							className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
						>
							<option value="">Select...</option>
							<option value="maintenance">Maintenance</option>
							<option value="weight_loss">Weight Loss</option>
							<option value="weight_gain">Weight Gain</option>
						</select>
					</Field>
					<Field>
						<FieldLabel htmlFor="user-height-cm">Height (cm)</FieldLabel>
						<Input id="user-height-cm" type="number" step="0.1" min="50" max="300" value={heightCm} onChange={(e) => setHeightCm(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-weight-kg">Weight (kg)</FieldLabel>
						<Input id="user-weight-kg" type="number" step="0.1" min="20" max="500" value={weightKg} onChange={(e) => setWeightKg(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="user-status">Status</FieldLabel>
						<Input id="user-status" value={status} onChange={(e) => setStatus(e.target.value)} />
					</Field>
					{error && <p style={{ color: "red" }}>{error}</p>}
					<DialogFooter>
						<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>Cancel</Button>
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
