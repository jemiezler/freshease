"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { createResource } from "@/lib/resource";
import { Button } from "@/components/ui/button";
import { PencilIcon, PlusIcon, TrashIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import DataTable from "./_components/products-table";
import { ColumnDef } from "@tanstack/react-table";
import { CreateProductDialog } from "./_components/create-product-dialog";
import { EditProductDialog } from "./_components/edit-product-dialog";
import { CreateCategoryDialog } from "./_components/create-category-dialog";
import { EditCategoryDialog } from "./_components/edit-catagory-dialog";
import type { Product, ProductPayload } from "@/types/product";
import type { Category, CategoryPayload } from "@/types/catagory";

const categories = createResource<Category, CategoryPayload, CategoryPayload>({
  basePath: "/product_categories",
});

const products = createResource<Product, ProductPayload, ProductPayload>({
  basePath: "/products",
});

export default function ProductsPage() {
  const [categoryItems, setCategoryItems] = useState<Category[]>([]);
  const [categoryLoading, setCategoryLoading] = useState(false);
  const [categoryError, setCategoryError] = useState<string | null>(null);
  const [createCategoryOpen, setCreateCategoryOpen] = useState(false);
  const [editCategoryId, setEditCategoryId] = useState<string | null>(null);

  const [items, setItems] = useState<Product[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [createOpen, setCreateOpen] = useState(false);
  const [editId, setEditId] = useState<string | null>(null);

  const loadCategories = useCallback(async () => {
    setCategoryLoading(true);
    setCategoryError(null);
    try {
      const res = await categories.list();
      setCategoryItems(res.data ?? []);
    } catch (e) {
      setCategoryError(e instanceof Error ? e.message : "Failed to load");
    } finally {
      setCategoryLoading(false);
    }
  }, []);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await products.list();
      setItems(res.data ?? []);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadCategories();
    void load();
  }, [loadCategories, load]);

  // Group products by category_id
  const productsByCategory = useMemo(() => {
    const grouped: Record<string, Product[]> = {};
    items.forEach((product) => {
      const categoryId = product.category_id || "uncategorized";
      if (!grouped[categoryId]) {
        grouped[categoryId] = [];
      }
      grouped[categoryId].push(product);
    });
    return grouped;
  }, [items]);

  const reloadAll = useCallback(async () => {
    await Promise.all([loadCategories(), load()]);
  }, [loadCategories, load]);

  const onDeleteCategory = useCallback(
    async (id: string) => {
      if (!confirm("Delete this category?")) return;
      try {
        await categories.delete(id);
        await reloadAll();
      } catch (e) {
        alert(e instanceof Error ? e.message : "Delete failed");
      }
    },
    [reloadAll]
  );

  const onDelete = useCallback(
    async (id: string) => {
      if (!confirm("Delete this product?")) return;
      try {
        await products.delete(id);
        await reloadAll();
      } catch (e) {
        alert(e instanceof Error ? e.message : "Delete failed");
      }
    },
    [reloadAll]
  );

  const productColumns = useMemo<ColumnDef<Product>[]>(
    () => [
      {
        accessorKey: "name",
        header: "Name",
        cell: ({ row }) => row.getValue("name") ?? "-",
      },
      {
        id: "price",
        header: "Price",
        cell: ({ row }) => {
          const product = row.original;
          return product.price != null ? product.price.toString() : "-";
        },
      },
      {
        id: "actions",
        header: "Actions",
        cell: ({ row }) => {
          const product = row.original;
          return (
            <div className="flex gap-2">
              <Button
                size="icon"
                variant="ghost"
                onClick={() => setEditId(product.id)}
              >
                <PencilIcon className="size-4" />
              </Button>
              <Button
                size="icon"
                variant="ghost"
                onClick={() => onDelete(product.id)}
              >
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
    <div className="space-y-8">
      <div>
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            marginBottom: 12,
          }}
        >
          <h1 style={{ fontSize: 20, fontWeight: 600 }}>Products</h1>
          <div className="flex gap-2">
            <Button onClick={() => setCreateCategoryOpen(true)}>
              <PlusIcon className="size-4" />
              New Category
            </Button>
          </div>
        </div>
        {(categoryError || error) && (
          <p style={{ color: "red" }}>{categoryError || error}</p>
        )}
        <div className="min-h-[200px]">
          {categoryLoading || loading ? (
            <div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
              <Spinner className="size-6" />
              <span>Loading…</span>
            </div>
          ) : categoryItems.length === 0 ? (
            <div className="flex h-full flex-col items-center justify-center gap-2 text-sm text-muted-foreground">
              <span>No categories found.</span>
            </div>
          ) : (
            <div className="overflow-hidden rounded-md border bg-white px-4">
              <Accordion type="single" collapsible className="w-full">
                {categoryItems.map((category) => {
                  const categoryProducts =
                    productsByCategory[category.id] || [];
                  return (
                    <AccordionItem key={category.id} value={category.id}>
                      <AccordionTrigger>
                        <div className="flex items-center justify-between w-full pr-4">
                          <div className="flex items-center gap-3">
                            <span className="font-medium">{category.name}</span>
                            <span className="text-sm text-muted-foreground">
                              ({categoryProducts.length} product
                              {categoryProducts.length !== 1 ? "s" : ""})
                            </span>
                          </div>
                        </div>
                      </AccordionTrigger>
                      <AccordionContent>
                        <div className="space-y-4">
                          {category.description && (
                            <div>
                              <p className="text-sm text-muted-foreground mb-2">
                                <strong>Description:</strong>{" "}
                                {category.description}
                              </p>
                            </div>
                          )}
                          <div className="flex mb-4 w-full justify-between">
                            <div className="flex gap-2">
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => setEditCategoryId(category.id)}
                              >
                                <PencilIcon className="size-4 mr-2" />
                                Edit Category
                              </Button>
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => onDeleteCategory(category.id)}
                              >
                                <TrashIcon className="size-4 mr-2 text-red-500" />
                                Delete Category
                              </Button>
                            </div>
                            <div>
                              <Button
                                size="sm"
                                variant="outline"
                                onClick={() => setCreateOpen(true)}
                              >
                                <PlusIcon className="size-4 mr-2" />
                                Add Product
                              </Button>
                            </div>
                          </div>
                          {categoryProducts.length > 0 ? (
                            <DataTable
                              columns={productColumns}
                              data={categoryProducts}
                            />
                          ) : (
                            <div className="text-sm text-muted-foreground py-4 text-center">
                              No products in this category.
                            </div>
                          )}
                        </div>
                      </AccordionContent>
                    </AccordionItem>
                  );
                })}
                {/* Show uncategorized products if any */}
                {productsByCategory["uncategorized"] &&
                  productsByCategory["uncategorized"].length > 0 && (
                    <AccordionItem key="uncategorized" value="uncategorized">
                      <AccordionTrigger>
                        <div className="flex items-center justify-between w-full pr-4">
                          <div className="flex items-center gap-3">
                            <span className="font-medium">Uncategorized</span>
                            <span className="text-sm text-muted-foreground">
                              ({productsByCategory["uncategorized"].length}{" "}
                              product
                              {productsByCategory["uncategorized"].length !== 1
                                ? "s"
                                : ""}
                              )
                            </span>
                          </div>
                        </div>
                      </AccordionTrigger>
                      <AccordionContent>
                        <DataTable
                          columns={productColumns}
                          data={productsByCategory["uncategorized"]}
                        />
                      </AccordionContent>
                    </AccordionItem>
                  )}
              </Accordion>
            </div>
          )}
        </div>
      </div>
      <CreateProductDialog
        open={createOpen}
        onOpenChange={setCreateOpen}
        onSaved={async () => {
          setCreateOpen(false);
          await reloadAll();
        }}
      />
      {editId && (
        <EditProductDialog
          id={editId}
          onOpenChange={(open) => {
            if (!open) setEditId(null);
          }}
          onSaved={async () => {
            setEditId(null);
            await reloadAll();
          }}
        />
      )}

      {/* Category Dialogs */}
      <CreateCategoryDialog
        open={createCategoryOpen}
        onOpenChange={setCreateCategoryOpen}
        onSaved={async () => {
          setCreateCategoryOpen(false);
          await reloadAll();
        }}
      />
      {editCategoryId && (
        <EditCategoryDialog
          id={editCategoryId}
          onOpenChange={(open) => {
            if (!open) setEditCategoryId(null);
          }}
          onSaved={async () => {
            setEditCategoryId(null);
            await reloadAll();
          }}
        />
      )}
    </div>
  );
}
