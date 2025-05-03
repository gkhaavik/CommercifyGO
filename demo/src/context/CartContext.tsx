"use client";

import { getProductById } from "@/actions/products";
import { Cart, CartItem, ProductVariant } from "@/types";
import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react";

interface CartContextType {
  cart: Cart;
  addToCart: (productId: string, quantity: number, variantId?: string) => void;
  updateCartItem: (itemId: string, quantity: number) => void;
  removeFromCart: (itemId: string) => void;
  clearCart: () => void;
  isLoading: boolean;
  error: string | null;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export function useCart() {
  const context = useContext(CartContext);
  if (!context) {
    throw new Error("useCart must be used within a CartProvider");
  }
  return context;
}

interface CartProviderProps {
  children: ReactNode;
}

export function CartProvider({ children }: CartProviderProps) {
  const [cart, setCart] = useState<Cart>({
    id: "cart-1",
    items: [],
    subtotal: 0,
    total: 0,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  });
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load cart from localStorage on component mount
  useEffect(() => {
    try {
      const savedCart = localStorage.getItem("commercify_cart");
      if (savedCart) {
        setCart(JSON.parse(savedCart));
      }
    } catch (err) {
      console.error("Failed to load cart from localStorage", err);
    }
  }, []);

  // Save cart to localStorage when it changes
  useEffect(() => {
    try {
      localStorage.setItem("commercify_cart", JSON.stringify(cart));
    } catch (err) {
      console.error("Failed to save cart to localStorage", err);
    }
  }, [cart]);

  // Recalculate totals whenever items change
  const calculateTotals = (items: CartItem[]) => {
    const subtotal = items.reduce(
      (sum, item) => sum + item.price * item.quantity,
      0
    );
    // Simple tax calculation - in a real app this would be more complex
    const tax = subtotal * 0.1;
    const total = subtotal + tax;

    return { subtotal, tax, total };
  };

  // Add item to cart
  const addToCart = async (
    productId: string,
    quantity: number,
    variantId?: string
  ) => {
    setIsLoading(true);
    setError(null);

    try {
      const product = await getProductById(productId);

      if (!product) {
        throw new Error("Product not found");
      }

      let selectedVariant: ProductVariant | undefined;
      let price = product.price;
      let stock = product.stock;

      // If product has variants and variantId is provided, find the variant
      if (product.has_variants && product.variants && variantId) {
        selectedVariant = product.variants.find((v) => v.id === variantId);

        if (!selectedVariant) {
          throw new Error("Variant not found");
        }

        price = selectedVariant.price;
        stock = selectedVariant.stock;
      }

      // Check if product/variant is in stock
      if (stock < quantity) {
        throw new Error("Not enough stock available");
      }

      // Check if product+variant combo is already in cart
      const existingItemIndex = cart.items.findIndex(
        (item) =>
          item.productId === productId &&
          ((!variantId && !item.variantId) || item.variantId === variantId)
      );

      let newItems: CartItem[];

      if (existingItemIndex > -1) {
        // Update quantity if product/variant is already in cart
        newItems = [...cart.items];
        newItems[existingItemIndex] = {
          ...newItems[existingItemIndex],
          quantity: newItems[existingItemIndex].quantity + quantity,
        };
      } else {
        // Add new item if product/variant is not in cart
        const newItem: CartItem = {
          id: `item-${Date.now()}`,
          productId,
          product,
          quantity,
          price,
          ...(variantId && { variantId, variant: selectedVariant }),
        };
        newItems = [...cart.items, newItem];
      }

      const { subtotal, tax, total } = calculateTotals(newItems);

      setCart((prevCart) => ({
        ...prevCart,
        items: newItems,
        subtotal,
        tax,
        total,
        updatedAt: new Date().toISOString(),
      }));
    } catch {
      setError("Failed to add item to cart");
    } finally {
      setIsLoading(false);
    }
  };

  // Update quantity of an item in cart
  const updateCartItem = async (itemId: string, quantity: number) => {
    setIsLoading(true);
    setError(null);

    try {
      const itemIndex = cart.items.findIndex((item) => item.id === itemId);

      if (itemIndex === -1) {
        throw new Error("Item not found in cart");
      }

      const cartItem = cart.items[itemIndex];
      const product = await getProductById(cartItem.productId);

      if (!product) {
        throw new Error("Product not found");
      }

      // Determine the actual stock based on variant (if any)
      let stock = product.stock;
      if (cartItem.variantId && product.variants) {
        const variant = product.variants.find(
          (v) => v.id === cartItem.variantId
        );
        if (variant) {
          stock = variant.stock;
        }
      }

      // Check if product is in stock
      if (stock < quantity) {
        throw new Error("Not enough stock available");
      }

      const newItems = [...cart.items];
      newItems[itemIndex] = {
        ...newItems[itemIndex],
        quantity,
      };

      const { subtotal, tax, total } = calculateTotals(newItems);

      setCart((prevCart) => ({
        ...prevCart,
        items: newItems,
        subtotal,
        tax,
        total,
        updatedAt: new Date().toISOString(),
      }));
    } catch {
      setError("Failed to update item in cart");
    } finally {
      setIsLoading(false);
    }
  };

  // Remove an item from cart
  const removeFromCart = (itemId: string) => {
    setIsLoading(true);
    setError(null);

    try {
      const newItems = cart.items.filter((item) => item.id !== itemId);

      const { subtotal, tax, total } = calculateTotals(newItems);

      setCart((prevCart) => ({
        ...prevCart,
        items: newItems,
        subtotal,
        tax,
        total,
        updatedAt: new Date().toISOString(),
      }));
    } catch {
      setError("Failed to remove item from cart");
    } finally {
      setIsLoading(false);
    }
  };

  // Clear the cart
  const clearCart = () => {
    setIsLoading(true);
    setError(null);

    try {
      setCart((prevCart) => ({
        ...prevCart,
        items: [],
        subtotal: 0,
        tax: 0,
        total: 0,
        updatedAt: new Date().toISOString(),
      }));
    } catch {
      setError("Failed to clear cart");
    } finally {
      setIsLoading(false);
    }
  };

  const value = {
    cart,
    addToCart,
    updateCartItem,
    removeFromCart,
    clearCart,
    isLoading,
    error,
  };

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}
