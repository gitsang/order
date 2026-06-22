import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import type { ApiResponse } from './types';
import { ApiError } from './types';

const TOKEN_KEY = 'auth_token';
const BASE_URL = '/api/v1';

function getToken(): string | null {
	if (!browser) return null;
	return localStorage.getItem(TOKEN_KEY);
}

function setToken(token: string): void {
	if (!browser) return;
	localStorage.setItem(TOKEN_KEY, token);
}

function clearToken(): void {
	if (!browser) return;
	localStorage.removeItem(TOKEN_KEY);
}

async function handleUnauthorized(): Promise<never> {
	clearToken();
	if (browser) {
		await goto('/login');
	}
	throw new ApiError(401, 'Unauthorized', 401);
}

interface RequestOptions extends Omit<RequestInit, 'body'> {
	body?: unknown;
}

async function request<T>(
	endpoint: string,
	options: RequestOptions = {}
): Promise<T> {
	const { body, headers: customHeaders, ...rest } = options;

	const headers = new Headers(customHeaders);
	headers.set('Accept', 'application/json');

	const token = getToken();
	if (token) {
		headers.set('Authorization', `Bearer ${token}`);
	}

	if (body !== undefined) {
		headers.set('Content-Type', 'application/json');
	}

	const response = await fetch(`${BASE_URL}${endpoint}`, {
		...rest,
		headers,
		body: body !== undefined ? JSON.stringify(body) : undefined
	});

	if (response.status === 401) {
		return handleUnauthorized();
	}

	const json: ApiResponse<T> = await response.json();

	if (response.ok && json.code === 0) {
		return json.data as T;
	}

	throw new ApiError(response.status, json.message, json.code);
}

export const api = {
	get<T>(endpoint: string, params?: Record<string, string | undefined>): Promise<T> {
		let url = endpoint;
		if (params) {
			const searchParams = new URLSearchParams();
			for (const [key, value] of Object.entries(params)) {
				if (value !== undefined) {
					searchParams.set(key, value);
				}
			}
			const qs = searchParams.toString();
			if (qs) url += `?${qs}`;
		}
		return request<T>(url, { method: 'GET' });
	},

	post<T>(endpoint: string, body?: unknown): Promise<T> {
		return request<T>(endpoint, { method: 'POST', body });
	},

	put<T>(endpoint: string, body?: unknown): Promise<T> {
		return request<T>(endpoint, { method: 'PUT', body });
	},

	delete<T>(endpoint: string): Promise<T> {
		return request<T>(endpoint, { method: 'DELETE' });
	}
};

export { getToken, setToken, clearToken };
