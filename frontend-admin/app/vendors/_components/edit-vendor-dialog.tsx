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
import type { Vendor, VendorPayload } from "@/types/vendor";
import type { EditDialogProps } from "@/types/dialog";

const vendors = createResource<Vendor, VendorPayload, VendorPayload>({
	basePath: "/vendors",
});

export function EditVendorDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [name, setName] = useState("");
	const [email, setEmail] = useState("");
	const [phone, setPhone] = useState("");
	const [address, setAddress] = useState("");
	const [city, setCity] = useState("");
	const [state, setState] = useState("");
	const [country, setCountry] = useState("");
	const [postalCode, setPostalCode] = useState("");
	const [website, setWebsite] = useState("");
	const [logoUrl, setLogoUrl] = useState("");
	const [, setLogoFile] = useState<File | null>(null);
	const [uploadingLogo, setUploadingLogo] = useState(false);
	const [description, setDescription] = useState("");
	const [isActive, setIsActive] = useState("active");
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await vendors.get(id);
				const v = res.data as Vendor | undefined;
				if (!cancelled && v) {
					setName(v.name ?? "");
					setEmail(v.email ?? "");
					setPhone(v.phone ?? "");
					setAddress(v.address ?? "");
					setCity(v.city ?? "");
					setState(v.state ?? "");
					setCountry(v.country ?? "");
					setPostalCode(v.postal_code ?? "");
					setWebsite(v.website ?? "");
					setLogoUrl(v.logo_url ?? "");
					setDescription(v.description ?? "");
					setIsActive(v.is_active ?? "active");
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

	async function handleLogoChange(e: React.ChangeEvent<HTMLInputElement>) {
		const file = e.target.files?.[0];
		if (!file) return;

		setUploadingLogo(true);
		setError(null);

		try {
			const data = await apiClient.uploadImage(file, "vendors/logos");
			setLogoUrl(data.url);
			setLogoFile(file as File);
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to upload logo");
		} finally {
			setUploadingLogo(false);
		}
	}

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: VendorPayload = {
				name: name || undefined,
				email: email || undefined,
				phone: phone || undefined,
				address: address || undefined,
				city: city || undefined,
				state: state || undefined,
				country: country || undefined,
				postal_code: postalCode || undefined,
				website: website || undefined,
				logo_url: logoUrl || undefined,
				description: description || undefined,
				is_active: isActive || undefined,
			};
			await vendors.update(id, payload);
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
					<DialogTitle>Edit Vendor</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading vendor…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-vendor-name">Name</FieldLabel>
							<Input id="edit-vendor-name" value={name} onChange={(e) => setName(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-email">Email</FieldLabel>
							<Input id="edit-vendor-email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-phone">Phone</FieldLabel>
							<Input id="edit-vendor-phone" value={phone} onChange={(e) => setPhone(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-address">Address</FieldLabel>
							<Input id="edit-vendor-address" value={address} onChange={(e) => setAddress(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-city">City</FieldLabel>
							<Input id="edit-vendor-city" value={city} onChange={(e) => setCity(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-state">State</FieldLabel>
							<Input id="edit-vendor-state" value={state} onChange={(e) => setState(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-country">Country</FieldLabel>
							<Input id="edit-vendor-country" value={country} onChange={(e) => setCountry(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-postal-code">Postal Code</FieldLabel>
							<Input id="edit-vendor-postal-code" value={postalCode} onChange={(e) => setPostalCode(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-website">Website</FieldLabel>
							<Input id="edit-vendor-website" type="url" value={website} onChange={(e) => setWebsite(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-logo">Logo</FieldLabel>
							<Input
								id="edit-vendor-logo"
								type="file"
								accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
								onChange={handleLogoChange}
								disabled={uploadingLogo}
							/>
							{uploadingLogo && (
								<div className="flex items-center gap-2 text-sm text-muted-foreground mt-2">
									<Spinner className="size-4" />
									<span>Uploading logo...</span>
								</div>
							)}
							{logoUrl && !uploadingLogo && (
								<div className="mt-2">
									<img src={logoUrl} alt="Logo preview" className="max-w-full h-32 object-contain border rounded" />
									<p className="text-xs text-muted-foreground mt-1">Current logo</p>
								</div>
							)}
							<p className="text-xs text-muted-foreground mt-1">Or enter URL manually:</p>
							<Input
								id="edit-vendor-logo-url"
								type="url"
								value={logoUrl}
								onChange={(e) => setLogoUrl(e.target.value)}
								placeholder="https://..."
								className="mt-1"
							/>
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-description">Description</FieldLabel>
							<Textarea id="edit-vendor-description" value={description} onChange={(e) => setDescription(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-vendor-is-active">Status</FieldLabel>
							<select
								id="edit-vendor-is-active"
								value={isActive}
								onChange={(e) => setIsActive(e.target.value)}
								className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
							>
								<option value="active">Active</option>
								<option value="inactive">Inactive</option>
							</select>
						</Field>
						{error && <p style={{ color: "red" }}>{error}</p>}
						<DialogFooter>
							<Button type="button" variant="secondary" onClick={() => onOpenChange(false)}>
								Cancel
							</Button>
						<Button type="submit" disabled={submitting || uploadingLogo} className="flex items-center gap-2">
							{(submitting || uploadingLogo) && <Spinner className="size-4" />}
							{submitting || uploadingLogo ? "Saving…" : "Save"}
						</Button>
						</DialogFooter>
					</form>
				)}
			</DialogContent>
		</Dialog>
	);
}

