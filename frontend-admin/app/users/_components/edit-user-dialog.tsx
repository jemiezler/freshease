"use client";

import { useState, useEffect } from "react";
import { createResource } from "@/lib/resource";
import { apiClient } from "@/lib/api";
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
import type { EditDialogProps } from "@/types/dialog";

const users = createResource<User, UserPayload, UserPayload>({
	basePath: "/users",
	updateMethod: "PUT",
});

export function EditUserDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [email, setEmail] = useState("");
	const [name, setName] = useState("");
	const [password, setPassword] = useState("");
	const [phone, setPhone] = useState("");
	const [bio, setBio] = useState("");
	const [avatar, setAvatar] = useState("");
	const [cover, setCover] = useState("");
	const [avatarFile, setAvatarFile] = useState<File | null>(null);
	const [coverFile, setCoverFile] = useState<File | null>(null);
	const [uploadingAvatar, setUploadingAvatar] = useState(false);
	const [uploadingCover, setUploadingCover] = useState(false);
	const [dateOfBirth, setDateOfBirth] = useState("");
	const [sex, setSex] = useState("");
	const [goal, setGoal] = useState("");
	const [heightCm, setHeightCm] = useState<string>("");
	const [weightKg, setWeightKg] = useState<string>("");
	const [status, setStatus] = useState("");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await users.get(id);
				const u = res.data as User | undefined;
				if (!cancelled && u) {
					setEmail(u.email ?? "");
					setName(u.name ?? "");
					setPhone(u.phone ?? "");
					setBio(u.bio ?? "");
					setAvatar(u.avatar ?? "");
					setCover(u.cover ?? "");
					setDateOfBirth(u.date_of_birth ?? "");
					setSex(u.sex ?? "");
					setGoal(u.goal ?? "");
					setHeightCm(u.height_cm != null ? String(u.height_cm) : "");
					setWeightKg(u.weight_kg != null ? String(u.weight_kg) : "");
					setStatus(u.status ?? "");
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

	async function handleAvatarChange(e: React.ChangeEvent<HTMLInputElement>) {
		const file = e.target.files?.[0];
		if (!file) return;

		setUploadingAvatar(true);
		setError(null);

		try {
			const data = await apiClient.uploadImage(file, "users/avatars");
			setAvatar(data.url);
			setAvatarFile(file);
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to upload avatar");
		} finally {
			setUploadingAvatar(false);
		}
	}

	async function handleCoverChange(e: React.ChangeEvent<HTMLInputElement>) {
		const file = e.target.files?.[0];
		if (!file) return;

		setUploadingCover(true);
		setError(null);

		try {
			const data = await apiClient.uploadImage(file, "users/covers");
			setCover(data.url);
			setCoverFile(file);
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to upload cover");
		} finally {
			setUploadingCover(false);
		}
	}

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
			await users.update(id, payload);
			await onSaved();
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to update");
		} finally {
			setSubmitting(false);
		}
	}

	return (
		<Dialog open onOpenChange={onOpenChange}>
			<DialogContent style={{ maxWidth: "600px", maxHeight: "90vh", overflowY: "auto" }}>
				<DialogHeader>
					<DialogTitle>Edit User</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading user…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-user-email">Email</FieldLabel>
							<Input id="edit-user-email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-name">Name</FieldLabel>
							<Input id="edit-user-name" value={name} onChange={(e) => setName(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-password">Password (leave blank to keep current)</FieldLabel>
							<Input id="edit-user-password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} minLength={8} maxLength={100} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-phone">Phone</FieldLabel>
							<Input id="edit-user-phone" value={phone} onChange={(e) => setPhone(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-bio">Bio</FieldLabel>
							<Textarea id="edit-user-bio" value={bio} onChange={(e) => setBio(e.target.value)} maxLength={500} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-avatar">Avatar</FieldLabel>
							<Input
								id="edit-user-avatar"
								type="file"
								accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
								onChange={handleAvatarChange}
								disabled={uploadingAvatar}
							/>
							{uploadingAvatar && (
								<div className="flex items-center gap-2 text-sm text-muted-foreground mt-2">
									<Spinner className="size-4" />
									<span>Uploading avatar...</span>
								</div>
							)}
							{avatar && !uploadingAvatar && (
								<div className="mt-2">
									<img src={avatar} alt="Avatar preview" className="max-w-full h-32 object-contain border rounded" />
									<p className="text-xs text-muted-foreground mt-1">Current avatar</p>
								</div>
							)}
							<p className="text-xs text-muted-foreground mt-1">Or enter URL manually:</p>
							<Input
								id="edit-user-avatar-url"
								type="url"
								value={avatar}
								onChange={(e) => setAvatar(e.target.value)}
								placeholder="https://..."
								className="mt-1"
							/>
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-cover">Cover</FieldLabel>
							<Input
								id="edit-user-cover"
								type="file"
								accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
								onChange={handleCoverChange}
								disabled={uploadingCover}
							/>
							{uploadingCover && (
								<div className="flex items-center gap-2 text-sm text-muted-foreground mt-2">
									<Spinner className="size-4" />
									<span>Uploading cover...</span>
								</div>
							)}
							{cover && !uploadingCover && (
								<div className="mt-2">
									<img src={cover} alt="Cover preview" className="max-w-full h-32 object-contain border rounded" />
									<p className="text-xs text-muted-foreground mt-1">Current cover</p>
								</div>
							)}
							<p className="text-xs text-muted-foreground mt-1">Or enter URL manually:</p>
							<Input
								id="edit-user-cover-url"
								type="url"
								value={cover}
								onChange={(e) => setCover(e.target.value)}
								placeholder="https://..."
								className="mt-1"
							/>
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-date-of-birth">Date of Birth</FieldLabel>
							<Input id="edit-user-date-of-birth" type="date" value={dateOfBirth} onChange={(e) => setDateOfBirth(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-sex">Sex</FieldLabel>
							<select
								id="edit-user-sex"
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
							<FieldLabel htmlFor="edit-user-goal">Goal</FieldLabel>
							<select
								id="edit-user-goal"
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
							<FieldLabel htmlFor="edit-user-height-cm">Height (cm)</FieldLabel>
							<Input id="edit-user-height-cm" type="number" step="0.1" min="50" max="300" value={heightCm} onChange={(e) => setHeightCm(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-weight-kg">Weight (kg)</FieldLabel>
							<Input id="edit-user-weight-kg" type="number" step="0.1" min="20" max="500" value={weightKg} onChange={(e) => setWeightKg(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-user-status">Status</FieldLabel>
							<Input id="edit-user-status" value={status} onChange={(e) => setStatus(e.target.value)} />
						</Field>
						{error && <p style={{ color: "red" }}>{error}</p>}
						<DialogFooter>
							<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>Cancel</Button>
						<Button type="submit" disabled={submitting || uploadingAvatar || uploadingCover} className="flex items-center gap-2">
							{(submitting || uploadingAvatar || uploadingCover) && <Spinner className="size-4" />}
							{submitting || uploadingAvatar || uploadingCover ? "Saving…" : "Save"}
						</Button>
						</DialogFooter>
					</form>
				)}
			</DialogContent>
		</Dialog>
	);
}
