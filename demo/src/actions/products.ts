"use server";

import { Product, ProductCategory } from "@/types";

// Base URL for the Commercify API
const API_BASE_URL = "http://localhost:6091/api";

type Params = {
  category?: string;
  search?: string;
  page?: number;
  limit?: number;
};

/**
 * Get all products with optional filtering
 */
export async function getProducts(params?: Params): Promise<Product[]> {
  try {
    const queryString = params
      ? `?${new URLSearchParams(params as Record<string, string>).toString()}`
      : "";
    const response = await fetch(`${API_BASE_URL}/products${queryString}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch products: ${response.status}`);
    }

    return response.json();
  } catch (error) {
    console.error("Error fetching products:", error);
    throw error;
  }
}

/**
 * Get product categories
 */
export async function getProductCategories(): Promise<ProductCategory[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/categories`, {
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch product categories: ${response.status}`);
    }

    return response.json();
  } catch (error) {
    console.error("Error fetching product categories:", error);
    throw error;
  }
}

/**
 * Get a product by ID
 */
export async function getProductById(id: string): Promise<Product> {
  try {
    const response = await fetch(`${API_BASE_URL}/products/${id}`, {
      headers: {
        "Content-Type": "application/json",
      },
    });

    console.log("Response:", response);

    if (!response.ok) {
      throw new Error(`Failed to fetch product: ${response.status}`);
    }

    return response.json();
  } catch (error) {
    console.error("Error fetching product:", error);
    throw error;
  }
}

export async function getCategoryById(id: string): Promise<ProductCategory> {
  try {
    const response = await fetch(`${API_BASE_URL}/categories/${id}`, {
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch category: ${response.status}`);
    }

    return response.json();
  } catch (error) {
    console.error("Error fetching category:", error);
    throw error;
  }
}
