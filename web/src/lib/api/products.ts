import { api } from './client';
import type { Product, ListProductsRequest, CreateProductRequest, UpdateProductRequest } from './types';

export const productsApi = {
	list(params?: ListProductsRequest): Promise<Product[]> {
		return api.get<Product[]>('/products', params as Record<string, string | undefined>);
	},

	get(id: string): Promise<Product> {
		return api.get<Product>(`/products/${id}`);
	},

	create(data: CreateProductRequest): Promise<Product> {
		return api.post<Product>('/products', data);
	},

	update(id: string, data: UpdateProductRequest): Promise<Product> {
		return api.put<Product>(`/products/${id}`, data);
	},

	delete(id: string): Promise<void> {
		return api.delete<void>(`/products/${id}`);
	}
};
