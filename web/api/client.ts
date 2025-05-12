import {
  CartDTO,
  OrderDTO,
  ProductDTO,
  UserDTO,
  CreateOrderRequest,
  AddToCartRequest,
  UpdateUserRequest,
  ResponseDTO,
  ListResponseDTO,
  ProductSearchRequest,
  PaginationDTO,
} from "../types/api";

export class CommercifyClient {
  private baseUrl: string;
  private token?: string;

  constructor(baseUrl: string, token?: string) {
    this.baseUrl = baseUrl;
    this.token = token;
  }

  private buildUrl(endpoint: string, params?: Record<string, any>): string {
    const url = new URL(`${this.baseUrl}${endpoint}`);
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          if (Array.isArray(value)) {
            value.forEach((v) => url.searchParams.append(key, String(v)));
          } else {
            url.searchParams.append(key, String(value));
          }
        }
      });
    }
    return url.toString();
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
    params?: Record<string, any>
  ): Promise<T> {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...(this.token && { Authorization: `Bearer ${this.token}` }),
      ...options.headers,
    };

    const url = this.buildUrl(endpoint, params);
    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      throw new Error(`API request failed: ${response.statusText}`);
    }

    return response.json();
  }

  // Cart endpoints
  async getCart(): Promise<ResponseDTO<CartDTO>> {
    return this.request<ResponseDTO<CartDTO>>("/api/cart");
  }

  async addToCart(data: AddToCartRequest): Promise<ResponseDTO<CartDTO>> {
    return this.request<ResponseDTO<CartDTO>>("/api/cart/items", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  // Order endpoints
  async createOrder(
    orderData: CreateOrderRequest
  ): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>("/api/orders", {
      method: "POST",
      body: JSON.stringify(orderData),
    });
  }

  async getOrder(orderId: string): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>(`/api/orders/${orderId}`);
  }

  // Product endpoints
  async getProducts(params?: {
    page?: number;
    page_size?: number;
    category_id?: number;
    currency?: string;
  }): Promise<ListResponseDTO<ProductDTO>> {
    return this.request<ListResponseDTO<ProductDTO>>(
      "/api/products",
      {},
      params
    );
  }

  async getProduct(
    productId: string,
    currency?: string
  ): Promise<ResponseDTO<ProductDTO>> {
    return this.request<ResponseDTO<ProductDTO>>(
      `/api/products/${productId}`,
      {},
      currency ? { currency } : undefined
    );
  }

  async searchProducts(
    params: ProductSearchRequest
  ): Promise<ListResponseDTO<ProductDTO>> {
    return this.request<ListResponseDTO<ProductDTO>>("/api/products/search", {
      method: "POST",
      body: JSON.stringify(params),
    });
  }

  // User endpoints
  async getCurrentUser(): Promise<ResponseDTO<UserDTO>> {
    return this.request<ResponseDTO<UserDTO>>("/api/user/me");
  }

  async updateUser(userData: UpdateUserRequest): Promise<ResponseDTO<UserDTO>> {
    return this.request<ResponseDTO<UserDTO>>("/api/user", {
      method: "PUT",
      body: JSON.stringify(userData),
    });
  }
}

// Example usage:
// const client = new CommercifyClient('https://api.commercify.com', 'your-auth-token');
//
// // Get products with pagination and filters
// const products = await client.getProducts({
//   page: 1,
//   page_size: 20,
//   category_id: 123,
//   currency: 'USD'
// });
//
// // Search products with advanced filters
// const searchResults = await client.searchProducts({
//   query: 'gaming laptop',
//   category_id: 123,
//   min_price: 500,
//   max_price: 2000,
//   page: 1,
//   page_size: 20
// });
