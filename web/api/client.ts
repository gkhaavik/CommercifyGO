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
import {
  CheckoutDTO,
  CreateGuestCheckoutRequest,
  CreateCheckoutRequest,
  AddCheckoutItemRequest,
  UpdateCheckoutItemRequest,
  SetShippingAddressRequest,
  SetBillingAddressRequest,
  SetCustomerDetailsRequest,
  SetShippingMethodRequest,
  ApplyDiscountRequest,
} from "../types/checkout";

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

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => null);

        // Create a more detailed error message
        const errorMessage =
          errorData?.error?.message || response.statusText || "Unknown error";
        const error = new Error(`API request failed: ${errorMessage}`);

        // Attach additional properties for error handling
        (error as any).status = response.status;
        (error as any).statusText = response.statusText;
        (error as any).errorData = errorData;

        throw error;
      }

      const data = await response.json();
      return data;
    } catch (error) {
      // If the error is already formatted by our code above, just rethrow it
      if ((error as any).status) {
        throw error;
      }

      // Otherwise, it's likely a network error or other issue
      console.error("API Request Error:", error);
      throw new Error(
        `API request failed: ${
          error instanceof Error ? error.message : "Network error"
        }`
      );
    }
  }

  // DEPRECATED: Cart endpoints - Use checkout endpoints instead
  /**
   * @deprecated The cart API is deprecated. Use getCheckout() instead.
   */
  async getCart(): Promise<ResponseDTO<CartDTO>> {
    console.warn('DEPRECATED: getCart() is deprecated. Use getCheckout() instead.');
    return this.request<ResponseDTO<CartDTO>>("/guest/cart", {
      method: "GET",
    });
  }

  /**
   * @deprecated The cart API is deprecated. Use addToCheckout() instead.
   */
  async addToCart(data: AddToCartRequest): Promise<ResponseDTO<CartDTO>> {
    console.warn('DEPRECATED: addToCart() is deprecated. Use addToCheckout() instead.');
    return this.request<ResponseDTO<CartDTO>>("/guest/cart/items", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  /**
   * @deprecated The cart API is deprecated. Use updateCheckoutItem() instead.
   */
  async updateCartItem(
    productId: string,
    data: UpdateCartItemRequest
  ): Promise<ResponseDTO<CartDTO>> {
    console.warn('DEPRECATED: updateCartItem() is deprecated. Use updateCheckoutItem() instead.');
    return this.request<ResponseDTO<CartDTO>>(
      `/guest/cart/items/${productId}`,
      {
        method: "PUT",
        body: JSON.stringify(data),
      }
    );
  }

  /**
   * @deprecated The cart API is deprecated. Use removeFromCheckout() instead.
   */
  async removeCartItem(productId: string): Promise<ResponseDTO<CartDTO>> {
    console.warn('DEPRECATED: removeCartItem() is deprecated. Use removeFromCheckout() instead.');
    return this.request<ResponseDTO<CartDTO>>(
      `/guest/cart/items/${productId}`,
      {
        method: "DELETE",
      }
    );
  }

  /**
   * @deprecated The cart API is deprecated. Use clearCheckout() instead.
   */
  async clearCart(): Promise<ResponseDTO<CartDTO>> {
    console.warn('DEPRECATED: clearCart() is deprecated. Use clearCheckout() instead.');
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

  // Guest Checkout API
  async createGuestCheckout(
    data: CreateGuestCheckoutRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>("/api/guest/checkout", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getGuestCheckout(): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>("/api/guest/checkout", {
      method: "GET",
    });
  }

  async addCheckoutItem(
    data: AddCheckoutItemRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>("/api/guest/checkout/items", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateCheckoutItem(
    productId: number,
    data: UpdateCheckoutItemRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      `/api/guest/checkout/items/${productId}`,
      {
        method: "PUT",
        body: JSON.stringify(data),
      }
    );
  }

  async removeCheckoutItem(
    productId: number
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      `/api/guest/checkout/items/${productId}`,
      {
        method: "DELETE",
      }
    );
  }

  async clearCheckout(): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>("/api/guest/checkout", {
      method: "DELETE",
    });
  }

  async setShippingAddress(
    data: SetShippingAddressRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      "/api/guest/checkout/shipping-address",
      {
        method: "PUT",
        body: JSON.stringify(data),
      }
    );
  }

  async setBillingAddress(
    data: SetBillingAddressRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      "/api/guest/checkout/billing-address",
      {
        method: "PUT",
        body: JSON.stringify(data),
      }
    );
  }

  async setCustomerDetails(
    data: SetCustomerDetailsRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      "/api/guest/checkout/customer-details",
      {
        method: "PUT",
        body: JSON.stringify(data),
      }
    );
  }

  async setShippingMethod(
    data: SetShippingMethodRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      "/api/guest/checkout/shipping-method",
      {
        method: "PUT",
        body: JSON.stringify(data),
      }
    );
  }

  async applyCheckoutDiscount(
    data: ApplyDiscountRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      "/api/guest/checkout/discount",
      {
        method: "POST",
        body: JSON.stringify(data),
      }
    );
  }

  async removeCheckoutDiscount(): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      "/api/guest/checkout/discount",
      {
        method: "DELETE",
      }
    );
  }

  async convertCheckoutToOrder(): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>("/api/guest/checkout/to-order", {
      method: "POST",
    });
  }

  // Authenticated Checkout API
  async getUserCheckout(): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>("/api/checkout", {
      method: "GET",
    });
  }

  async createCheckout(
    data: CreateCheckoutRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>("/api/checkout", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async convertGuestCheckoutToUserCheckout(): Promise<
    ResponseDTO<CheckoutDTO>
  > {
    return this.request<ResponseDTO<CheckoutDTO>>("/api/checkout/convert", {
      method: "POST",
    });
  }

  // Admin checkout endpoints
  async getAdminCheckouts(params?: {
    page?: number;
    page_size?: number;
    status?: string;
  }): Promise<ListResponseDTO<CheckoutDTO>> {
    return this.request<ListResponseDTO<CheckoutDTO>>(
      "/api/admin/checkouts",
      {},
      params
    );
  }

  async getAdminCheckoutById(
    checkoutId: number
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      `/api/admin/checkouts/${checkoutId}`,
      {
        method: "GET",
      }
    );
  }

  async deleteAdminCheckout(checkoutId: number): Promise<ResponseDTO<string>> {
    return this.request<ResponseDTO<string>>(
      `/api/admin/checkouts/${checkoutId}`,
      {
        method: "DELETE",
      }
    );
  }

  async getCheckoutsByUser(
    userId: number,
    params?: {
      page?: number;
      page_size?: number;
      status?: string;
    }
  ): Promise<ListResponseDTO<CheckoutDTO>> {
    return this.request<ListResponseDTO<CheckoutDTO>>(
      `/api/admin/users/${userId}/checkouts`,
      {},
      params
    );
  }

  async getAbandonedCheckouts(): Promise<ListResponseDTO<CheckoutDTO>> {
    return this.request<ListResponseDTO<CheckoutDTO>>(
      `/api/admin/checkouts/abandoned`,
      {
        method: "GET",
      }
    );
  }

  async getExpiredCheckouts(): Promise<ListResponseDTO<CheckoutDTO>> {
    return this.request<ListResponseDTO<CheckoutDTO>>(
      `/api/admin/checkouts/expired`,
      {
        method: "GET",
      }
    );
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
