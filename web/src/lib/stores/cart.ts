import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { CartItem, CartState, Product } from '$lib/api/types';

const CART_KEY = 'coffee_cart';

function calculateTotals(items: CartItem[]): Pick<CartState, 'totalItems' | 'totalAmount'> {
	const totalItems = items.reduce((sum, item) => sum + item.quantity, 0);
	const totalAmount = items.reduce((sum, item) => sum + item.product.price * item.quantity, 0);
	return { totalItems, totalAmount };
}

function loadCart(): CartState {
	if (!browser) {
		return { items: [], totalItems: 0, totalAmount: 0 };
	}

	try {
		const cartStr = localStorage.getItem(CART_KEY);
		if (!cartStr) {
			return { items: [], totalItems: 0, totalAmount: 0 };
		}
		const items: CartItem[] = JSON.parse(cartStr);
		return { items, ...calculateTotals(items) };
	} catch {
		return { items: [], totalItems: 0, totalAmount: 0 };
	}
}

function persistCart(items: CartItem[]) {
	if (browser) {
		localStorage.setItem(CART_KEY, JSON.stringify(items));
	}
}

const cartState = writable<CartState>(loadCart());

export const cart = {
	subscribe: cartState.subscribe,

	add(product: Product, quantity: number = 1) {
		cartState.update((state) => {
			const existing = state.items.findIndex((item) => item.product.id === product.id);
			let items: CartItem[];
			if (existing >= 0) {
				items = state.items.map((item, i) =>
					i === existing ? { ...item, quantity: item.quantity + quantity } : item,
				);
			} else {
				items = [...state.items, { product, quantity }];
			}
			persistCart(items);
			return { items, ...calculateTotals(items) };
		});
	},

	remove(productId: string) {
		cartState.update((state) => {
			const items = state.items.filter((item) => item.product.id !== productId);
			persistCart(items);
			return { items, ...calculateTotals(items) };
		});
	},

	updateQuantity(productId: string, quantity: number) {
		if (quantity <= 0) {
			cart.remove(productId);
			return;
		}
		cartState.update((state) => {
			const items = state.items.map((item) =>
				item.product.id === productId ? { ...item, quantity } : item,
			);
			persistCart(items);
			return { items, ...calculateTotals(items) };
		});
	},

	clear() {
		cartState.update(() => {
			const items: CartItem[] = [];
			persistCart(items);
			return { items, totalItems: 0, totalAmount: 0 };
		});
	},
};
