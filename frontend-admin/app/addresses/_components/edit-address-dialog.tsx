"use client";

import { useState, useEffect } from "react";
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
import type { EditDialogProps } from "@/types/dialog";

const addresses = createResource<Address, AddressPayload, AddressPayload>({
	basePath: "/addresses",
});

export function EditAddressDialog({
	id,
	onOpenChange,
	onSaved,
}: EditDialogProps) {
	const [line1, setLine1] = useState("");
	const [line2, setLine2] = useState("");
	const [city, setCity] = useState("");
	const [province, setProvince] = useState("");
	const [country, setCountry] = useState("");
	const [zip, setZip] = useState("");
	const [isDefault, setIsDefault] = useState(false);
	const [loading, setLoading] = useState(true);
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		let cancelled = false;
		(async () => {
			try {
				const res = await addresses.get(id);
				const a = res.data as Address | undefined;
				if (!cancelled && a) {
					setLine1(a.line1 ?? "");
					setLine2(a.line2 ?? "");
					setCity(a.city ?? "");
					setProvince(a.province ?? "");
					setCountry(a.country ?? "");
					setZip(a.zip ?? "");
					setIsDefault(a.is_default ?? false);
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
			const payload: AddressPayload = {
				line1: line1 || undefined,
				line2: line2 || undefined,
				city: city || undefined,
				province: province || undefined,
				country: country || undefined,
				zip: zip || undefined,
				is_default: isDefault,
			};
			await addresses.update(id, payload);
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
					<DialogTitle>Edit Address</DialogTitle>
				</DialogHeader>
				{loading ? (
					<div className="flex items-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-4" />
						<span>Loading address…</span>
					</div>
				) : (
					<form onSubmit={onSubmit} style={{ display: "grid", gap: 12 }}>
						<Field>
							<FieldLabel htmlFor="edit-addr-line1">Line 1</FieldLabel>
							<Input id="edit-addr-line1" value={line1} onChange={(e) => setLine1(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-addr-line2">Line 2</FieldLabel>
							<Input id="edit-addr-line2" value={line2} onChange={(e) => setLine2(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-addr-city">City</FieldLabel>
							<Input id="edit-addr-city" value={city} onChange={(e) => setCity(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-addr-province">Province</FieldLabel>
							<Input id="edit-addr-province" value={province} onChange={(e) => setProvince(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-addr-country">Country</FieldLabel>
							<Input id="edit-addr-country" value={country} onChange={(e) => setCountry(e.target.value)} />
						</Field>
						<Field>
							<FieldLabel htmlFor="edit-addr-zip">Zip</FieldLabel>
							<Input id="edit-addr-zip" value={zip} onChange={(e) => setZip(e.target.value)} />
						</Field>
						<Field>
							<div className="flex items-center gap-2">
								<input
									type="checkbox"
									id="edit-addr-is-default"
									checked={isDefault}
									onChange={(e) => setIsDefault(e.target.checked)}
									className="h-4 w-4 rounded border-gray-300"
								/>
								<FieldLabel htmlFor="edit-addr-is-default" className="cursor-pointer">Set as default address</FieldLabel>
							</div>
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
