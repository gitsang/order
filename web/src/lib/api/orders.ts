import { api } from './client';
import type { Order, CreateOrderRequest, ListOrdersRequest } from './types';

export const ordersApi = {
	create(req: CreateOrderRequest): Promise<Order> {
		return api.post<Order>('/orders', req);
	},

	list(params?: ListOrdersRequest): Promise<Order[]> {
		const query: Record<string, string | undefined> = {};
		if (params?.limit !== undefined) query.limit = String(params.limit);
		if (params?.offset !== undefined) query.offset = String(params.offset);
		return api.get<Order[]>('/orders', query);
	},

	get(id: string): Promise<Order> {
		return api.get<Order>(`/orders/${id}`);
	}
};
