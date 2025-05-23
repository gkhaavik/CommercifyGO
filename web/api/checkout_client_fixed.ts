import { 
  CheckoutDTO} from "../types/checkout";
import { ListResponseDTO, OrderDTO, ResponseDTO } from "../types/api";

// These are extension methods for the CommercifyClient class
// Add these methods to the CommercifyClient class in client.ts

// Guest Checkout API
async createGuestCheckout(data: CreateGuestCheckoutRequest): Promise<ResponseDTO<CheckoutDTO>> {
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

async addCheckoutItem(data: AddCheckoutItemRequest): Promise<ResponseDTO<CheckoutDTO>> {
  return this.request<ResponseDTO<CheckoutDTO>>(
    "/api/guest/checkout/items", 
    {
      method: "POST",
      body: JSON.stringify(data),
    }
  );
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

async removeCheckoutItem(productId: number): Promise<ResponseDTO<CheckoutDTO>> {
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
  return this.request<ResponseDTO<OrderDTO>>(
    "/api/guest/checkout/to-order",
    {
      method: "POST",
    }
  );
}

// Authenticated Checkout API
async getUserCheckout(): Promise<ResponseDTO<CheckoutDTO>> {
  return this.request<ResponseDTO<CheckoutDTO>>("/api/checkout", {
    method: "GET",
  });
}

async createCheckout(data: CreateCheckoutRequest): Promise<ResponseDTO<CheckoutDTO>> {
  return this.request<ResponseDTO<CheckoutDTO>>("/api/checkout", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

async convertGuestCheckoutToUserCheckout(): Promise<ResponseDTO<CheckoutDTO>> {
  return this.request<ResponseDTO<CheckoutDTO>>(
    "/api/checkout/convert",
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
  return this.request<ListResponseDTO<CheckoutDTO>>("/api/admin/checkouts", {}, params);
}

async getCheckoutById(checkoutId: number): Promise<ResponseDTO<CheckoutDTO>> {
  return this.request<ResponseDTO<CheckoutDTO>>(`/api/admin/checkouts/${checkoutId}`, {
    method: "GET",
  });
}

async getCheckoutsByUser(userId: number, params?: {
  page?: number;
  page_size?: number;
  status?: string;
}): Promise<ListResponseDTO<CheckoutDTO>> {
  return this.request<ListResponseDTO<CheckoutDTO>>(`/api/admin/users/${userId}/checkouts`, {}, params);
}
