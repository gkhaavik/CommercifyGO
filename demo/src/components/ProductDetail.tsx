"use client";

import { addToCart } from "@/actions/cart";
import { Product, ProductVariant } from "@/types";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

export default function ProductDetail({ product }: { product: Product }) {
  const router = useRouter();
  const [selectedImage, setSelectedImage] = useState(0);
  const [quantity, setQuantity] = useState(1);
  const [addingToCart, setAddingToCart] = useState(false);
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(
    null
  );
  const [selectedAttributes, setSelectedAttributes] = useState<
    Record<string, string>
  >({});
  const [cartMessage, setCartMessage] = useState<{
    type: "success" | "error";
    text: string;
  } | null>(null);

  // Get unique attribute names across all variants
  const attributeNames = useMemo(() => {
    if (
      !product.has_variants ||
      !product.variants ||
      product.variants.length === 0
    ) {
      return [];
    }

    const names = new Set<string>();
    product.variants.forEach((variant) => {
      variant.attributes.forEach((attr) => {
        names.add(attr.name);
      });
    });

    return Array.from(names);
  }, [product.has_variants, product.variants]);

  // Get available values for each attribute
  const attributeValues = useMemo(() => {
    if (
      !product.has_variants ||
      !product.variants ||
      product.variants.length === 0
    ) {
      return {};
    }

    const values: Record<string, Set<string>> = {};

    attributeNames.forEach((name) => {
      values[name] = new Set<string>();
    });

    product.variants.forEach((variant) => {
      variant.attributes.forEach((attr) => {
        if (values[attr.name]) {
          values[attr.name].add(attr.value);
        }
      });
    });

    // Convert Sets to Arrays for easier rendering
    const result: Record<string, string[]> = {};
    Object.keys(values).forEach((key) => {
      result[key] = Array.from(values[key]);
    });

    return result;
  }, [product.has_variants, product.variants, attributeNames]);

  // Find the matching variant when attributes are selected
  useEffect(() => {
    if (
      !product.has_variants ||
      !product.variants ||
      product.variants.length === 0
    ) {
      return;
    }

    // Find default variant if no attributes are selected yet
    if (Object.keys(selectedAttributes).length === 0) {
      const defaultVariant = product.variants.find((v) => v.is_default);
      if (defaultVariant) {
        setSelectedVariant(defaultVariant);
        // Pre-select the default variant attributes
        const defaultAttributes: Record<string, string> = {};
        defaultVariant.attributes.forEach((attr) => {
          defaultAttributes[attr.name] = attr.value;
        });
        setSelectedAttributes(defaultAttributes);
      }
      return;
    }

    // Find variant that matches all selected attributes
    const matchingVariant = product.variants.find((variant) => {
      // For each selected attribute, check if variant has matching attribute value
      return Object.entries(selectedAttributes).every(([name, value]) => {
        return variant.attributes.some(
          (attr) => attr.name === name && attr.value === value
        );
      });
    });

    setSelectedVariant(matchingVariant || null);
  }, [product.has_variants, product.variants, selectedAttributes]);

  // When variant changes, update the selected image if the variant has images
  useEffect(() => {
    if (
      selectedVariant &&
      selectedVariant.images &&
      selectedVariant.images.length > 0
    ) {
      setSelectedImage(0); // Reset to first image of the variant
    }
  }, [selectedVariant]);

  const handleAttributeChange = (name: string, value: string) => {
    setSelectedAttributes((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleAddToCart = async () => {
    if (!product) return;

    setAddingToCart(true);
    setCartMessage(null);

    try {
      // Get token from localStorage if user is logged in
      const token =
        typeof window !== "undefined"
          ? localStorage.getItem("commercify_token")
          : null;

      // Use the server action to add to cart with the selected variant if available
      await addToCart(
        product.id,
        quantity,
        token || undefined,
        selectedVariant ? selectedVariant.id : undefined
      );

      setCartMessage({
        type: "success",
        text: `${quantity} ${
          quantity > 1 ? "items" : "item"
        } added to your cart.`,
      });

      // Reset quantity after successful add
      setQuantity(1);
    } catch (err) {
      console.error("Failed to add to cart:", err);
      setCartMessage({
        type: "error",
        text: "Failed to add to cart. Please try again.",
      });
    } finally {
      setAddingToCart(false);
    }
  };

  const handleQuantityChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!product) return;

    const value = parseInt(e.target.value);
    const maxStock = selectedVariant ? selectedVariant.stock : product.stock;

    if (value > 0 && value <= maxStock) {
      setQuantity(value);
    }
  };

  const navigateToCart = () => {
    router.push("/cart");
  };

  console.log("Product Detail:", product);
  console.log("Selected Variant:", selectedVariant);

  // Display images based on the selected variant or product
  const displayImages = useMemo(() => {
    if (
      selectedVariant &&
      selectedVariant.images &&
      selectedVariant.images.length > 0
    ) {
      return selectedVariant.images;
    }
    return product.images;
  }, [selectedVariant, product.images]);

  // Get current price and stock based on variant selection
  const currentPrice = selectedVariant ? selectedVariant.price : product.price;
  const currentStock = selectedVariant ? selectedVariant.stock : product.stock;
  const comparePrice = selectedVariant?.compare_price;

  // Format prices with currency
  const formattedPrice = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: product.currency || "USD",
  }).format(currentPrice);

  const formattedComparePrice = comparePrice
    ? new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: product.currency || "USD",
      }).format(comparePrice)
    : null;

  return (
    <div className="bg-gray-50 min-h-screen py-8">
      <div className="max-w-7xl mx-auto px-4">
        {cartMessage && (
          <div
            className={`mb-4 p-4 rounded-md ${
              cartMessage.type === "success"
                ? "bg-green-50 border-green-500"
                : "bg-red-50 border-red-500"
            } border-l-4`}
          >
            <div className="flex">
              <div className="flex-shrink-0">
                {cartMessage.type === "success" ? (
                  <svg
                    className="h-5 w-5 text-green-400"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                      clipRule="evenodd"
                    />
                  </svg>
                ) : (
                  <svg
                    className="h-5 w-5 text-red-400"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                      clipRule="evenodd"
                    />
                  </svg>
                )}
              </div>
              <div className="ml-3 flex justify-between items-center w-full">
                <p
                  className={`text-sm ${
                    cartMessage.type === "success"
                      ? "text-green-700"
                      : "text-red-700"
                  }`}
                >
                  {cartMessage.text}
                </p>
                {cartMessage.type === "success" && (
                  <button
                    onClick={navigateToCart}
                    className="text-sm text-blue-600 hover:text-blue-800 font-medium"
                  >
                    View Cart
                  </button>
                )}
              </div>
            </div>
          </div>
        )}

        <div className="bg-white rounded-lg shadow-md p-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            {/* Product Images */}
            <div>
              <div className="relative h-96 rounded-lg overflow-hidden mb-4">
                <Image
                  src={"/" + displayImages[selectedImage]}
                  alt={product.name}
                  className="object-contain"
                  fill
                  priority
                />
              </div>

              <div className="flex gap-2 overflow-x-auto pb-2">
                {displayImages.map((image, index) => (
                  <div
                    key={index}
                    className={`relative w-20 h-20 border-2 rounded cursor-pointer ${
                      selectedImage === index
                        ? "border-blue-500"
                        : "border-transparent"
                    }`}
                    onClick={() => setSelectedImage(index)}
                  >
                    <Image
                      src={"/" + image}
                      alt={`${product.name} thumbnail ${index + 1}`}
                      className="object-cover"
                      fill
                    />
                  </div>
                ))}
              </div>
            </div>

            {/* Product Details */}
            <div>
              <h1 className="text-3xl font-bold text-gray-800">
                {product.name}
              </h1>

              <div className="mt-4 flex items-center">
                {product.rating && (
                  <div className="flex items-center">
                    <div className="flex">
                      {[...Array(5)].map((_, i) => (
                        <svg
                          key={i}
                          className={`w-5 h-5 ${
                            i < Math.floor(product.rating || 0)
                              ? "text-yellow-500"
                              : "text-gray-300"
                          }`}
                          fill="currentColor"
                          viewBox="0 0 20 20"
                          xmlns="http://www.w3.org/2000/svg"
                        >
                          <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"></path>
                        </svg>
                      ))}
                    </div>
                    <span className="text-gray-600 ml-2">
                      {product.rating} ({product.reviewCount || 0} reviews)
                    </span>
                  </div>
                )}
              </div>

              <div className="mt-6">
                <div className="flex items-center">
                  <h2
                    className={`text-xl font-bold text-gray-800 ${
                      formattedComparePrice ? "text-red-600" : ""
                    }`}
                  >
                    {formattedPrice}
                  </h2>

                  {formattedComparePrice && (
                    <span className="ml-2 text-gray-500 line-through">
                      {formattedComparePrice}
                    </span>
                  )}
                </div>

                <div className="mt-2">
                  <span
                    className={`inline-block px-2 py-1 rounded-md text-sm ${
                      currentStock > 10
                        ? "bg-green-100 text-green-800"
                        : currentStock > 0
                        ? "bg-orange-100 text-orange-800"
                        : "bg-red-100 text-red-800"
                    }`}
                  >
                    {currentStock > 10
                      ? "In Stock"
                      : currentStock > 0
                      ? `Low Stock: ${currentStock} remaining`
                      : "Out of Stock"}
                  </span>
                </div>
              </div>

              <div className="mt-6">
                <p className="text-gray-700">{product.description}</p>
              </div>

              {/* Variant Selection */}
              {product.has_variants &&
                product.variants &&
                product.variants.length > 0 && (
                  <div className="mt-6 border-t border-gray-200 pt-4">
                    {attributeNames.map((attrName) => (
                      <div key={attrName} className="mb-4">
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                          {attrName.charAt(0).toUpperCase() + attrName.slice(1)}
                        </label>
                        <div className="flex flex-wrap gap-2">
                          {attributeValues[attrName]?.map((value) => (
                            <button
                              key={`${attrName}-${value}`}
                              onClick={() =>
                                handleAttributeChange(attrName, value)
                              }
                              className={`px-3 py-1 border rounded-md text-sm 
                              ${
                                selectedAttributes[attrName] === value
                                  ? "border-blue-500 bg-blue-50 text-blue-700"
                                  : "border-gray-300 hover:border-gray-400"
                              }`}
                            >
                              {value}
                            </button>
                          ))}
                        </div>
                      </div>
                    ))}

                    {!selectedVariant &&
                      Object.keys(selectedAttributes).length > 0 && (
                        <p className="text-red-500 text-sm mt-2">
                          This combination is not available
                        </p>
                      )}
                  </div>
                )}

              {currentStock > 0 && (
                <div className="mt-8 flex items-center">
                  <span className="mr-4">Quantity:</span>
                  <div className="flex items-center border border-gray-300 rounded">
                    <button
                      className="px-3 py-1 border-r border-gray-300"
                      onClick={() => quantity > 1 && setQuantity(quantity - 1)}
                      disabled={addingToCart}
                    >
                      -
                    </button>
                    <input
                      type="number"
                      className="w-12 py-1 text-center border-none focus:outline-none"
                      value={quantity}
                      min={1}
                      max={currentStock}
                      onChange={handleQuantityChange}
                      disabled={addingToCart}
                    />
                    <button
                      className="px-3 py-1 border-l border-gray-300"
                      onClick={() =>
                        quantity < currentStock && setQuantity(quantity + 1)
                      }
                      disabled={addingToCart}
                    >
                      +
                    </button>
                  </div>
                </div>
              )}

              <div className="mt-8">
                <button
                  onClick={handleAddToCart}
                  disabled={
                    currentStock === 0 ||
                    addingToCart ||
                    (product.has_variants && !selectedVariant)
                  }
                  className={`w-full py-3 px-4 rounded-md text-white font-medium ${
                    currentStock === 0 ||
                    addingToCart ||
                    (product.has_variants && !selectedVariant)
                      ? "bg-gray-300 cursor-not-allowed"
                      : "bg-blue-600 hover:bg-blue-700"
                  }`}
                >
                  {addingToCart ? (
                    <span className="flex items-center justify-center">
                      <svg
                        className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                      >
                        <circle
                          className="opacity-25"
                          cx="12"
                          cy="12"
                          r="10"
                          stroke="currentColor"
                          strokeWidth="4"
                        ></circle>
                        <path
                          className="opacity-75"
                          fill="currentColor"
                          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        ></path>
                      </svg>
                      Adding to Cart...
                    </span>
                  ) : currentStock === 0 ? (
                    "Out of Stock"
                  ) : product.has_variants && !selectedVariant ? (
                    "Select Options"
                  ) : (
                    "Add to Cart"
                  )}
                </button>
              </div>

              <div className="mt-8 border-t border-gray-200 pt-4">
                <h3 className="font-semibold text-gray-800">Product Details</h3>
                <ul className="mt-2 space-y-1 text-gray-600">
                  <li>Category: {product.category}</li>
                  {selectedVariant ? (
                    <li>SKU: {selectedVariant.sku}</li>
                  ) : (
                    <li>ID: {product.id}</li>
                  )}
                  <li>
                    Added on: {new Date(product.created_at).toLocaleDateString()}
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
