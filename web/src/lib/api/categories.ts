import { api } from './client';
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from './types';

export const categoriesApi = {
	list(): Promise<Category[]> {
		return api.get<Category[]>('/categories');
	},

	get(id: string): Promise<Category> {
		return api.get<Category>(`/categories/${id}`);
	},

	create(data: CreateCategoryRequest): Promise<Category> {
		return api.post<Category>('/categories', data);
	},

	update(id: string, data: UpdateCategoryRequest): Promise<Category> {
		return api.put<Category>(`/categories/${id}`, data);
	},

	delete(id: string): Promise<void> {
		return api.delete<void>(`/categories/${id}`);
	}
};
