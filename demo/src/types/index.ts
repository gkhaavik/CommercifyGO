// Product Types
export interface Product {
  id: string;
  name: string;
  product_number: string;
  description: string;
  category_id: string;
  price: number;
  currency: string;
  images: string[];
  category: string;
  stock: number;
  rating?: number;
  reviewCount?: number;
  featured?: boolean;
  has_variants?: boolean;
  variants?: ProductVariant[];
  created_at: string;
  updated_at: string;
}

export interface VariantAttribute {
  name: string;
  value: string;
}

export interface ProductVariant {
  id: string;
  product_id: string;
  sku: string;
  category_id: string;
  price: number;
  compare_price?: number;
  stock: number;
  attributes: VariantAttribute[];
  images?: string[];
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export interface ProductCategory {
  id: string;
  name: string;
  description?: string;
  parentId?: string;
}

// User Types
export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: "customer" | "admin";
  addresses?: Address[];
  createdAt: string;
  updatedAt: string;
}

export interface Address {
  id: string;
  userId: string;
  type: "billing" | "shipping";
  name: string;
  line1: string;
  line2?: string;
  city: string;
  state?: string;
  postalCode: string;
  country: string;
  phone?: string;
  isDefault: boolean;
}

// Cart Types
export interface CartItem {
  id: string;
  productId: string;
  product: Product;
  variantId?: string;
  variant?: ProductVariant;
  quantity: number;
  price: number;
}

export interface Cart {
  id: string;
  userId?: string;
  items: CartItem[];
  subtotal: number;
  tax?: number;
  total: number;
  createdAt: string;
  updatedAt: string;
}

// Order Types
export type OrderStatus =
  | "pending"
  | "processing"
  | "shipped"
  | "delivered"
  | "cancelled";

export interface OrderItem {
  id: string;
  orderId: string;
  productId: string;
  product: Product;
  quantity: number;
  price: number;
  subtotal: number;
}

export interface Order {
  id: string;
  userId: string;
  items: OrderItem[];
  subtotal: number;
  tax: number;
  shipping: number;
  total: number;
  status: OrderStatus;
  shippingAddress: Address;
  billingAddress: Address;
  paymentMethod: string;
  paymentStatus: "pending" | "paid" | "failed";
  createdAt: string;
  updatedAt: string;
}

// Authentication Types
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}
