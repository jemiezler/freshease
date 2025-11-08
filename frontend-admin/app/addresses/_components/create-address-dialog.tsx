"use client";

import { useState } from "react";
import { createResource } from "@/lib/resource";
import { Input } from "@/components/ui/input";
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
import type { Address, AddressPayload } from "@/types/address";
import type { DialogProps } from "@/types/dialog";

const addresses = createResource<Address, AddressPayload, AddressPayload>({
	basePath: "/addresses",
});

export function CreateAddressDialog({
	open,
	onOpenChange,
	onSaved,
}: DialogProps) {
	const [line1, setLine1] = useState("");
	const [line2, setLine2] = useState("");
	const [city, setCity] = useState("");
	const [province, setProvince] = useState("");
	const [country, setCountry] = useState("");
	const [zip, setZip] = useState("");
	const [isDefault, setIsDefault] = useState(false);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	async function onSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSubmitting(true);
		setError(null);
		try {
			const payload: AddressPayload = {
				line1: line1 || undefined,
				line2: line2 || undefined,
				city: city || undefined,
				province: province || undefined,
				country: country || undefined,
				zip: zip || undefined,
				is_default: isDefault,
			};
			await addresses.create(payload);
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
					<DialogTitle>New Address</DialogTitle>
				</DialogHeader>
				<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
					<Field>
						<FieldLabel htmlFor="addr-line1">Line 1 *</FieldLabel>
						<Input id="addr-line1" value={line1} onChange={(e) => setLine1(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="addr-line2">Line 2</FieldLabel>
						<Input id="addr-line2" value={line2} onChange={(e) => setLine2(e.target.value)} />
					</Field>
					<Field>
						<FieldLabel htmlFor="addr-city">City *</FieldLabel>
						<Input id="addr-city" value={city} onChange={(e) => setCity(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="addr-province">Province *</FieldLabel>
						<Input id="addr-province" value={province} onChange={(e) => setProvince(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="addr-country">Country *</FieldLabel>
						<Input id="addr-country" value={country} onChange={(e) => setCountry(e.target.value)} required />
					</Field>
					<Field>
						<FieldLabel htmlFor="addr-zip">Zip *</FieldLabel>
						<Input id="addr-zip" value={zip} onChange={(e) => setZip(e.target.value)} required />
					</Field>
					<Field>
						<div className="flex items-center gap-2">
							<input
								type="checkbox"
								id="addr-is-default"
								checked={isDefault}
								onChange={(e) => setIsDefault(e.target.checked)}
								className="h-4 w-4 rounded border-gray-300"
							/>
							<FieldLabel htmlFor="addr-is-default" className="cursor-pointer">Set as default address</FieldLabel>
						</div>
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
