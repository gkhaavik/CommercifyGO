"use server";

import { CartItem } from "@/types";

// Base URL for the Commercify API
const API_BASE_URL = "http://localhost:6091/api";

/**
 * Add an item to the cart
 */
export async function addToCart(
  productId: string,
  quantity: number,
  token?: string,
  variantId?: string
): Promise<CartItem> {
  try {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const payload = {
      productId,
      quantity,
      ...(variantId && { variantId }),
    };

    const response = await fetch(`${API_BASE_URL}/cart/items`, {
      method: "POST",
      headers,
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error(`Failed to add item to cart: ${response.status}`);
    }

    return response.json();
  } catch (error) {
    console.error("Error adding item to cart:", error);
    throw error;
  }
}

/**
 * Update cart item quantity
 */
export async function updateCartItem(
  itemId: string,
  quantity: number,
  token?: string
): Promise<CartItem> {
  try {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}/cart/items/${itemId}`, {
      method: "PATCH",
      headers,
      body: JSON.stringify({ quantity }),
    });

    if (!response.ok) {
      throw new Error(`Failed to update cart item: ${response.status}`);
    }

    return response.json();
  } catch (error) {
    console.error("Error updating cart item:", error);
    throw error;
  }
}

/**
 * Remove an item from the cart
 */
export async function removeCartItem(
  itemId: string,
  token?: string
): Promise<void> {
  try {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}/cart/items/${itemId}`, {
      method: "DELETE",
      headers,
    });

    if (!response.ok) {
      throw new Error(`Failed to remove cart item: ${response.status}`);
    }
  } catch (error) {
    console.error("Error removing cart item:", error);
    throw error;
  }
}

/**
 * Get the current cart contents
 */
export async function getCart(token?: string): Promise<{
  items: CartItem[];
  subtotal: number;
  tax: number;
  total: number;
}> {
  try {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}/cart`, {
      headers,
      cache: "no-store",
    });

    if (!response.ok) {
      throw new Error(`Failed to get cart: ${response.status}`);
    }

    const cart = await response.json();
    
    // If the API doesn't return calculated values, calculate them here
    if (!cart.subtotal) {
      cart.subtotal = cart.items.reduce(
        (sum: number, item: CartItem) => sum + item.price * item.quantity,
        0
      );
    }
    
    if (!cart.tax) {
      cart.tax = cart.subtotal * 0.08; // 8% tax rate for demo
    }
    
    if (!cart.total) {
      cart.total = cart.subtotal + cart.tax;
    }

    return cart;
  } catch (error) {
    console.error("Error getting cart:", error);
    // Return an empty cart on error
    return {
      items: [],
      subtotal: 0,
      tax: 0,
      total: 0,
    };
  }
}

/**
 * Clear the entire cart
 */
export async function clearCart(token?: string): Promise<void> {
  try {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE_URL}/cart`, {
      method: "DELETE",
      headers,
    });

    if (!response.ok) {
      throw new Error(`Failed to clear cart: ${response.status}`);
    }
  } catch (error) {
    console.error("Error clearing cart:", error);
    throw error;
  }
}