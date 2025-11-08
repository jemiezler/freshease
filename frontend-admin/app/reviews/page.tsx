"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, TrashIcon, PlusIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import DataTable from "./_components/reviews-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateReviewDialog } from "./_components/create-review-dialog";
import { EditReviewDialog } from "./_components/edit-review-dialog";
import type { Review, ReviewPayload } from "@/types/review";

const reviews = createResource<Review, ReviewPayload, ReviewPayload>({
	basePath: "/reviews",
});

export default function ReviewsPage() {
	const [items, setItems] = useState<Review[]>([]);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [createOpen, setCreateOpen] = useState(false);
	const [editId, setEditId] = useState<string | null>(null);

	const load = useCallback(async () => {
		setLoading(true);
		setError(null);
		try {
			const res = await reviews.list();
			setItems(res.data ?? []);
		} catch (e) {
			setError(e instanceof Error ? e.message : "Failed to load");
		} finally {
			setLoading(false);
		}
	}, []);

	useEffect(() => {
		void load();
	}, [load]);

	const onDelete = useCallback(
		async (id: string) => {
			if (!confirm("Delete this review?")) return;
			try {
				await reviews.delete(id);
				await load();
			} catch (e) {
				alert(e instanceof Error ? e.message : "Delete failed");
			}
		},
		[load]
	);

	const columns = useMemo<ColumnDef<Review>[]>(
		() => [
			{
				accessorKey: "rating",
				header: "Rating",
				cell: ({ row }) => {
					const rating = row.getValue("rating") as number;
					return (
						<div className="flex items-center gap-1">
							{Array.from({ length: 5 }).map((_, i) => (
								<span key={i} className={i < rating ? "text-yellow-400" : "text-zinc-300"}>
									★
								</span>
							))}
							<span className="ml-2 text-sm text-zinc-600">({rating})</span>
						</div>
					);
				},
			},
			{
				accessorKey: "comment",
				header: "Comment",
				cell: ({ row }) => {
					const comment = row.getValue("comment") as string;
					return comment ? (comment.length > 50 ? comment.slice(0, 50) + "..." : comment) : "-";
				},
			},
			{
				accessorKey: "product_id",
				header: "Product ID",
				cell: ({ row }) => {
					const productId = row.getValue("product_id") as string;
					return <span className="font-mono text-xs">{productId.slice(0, 8)}...</span>;
				},
			},
			{
				accessorKey: "created_at",
				header: "Created At",
				cell: ({ row }) => {
					const date = row.getValue("created_at") as string;
					return date ? new Date(date).toLocaleDateString() : "-";
				},
			},
			{
				id: "actions",
				header: "Actions",
				cell: ({ row }) => {
					const review = row.original;
					return (
						<div className="flex gap-2">
							<Button size="icon" variant="ghost" onClick={() => setEditId(review.id)}>
								<PencilIcon className="size-4" />
							</Button>
							<Button size="icon" variant="ghost" onClick={() => onDelete(review.id)}>
								<TrashIcon className="size-4 text-red-500" />
							</Button>
						</div>
					);
				},
			},
		],
		[onDelete]
	);

	return (
		<div>
			<div className="mb-6 flex items-center justify-between">
				<h1 className="text-3xl font-bold text-zinc-900">Reviews</h1>
				<Button onClick={() => setCreateOpen(true)}>
					<PlusIcon className="size-4 mr-2" />
					New Review
				</Button>
			</div>
			{error && <p className="mb-4 text-red-500">{error}</p>}
			<div className="min-h-[200px]">
				{loading ? (
					<div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
						<Spinner className="size-6" />
						<span>Loading reviews…</span>
					</div>
				) : (
					<DataTable columns={columns} data={items} />
				)}
			</div>
			<CreateReviewDialog
				open={createOpen}
				onOpenChange={setCreateOpen}
				onSaved={async () => {
					setCreateOpen(false);
					await load();
				}}
			/>
			{editId && (
				<EditReviewDialog
					id={editId}
					onOpenChange={(open) => {
						if (!open) setEditId(null);
					}}
					onSaved={async () => {
						setEditId(null);
						await load();
					}}
				/>
			)}
		</div>
	);
}

