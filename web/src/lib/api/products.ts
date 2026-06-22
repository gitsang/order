import { api } from './client';
import type { Product, ListProductsRequest } from './types';

export const productsApi = {
	list(params?: ListProductsRequest): Promise<Product[]> {
		return api.get<Product[]>('/products', params as Record<string, string | undefined>);
	},

	get(id: string): Promise<Product> {
		return api.get<Product>(`/products/${id}`);
	}
};
