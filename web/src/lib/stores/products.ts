import { writable } from 'svelte/store';
import type { Product } from '$lib/api/types';

interface ProductsState {
	products: Product[];
	loading: boolean;
	error: string | null;
	lastFetched: number | null;
}

const initialState: ProductsState = {
	products: [],
	loading: false,
	error: null,
	lastFetched: null,
};

const productsState = writable<ProductsState>(initialState);

export const products = {
	subscribe: productsState.subscribe,

	setProducts(products: Product[]) {
		productsState.set({
			products,
			loading: false,
			error: null,
			lastFetched: Date.now(),
		});
	},

	setLoading() {
		productsState.update((state) => ({ ...state, loading: true, error: null }));
	},

	setError(error: string) {
		productsState.update((state) => ({ ...state, loading: false, error }));
	},

	getProductById(id: string): Product | undefined {
		let result: Product | undefined;
		productsState.subscribe((state) => {
			result = state.products.find((p) => p.id === id);
		})();
		return result;
	},
};
