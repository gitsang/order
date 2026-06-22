import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { AuthState, User } from '$lib/api/types';

const TOKEN_KEY = 'auth_token';
const USER_KEY = 'auth_user';

function loadInitialState(): AuthState {
	if (!browser) {
		return { user: null, token: null, isAuthenticated: false };
	}

	try {
		const token = localStorage.getItem(TOKEN_KEY);
		const userStr = localStorage.getItem(USER_KEY);
		const user: User | null = userStr ? JSON.parse(userStr) : null;
		return {
			user,
			token,
			isAuthenticated: !!(token && user),
		};
	} catch {
		return { user: null, token: null, isAuthenticated: false };
	}
}

const authState = writable<AuthState>(loadInitialState());

export const auth = {
	subscribe: authState.subscribe,

	login(token: string, user: User) {
		if (browser) {
			localStorage.setItem(TOKEN_KEY, token);
			localStorage.setItem(USER_KEY, JSON.stringify(user));
		}
		authState.set({ user, token, isAuthenticated: true });
	},

	logout() {
		if (browser) {
			localStorage.removeItem(TOKEN_KEY);
			localStorage.removeItem(USER_KEY);
		}
		authState.set({ user: null, token: null, isAuthenticated: false });
	},

	setUser(user: User) {
		authState.update((state) => ({ ...state, user }));
	},
};
