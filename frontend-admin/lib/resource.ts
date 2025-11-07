import { apiClient } from "./api";

export type UpdateMethod = "PUT" | "PATCH";

export type ResourceConfig = {
  basePath: string; // e.g. "/products" or "/cart_items"
  updateMethod?: UpdateMethod; // default PATCH
};

export function createResource<TItem, TCreate = Partial<TItem>, TUpdate = Partial<TItem>>(config: ResourceConfig) {
  const updateMethod = config.updateMethod ?? "PATCH";
  return {
    list: (params?: Record<string, string | number | boolean | undefined>) => {
      const q = params
        ? "?" + new URLSearchParams(Object.entries(params).reduce<Record<string,string>>((a,[k,v])=>{ if (v!==undefined) a[k]=String(v); return a;},{})).toString()
        : "";
      return apiClient.get<{ data?: TItem[] }>(`${config.basePath}${q}`);
    },
    get: (id: string) => apiClient.get<{ data?: TItem }>(`${config.basePath}/${id}`),
    create: (payload: TCreate) => apiClient.post<{ data?: TItem }>(`${config.basePath}`, payload as unknown as TCreate),
    update: (id: string, payload: TUpdate) =>
      updateMethod === "PUT"
        ? apiClient.post<{ data?: TItem }>(`${config.basePath}/${id}?_method=PUT`, payload as unknown as TUpdate)
        : apiClient.patch<{ data?: TItem }>(`${config.basePath}/${id}`, payload as unknown as TUpdate),
    delete: (id: string) => apiClient.delete<{ data?: unknown }>(`${config.basePath}/${id}`),
  };
}
