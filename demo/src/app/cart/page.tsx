"use client";

import { getCart, updateCartItem, removeCartItem } from "@/actions/cart";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Image from "next/image";
import Link from "next/link";
import { CartItem } from "@/types";

export default function CartPage() {
  const router = useRouter();
  const [cart, setCart] = useState<{
    items: CartItem[];
    subtotal: number;
    tax: number;
    total: number;
  }>({
    items: [],
    subtotal: 0,
    tax: 0,
    total: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [itemsLoading, setItemsLoading] = useState<Record<string, boolean>>({});

  // Fetch cart data when the component mounts
  useEffect(() => {
    const fetchCart = async () => {
      setLoading(true);
      setError(null);
      try {
        // Get token from localStorage if user is logged in
        const token =
          typeof window !== "undefined"
            ? localStorage.getItem("commercify_token")
            : null;
        const cartData = await getCart(token || undefined);
        setCart(cartData);
      } catch (err) {
        console.error("Failed to fetch cart:", err);
        setError("Failed to load cart. Please try again later.");
      } finally {
        setLoading(false);
      }
    };

    fetchCart();
  }, []);

  // Handle quantity change
  const handleQuantityChange = async (
    itemId: string,
    currentQuantity: number,
    newQuantity: number
  ) => {
    if (newQuantity < 1) return;

    // Optimistically update the UI
    setCart((prevCart) => ({
      ...prevCart,
      items: prevCart.items.map((item) =>
        item.id === itemId ? { ...item, quantity: newQuantity } : item
      ),
    }));

    // Mark this item as loading
    setItemsLoading((prev) => ({ ...prev, [itemId]: true }));

    try {
      // Get token from localStorage if user is logged in
      const token =
        typeof window !== "undefined"
          ? localStorage.getItem("commercify_token")
          : null;

      // Update the cart item on the server
      await updateCartItem(itemId, newQuantity, token || undefined);

      // Refresh the cart to get updated totals
      const updatedCart = await getCart(token || undefined);
      setCart(updatedCart);
    } catch (err) {
      console.error("Failed to update item quantity:", err);

      // Revert the optimistic update
      setCart((prevCart) => ({
        ...prevCart,
        items: prevCart.items.map((item) =>
          item.id === itemId ? { ...item, quantity: currentQuantity } : item
        ),
      }));

      setError("Failed to update quantity. Please try again.");
    } finally {
      // Mark the item as no longer loading
      setItemsLoading((prev) => ({ ...prev, [itemId]: false }));
    }
  };

  // Handle item removal
  const handleRemoveItem = async (itemId: string) => {
    // Optimistically update the UI
    const removedItem = cart.items.find((item) => item.id === itemId);
    setCart((prevCart) => ({
      ...prevCart,
      items: prevCart.items.filter((item) => item.id !== itemId),
    }));

    // Mark this item as loading
    setItemsLoading((prev) => ({ ...prev, [itemId]: true }));

    try {
      // Get token from localStorage if user is logged in
      const token =
        typeof window !== "undefined"
          ? localStorage.getItem("commercify_token")
          : null;

      // Remove the cart item on the server
      await removeCartItem(itemId, token || undefined);

      // Refresh the cart to get updated totals
      const updatedCart = await getCart(token || undefined);
      setCart(updatedCart);
    } catch (err) {
      console.error("Failed to remove item:", err);

      // Revert the optimistic update if we have the removed item
      if (removedItem) {
        setCart((prevCart) => ({
          ...prevCart,
          items: [...prevCart.items, removedItem],
        }));
      }

      setError("Failed to remove item. Please try again.");
    } finally {
      // Mark the item as no longer loading
      setItemsLoading((prev) => ({ ...prev, [itemId]: false }));
    }
  };

  const navigateToCheckout = () => {
    router.push("/checkout");
  };

  // Format currency
  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(amount);
  };

  if (loading) {
    return (
      <div className="bg-gray-50 min-h-screen py-8 flex justify-center items-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-gray-50 min-h-screen py-8">
        <div className="max-w-3xl mx-auto px-4">
          <div className="bg-red-50 border-l-4 border-red-400 p-4 mb-6">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg
                  className="h-5 w-5 text-red-400"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    fillRule="evenodd"
                    d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                    clipRule="evenodd"
                  />
                </svg>
              </div>
              <div className="ml-3">
                <p className="text-sm text-red-700">{error}</p>
              </div>
            </div>
          </div>
          <div className="flex justify-center">
            <button
              onClick={() => window.location.reload()}
              className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
            >
              Try Again
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (cart.items.length === 0) {
    return (
      <div className="bg-gray-50 min-h-screen py-8">
        <div className="max-w-3xl mx-auto px-4">
          <div className="bg-white rounded-lg shadow-md p-8 text-center">
            <svg
              className="w-16 h-16 text-gray-400 mx-auto mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
              ></path>
            </svg>
            <h2 className="text-2xl font-bold text-gray-800">
              Your cart is empty
            </h2>
            <p className="text-gray-600 mt-2">
              Looks like you haven&apos;t added any items to your cart yet.
            </p>
            <div className="mt-6">
              <Link
                href="/products"
                className="bg-blue-600 text-white px-6 py-3 rounded-md hover:bg-blue-700 inline-block"
              >
                Browse Products
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-50 min-h-screen py-8">
      <div className="max-w-6xl mx-auto px-4">
        <h1 className="text-3xl font-bold text-gray-800 mb-8">Shopping Cart</h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Cart Items */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow-md overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Product
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Price
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Quantity
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Total
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      <span className="sr-only">Actions</span>
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {cart.items.map((item) => (
                    <tr
                      key={item.id}
                      className={itemsLoading[item.id] ? "opacity-50" : ""}
                    >
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center">
                          <div className="flex-shrink-0 h-14 w-14 relative">
                            <Image
                              src={
                                item.product.images[0] ||
                                "https://via.placeholder.com/150"
                              }
                              alt={item.product.name}
                              className="object-cover"
                              fill
                            />
                          </div>
                          <div className="ml-4">
                            <div className="text-sm font-medium text-gray-900">
                              {item.product.name}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-900">
                          {formatCurrency(item.price)}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center">
                          <button
                            onClick={() =>
                              handleQuantityChange(
                                item.id,
                                item.quantity,
                                item.quantity - 1
                              )
                            }
                            disabled={
                              item.quantity <= 1 || itemsLoading[item.id]
                            }
                            className="p-1 rounded-md hover:bg-gray-100 disabled:opacity-50"
                          >
                            <svg
                              className="w-4 h-4 text-gray-500"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                d="M20 12H4"
                              />
                            </svg>
                          </button>
                          <span className="mx-2 text-gray-600">
                            {item.quantity}
                          </span>
                          <button
                            onClick={() =>
                              handleQuantityChange(
                                item.id,
                                item.quantity,
                                item.quantity + 1
                              )
                            }
                            disabled={itemsLoading[item.id]}
                            className="p-1 rounded-md hover:bg-gray-100 disabled:opacity-50"
                          >
                            <svg
                              className="w-4 h-4 text-gray-500"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                d="M12 4v16m8-8H4"
                              />
                            </svg>
                          </button>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm font-medium text-gray-900">
                          {formatCurrency(item.price * item.quantity)}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <button
                          onClick={() => handleRemoveItem(item.id)}
                          disabled={itemsLoading[item.id]}
                          className="text-red-600 hover:text-red-800 disabled:opacity-50"
                        >
                          <svg
                            className="w-5 h-5"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth="2"
                              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            />
                          </svg>
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Order Summary */}
          <div>
            <div className="bg-white rounded-lg shadow-md p-6">
              <h2 className="text-lg font-semibold text-gray-800 mb-4">
                Order Summary
              </h2>

              <div className="space-y-3 border-b border-gray-200 pb-4">
                <div className="flex justify-between">
                  <span className="text-gray-600">Subtotal</span>
                  <span className="text-gray-800 font-medium">
                    {formatCurrency(cart.subtotal)}
                  </span>
                </div>

                <div className="flex justify-between">
                  <span className="text-gray-600">Tax</span>
                  <span className="text-gray-800 font-medium">
                    {formatCurrency(cart.tax)}
                  </span>
                </div>

                <div className="flex justify-between">
                  <span className="text-gray-600">Shipping</span>
                  <span className="text-gray-800 font-medium">Free</span>
                </div>
              </div>

              <div className="flex justify-between items-center mt-4">
                <span className="text-lg font-semibold text-gray-800">
                  Total
                </span>
                <span className="text-xl font-bold text-gray-900">
                  {formatCurrency(cart.total)}
                </span>
              </div>

              <button
                onClick={navigateToCheckout}
                className="w-full bg-blue-600 text-white rounded-md py-3 mt-6 font-medium hover:bg-blue-700 transition"
              >
                Proceed to Checkout
              </button>

              <div className="flex justify-center mt-4">
                <Link
                  href="/products"
                  className="text-blue-600 hover:text-blue-800 text-sm font-medium"
                >
                  Continue Shopping
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
