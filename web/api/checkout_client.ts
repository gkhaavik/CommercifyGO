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
  ApplyDiscountRequest
} from "../types/checkout";
import { ListResponseDTO, OrderDTO, ResponseDTO } from "../types/api";

// Add this implementation to the CommercifyClient class in client.ts

  // Guest Checkout API
  async createGuestCheckout(data: CreateGuestCheckoutRequest): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>("/api/guest/checkout", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async getGuestCheckout(sessionId: string): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(`/api/guest/checkout/${sessionId}`, {
      method: "GET",
    });
  }

  async updateCheckoutItem(
    productId: string,
    data: UpdateCheckoutItemRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      `/guest/checkout/items/${productId}`,
      {
        method: "PUT",
        body: JSON.stringify(data),
      }
    );
  }

  async removeCheckoutItem(productId: string): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      `/guest/checkout/items/${productId}`,
      {
        method: "DELETE",
      }
    );
  }

  async clearCheckout(): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>("/guest/checkout", {
      method: "DELETE",
    });
  }

  async setShippingAddress(
    data: SetShippingAddressRequest
  ): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      "/guest/checkout/shipping-address",
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
      "/guest/checkout/billing-address",
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
      "/guest/checkout/customer-details",
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
      "/guest/checkout/shipping-method",
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
      "/guest/checkout/discount",
      {
        method: "POST",
        body: JSON.stringify(data),
      }
    );
  }

  async removeCheckoutDiscount(): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      "/guest/checkout/discount",
      {
        method: "DELETE",
      }
    );
  }

  async convertCheckoutToOrder(): Promise<ResponseDTO<OrderDTO>> {
    return this.request<ResponseDTO<OrderDTO>>(
      "/guest/checkout/to-order",
      {
        method: "POST",
      }
    );
  }

  async convertGuestCheckoutToUserCheckout(): Promise<ResponseDTO<CheckoutDTO>> {
    return this.request<ResponseDTO<CheckoutDTO>>(
      "/checkout/convert",
      {
        method: "POST",
      }
    );
  }

  // Admin checkout endpoints
  async getAllCheckouts(params?: {
    page?: number;
    page_size?: number;
    status?: string;
  }): Promise<ListResponseDTO<CheckoutDTO>> {
    return this.request<ListResponseDTO<CheckoutDTO>>(
      "/admin/checkouts",
      {
        method: "GET",
      },
      params
    );
  }
