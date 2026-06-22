import { api, setToken, clearToken } from './client';
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse, User } from './types';

export const authApi = {
	async login(req: LoginRequest): Promise<LoginResponse> {
		const res = await api.post<LoginResponse>('/auth/login', req);
		setToken(res.token);
		return res;
	},

	getMe(): Promise<User> {
		return api.get<User>('/auth/me');
	},

	register(req: RegisterRequest): Promise<RegisterResponse> {
		return api.post<RegisterResponse>('/auth/register', req);
	},

	logout(): void {
		clearToken();
	}
};
