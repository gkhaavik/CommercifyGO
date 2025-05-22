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
  UserLoginRequest,
  UserLoginResponse,
  CreateUserRequest,
  CreateProductRequest,
  UpdateProductRequest,
  ProcessPaymentRequest,
  UpdateCartItemRequest,
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
    return this.request<ResponseDTO<CartDTO>>("/guest/cart", {
      method: "GET",
    });
  }

  async addToCart(data: AddToCartRequest): Promise<ResponseDTO<CartDTO>> {
    return this.request<ResponseDTO<CartDTO>>("/guest/cart/items", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateCartItem(
    productId: string,
    data: UpdateCartItemRequest
  ): Promise<ResponseDTO<CartDTO>> {
    return this.request<ResponseDTO<CartDTO>>(
      `/guest/cart/items/${productId}`,
      {
        method: "PUT",
        body: JSON.stringify(data),
      }
    );
  }

  async removeCartItem(productId: string): Promise<ResponseDTO<CartDTO>> {
    return this.request<ResponseDTO<CartDTO>>(
      `/guest/cart/items/${productId}`,
      {
        method: "DELETE",
      }
    );
  }

  async clearCart(): Promise<ResponseDTO<CartDTO>> {
    return this.request<ResponseDTO<CartDTO>>("/guest/cart", {
      method: "DELETE",
    });
  }

  // Order endpoints
  async createOrder(
    orderData: CreateOrderRequest
  ): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>("/guest/orders", {
      method: "POST",
      body: JSON.stringify(orderData),
    });
  }

  async getOrder(orderId: string): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>(`/orders/${orderId}`, {
      method: "GET",
    });
  }

  async getOrders(params?: {
    page?: number;
    page_size?: number;
  }): Promise<ListResponseDTO<OrderDTO>> {
    return this.request<ListResponseDTO<OrderDTO>>(
      "/orders",
      {
        method: "GET",
      },
      params
    );
  }

  async getUserOrders(params?: {
    page?: number;
    page_size?: number;
  }): Promise<ListResponseDTO<OrderDTO>> {
    return this.request<ListResponseDTO<OrderDTO>>(
      "/orders",
      {
        method: "GET",
      },
      params
    );
  }

  async processPayment(
    orderId: string,
    paymentData: ProcessPaymentRequest
  ): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>(
      `/guest/orders/${orderId}/payment`,
      {
        method: "POST",
        body: JSON.stringify(paymentData),
      }
    );
  }

  async capturePayment(paymentId: string): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>(
      `/admin/payments/${paymentId}/capture`,
      {
        method: "POST",
      }
    );
  }

  async cancelPayment(paymentId: string): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>(
      `/admin/payments/${paymentId}/cancel`,
      {
        method: "POST",
      }
    );
  }

  async refundPayment(paymentId: string): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>(
      `/admin/payments/${paymentId}/refund`,
      {
        method: "POST",
      }
    );
  }

  async forceApproveMobilePayPayment(
    paymentId: string
  ): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>(
      `/admin/payments/${paymentId}/force-approve`,
      {
        method: "POST",
      }
    );
  }
  // Product endpoints
  async getProducts(params?: {
    page?: number;
    page_size?: number;
    category_id?: number;
    currency?: string;
  }): Promise<ListResponseDTO<ProductDTO>> {
    return this.request<ListResponseDTO<ProductDTO>>("/products", {}, params);
  }

  async getProduct(
    productId: string,
    currency?: string
  ): Promise<ResponseDTO<ProductDTO>> {
    return this.request<ResponseDTO<ProductDTO>>(
      `/products/${productId}`,
      {
        method: "GET",
      },
      currency ? { currency } : undefined
    );
  }

  async searchProducts(params: {
    query?: string;
    category_id?: number;
    min_price?: number;
    max_price?: number;
    page?: number;
    page_size?: number;
  }): Promise<ListResponseDTO<ProductDTO>> {
    return this.request<ListResponseDTO<ProductDTO>>(
      "/products/search",
      {
        method: "GET",
      },
      params
    );
  }

  async createProduct(
    productData: CreateProductRequest
  ): Promise<ResponseDTO<ProductDTO>> {
    return this.request<ResponseDTO<ProductDTO>>("/products", {
      method: "POST",
      body: JSON.stringify(productData),
    });
  }

  async updateProduct(
    productId: string,
    productData: UpdateProductRequest
  ): Promise<ResponseDTO<ProductDTO>> {
    return this.request<ResponseDTO<ProductDTO>>(`/products/${productId}`, {
      method: "PUT",
      body: JSON.stringify(productData),
    });
  }

  async deleteProduct(productId: string): Promise<ResponseDTO<ProductDTO>> {
    return this.request<ResponseDTO<ProductDTO>>(`/products/${productId}`, {
      method: "DELETE",
    });
  }

  // User endpoints
  async getCurrentUser(): Promise<ResponseDTO<UserDTO>> {
    return this.request<ResponseDTO<UserDTO>>("/users/me");
  }

  async updateUser(userData: UpdateUserRequest): Promise<ResponseDTO<UserDTO>> {
    return this.request<ResponseDTO<UserDTO>>("/users/me", {
      method: "PUT",
      body: JSON.stringify(userData),
    });
  }

  async signIn(
    credentials: UserLoginRequest
  ): Promise<ResponseDTO<UserLoginResponse>> {
    return this.request<ResponseDTO<UserLoginResponse>>("/auth/signin", {
      method: "POST",
      body: JSON.stringify(credentials),
    });
  }
  async signUp(
    userData: CreateUserRequest
  ): Promise<ResponseDTO<UserLoginResponse>> {
    return this.request<ResponseDTO<UserLoginResponse>>("/auth/signup", {
      method: "POST",
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
