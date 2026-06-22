import { api, setToken, clearToken } from './client';
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from './types';

export const authApi = {
	async login(req: LoginRequest): Promise<LoginResponse> {
		const res = await api.post<LoginResponse>('/auth/login', req);
		setToken(res.token);
		return res;
	},

	register(req: RegisterRequest): Promise<RegisterResponse> {
		return api.post<RegisterResponse>('/auth/register', req);
	},

	logout(): void {
		clearToken();
	}
};
