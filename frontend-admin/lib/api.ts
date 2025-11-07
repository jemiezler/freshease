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
}

export const apiClient = new ApiClient();
