import { AddressDTO, CustomerDetailsDTO } from "./api";

export interface CheckoutItemDTO {
  id: number;
  product_id: number;
  variant_id?: number;
  product_name: string;
  variant_name?: string;
  sku: string;
  price: number;
  quantity: number;
  weight: number;
  subtotal: number;
  created_at: string;
  updated_at: string;
}

export interface AppliedDiscountDTO {
  discount_id: number;
  code: string;
  type: string;
  method: string;
  value: number;
  amount_saved: number;
}

export interface ShippingMethodDTO {
  id: number;
  name: string;
  description: string;
  price: number;
  estimated_delivery_days: number;
  carrier: string;
}

export interface CheckoutDTO {
  id: number;
  user_id?: number;
  session_id?: string;
  items: CheckoutItemDTO[];
  status: string;
  shipping_address: AddressDTO;
  billing_address: AddressDTO;
  shipping_method_id?: number;
  shipping_method?: ShippingMethodDTO;
  payment_provider?: string;
  total_amount: number;
  shipping_cost: number;
  total_weight: number;
  customer_details: CustomerDetailsDTO;
  currency: string;
  discount_code?: string;
  discount_amount: number;
  final_amount: number;
  applied_discount?: AppliedDiscountDTO;
  created_at: string;
  updated_at: string;
  last_activity_at: string;
  expires_at: string;
  completed_at?: string;
  converted_order_id?: number;
}

// Request interfaces
export interface CreateGuestCheckoutRequest {
  session_id: string;
  currency: string;
}

export interface CreateCheckoutRequest {
  currency: string;
}

export interface AddCheckoutItemRequest {
  product_id: number;
  variant_id?: number;
  quantity: number;
}

export interface UpdateCheckoutItemRequest {
  quantity: number;
}

export interface SetShippingAddressRequest {
  street: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
}

export interface SetBillingAddressRequest {
  street: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
}

export interface SetCustomerDetailsRequest {
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
}

export interface SetShippingMethodRequest {
  shipping_method_id: number;
}

export interface ApplyDiscountRequest {
  discount_code: string;
}
