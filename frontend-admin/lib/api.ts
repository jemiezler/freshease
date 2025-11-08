export type ApiClientOptions = {
  baseUrl?: string;
  getToken?: () => string | null;
};

const defaultBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080/api";

function defaultGetToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("admin_token");
}

export class ApiClient {
  private readonly baseUrl: string;
  private readonly getToken: () => string | null;

  constructor(options?: ApiClientOptions) {
    this.baseUrl = options?.baseUrl ?? defaultBaseUrl;
    this.getToken = options?.getToken ?? defaultGetToken;
  }

  private buildHeaders(extra?: HeadersInit): HeadersInit {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      Accept: "application/json",
      ...Object(extra),
    };
    const token = this.getToken();
    if (token) headers["Authorization"] = `Bearer ${token}`;
    return headers;
  }

  async get<T>(path: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      ...init,
      method: "GET",
      headers: this.buildHeaders(init?.headers),
      cache: "no-store",
    });
    if (!res.ok) throw new Error(`GET ${path} failed: ${res.status}`);
    return (await res.json()) as T;
  }

  async post<T, B = unknown>(path: string, body: B, init?: RequestInit): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      ...init,
      method: "POST",
      headers: this.buildHeaders(init?.headers),
      body: JSON.stringify(body),
    });
    if (!res.ok) throw new Error(`POST ${path} failed: ${res.status}`);
    return (await res.json()) as T;
  }

  async patch<T, B = unknown>(path: string, body: B, init?: RequestInit): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      ...init,
      method: "PATCH",
      headers: this.buildHeaders(init?.headers),
      body: JSON.stringify(body),
    });
    if (!res.ok) throw new Error(`PATCH ${path} failed: ${res.status}`);
    return (await res.json()) as T;
  }

  async delete<T>(path: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      ...init,
      method: "DELETE",
      headers: this.buildHeaders(init?.headers),
    });
    if (!res.ok) throw new Error(`DELETE ${path} failed: ${res.status}`);
    return (await res.json()) as T;
  }

  /**
   * Upload an image file to MinIO
   * @param file - The image file to upload
   * @param folder - The folder path in MinIO (e.g., "products", "users/avatars", "vendors/logos")
   * @returns Promise with the upload response containing object_name and url
   */
  async uploadImage(file: File, folder: string): Promise<{ object_name: string; url: string; message: string }> {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("folder", folder);

    const token = this.getToken();
    const headers: HeadersInit = {};
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
    // Don't set Content-Type header - browser will set it with boundary for multipart/form-data

    const res = await fetch(`${this.baseUrl}/uploads/images`, {
      method: "POST",
      headers,
      body: formData,
    });

    if (!res.ok) {
      const errorData = await res.json().catch(() => ({ message: "Upload failed" }));
      throw new Error(errorData.message || `Upload failed: ${res.status}`);
    }

    return (await res.json()) as { object_name: string; url: string; message: string };
  }

  /**
   * Upload an image and create/update a resource with multipart/form-data
   * @param path - API path (e.g., "/products", "/products/:id")
   * @param file - The image file to upload
   * @param data - Additional form data fields
   * @param method - HTTP method (POST or PATCH)
   * @returns Promise with the response
   */
  async postWithImage<T>(
    path: string,
    file: File | null,
    data: Record<string, unknown>,
    method: "POST" | "PATCH" = "POST"
  ): Promise<T> {
    const formData = new FormData();
    
    // Add file if provided
    if (file) {
      formData.append("image", file);
    }

    // Add other fields as JSON string in payload field, or as individual fields
    formData.append("payload", JSON.stringify(data));

    const token = this.getToken();
    const headers: HeadersInit = {};
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
    // Don't set Content-Type - browser will set it with boundary

    const res = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers,
      body: formData,
    });

    if (!res.ok) {
      const errorData = await res.json().catch(() => ({ message: "Request failed" }));
      throw new Error(errorData.message || `${method} ${path} failed: ${res.status}`);
    }

    return (await res.json()) as T;
  }
}

export const apiClient = new ApiClient();
