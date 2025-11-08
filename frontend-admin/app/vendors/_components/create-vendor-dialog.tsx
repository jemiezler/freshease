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
import type { Vendor, VendorPayload } from "@/types/vendor";
import type { DialogProps } from "@/types/dialog";

const vendors = createResource<Vendor, VendorPayload, VendorPayload>({
	basePath: "/vendors",
});

export function CreateVendorDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
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
	const [description, setDescription] = useState("");
	const [isActive, setIsActive] = useState("active");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

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
				is_active: isActive,
			};
			await vendors.create(payload);
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
					<DialogTitle>New Vendor</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="vendor-name">Name</FieldLabel>
						<Input id="vendor-name" value={name} onChange={(e) => setName(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-email">Email</FieldLabel>
						<Input id="vendor-email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-phone">Phone</FieldLabel>
						<Input id="vendor-phone" value={phone} onChange={(e) => setPhone(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-address">Address</FieldLabel>
						<Input id="vendor-address" value={address} onChange={(e) => setAddress(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-city">City</FieldLabel>
						<Input id="vendor-city" value={city} onChange={(e) => setCity(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-state">State</FieldLabel>
						<Input id="vendor-state" value={state} onChange={(e) => setState(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-country">Country</FieldLabel>
						<Input id="vendor-country" value={country} onChange={(e) => setCountry(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-postal-code">Postal Code</FieldLabel>
						<Input id="vendor-postal-code" value={postalCode} onChange={(e) => setPostalCode(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-website">Website</FieldLabel>
						<Input id="vendor-website" type="url" value={website} onChange={(e) => setWebsite(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-logo-url">Logo URL</FieldLabel>
						<Input id="vendor-logo-url" type="url" value={logoUrl} onChange={(e) => setLogoUrl(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-description">Description</FieldLabel>
						<Textarea id="vendor-description" value={description} onChange={(e) => setDescription(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="vendor-is-active">Status</FieldLabel>
						<select
							id="vendor-is-active"
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

