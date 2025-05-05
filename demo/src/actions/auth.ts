"use server";

import { User } from "@/types";

// Base URL for the Commercify API
const API_BASE_URL = "http://localhost:6091/api";

/**
 * User login
 */
export async function login(email: string, password: string): Promise<{
  user: User;
  token: string;
}> {
  try {
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email, password }),
      cache: "no-store",
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(
        errorData.message || `Login failed with status ${response.status}`
      );
    }

    return response.json();
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
}

/**
 * User registration
 */
export async function register(userData: {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
}): Promise<{
  user: User;
  token: string;
}> {
  try {
    const response = await fetch(`${API_BASE_URL}/auth/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(userData),
      cache: "no-store",
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(
        errorData.message || `Registration failed with status ${response.status}`
      );
    }

    return response.json();
  } catch (error) {
    console.error("Registration error:", error);
    throw error;
  }
}

/**
 * Get current user profile
 */
export async function getUserProfile(token: string): Promise<User> {
  try {
    const response = await fetch(`${API_BASE_URL}/users/me`, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      cache: "no-store",
    });

    if (!response.ok) {
      throw new Error(`Failed to get user profile: ${response.status}`);
    }

    return response.json();
  } catch (error) {
    console.error("Error getting user profile:", error);
    throw error;
  }
}