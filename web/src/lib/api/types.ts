export interface ApiResponse<T = unknown> {
	code: number;
	message: string;
	data?: T;
}

export class ApiError extends Error {
	constructor(
		public readonly status: number,
		message: string,
		public readonly code: number
	) {
		super(message);
		this.name = 'ApiError';
	}
}

export interface User {
	id: string;
	username: string;
	name: string;
	phone: string;
	role: string;
	created_at?: string;
	updated_at?: string;
}

export interface Category {
	id: string;
	name: string;
	sort_order: number;
}

export interface Product {
	id: string;
	category_id: string;
	name: string;
	description: string;
	price: number;
	image: string;
	status: string;
	sort_order: number;
	created_at?: string;
	updated_at?: string;
	category?: Category;
}

export interface OrderItem {
	id: string;
	order_id: string;
	product_id: string;
	quantity: number;
	price: number;
	product?: Product;
}

export interface Order {
	id: string;
	user_id: string;
	order_no: string;
	total_amount: number;
	status: string;
	remark: string;
	items?: OrderItem[];
	created_at?: string;
	updated_at?: string;
}

export interface LoginRequest {
	username: string;
	password: string;
}

export interface LoginResponse {
	token: string;
}

export interface RegisterRequest {
	username: string;
	password: string;
	name: string;
	phone: string;
}

export interface RegisterResponse {
	id: string;
	username: string;
}

export interface ListProductsRequest {
	category_id?: string;
}

export interface CreateOrderItemRequest {
	product_id: string;
	quantity: number;
}

export interface CreateOrderRequest {
	items: CreateOrderItemRequest[];
	remark?: string;
}

export interface ListOrdersRequest {
	limit?: number;
	offset?: number;
}

// Cart item (client-side, extends product with quantity)
export interface CartItem {
	product: Product;
	quantity: number;
}

// Auth state
export interface AuthState {
	user: User | null;
	token: string | null;
	isAuthenticated: boolean;
}

// Cart state
export interface CartState {
	items: CartItem[];
	totalItems: number;
	totalAmount: number;
}
