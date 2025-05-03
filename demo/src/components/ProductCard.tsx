"use client";

import { Product } from "@/types";
import Image from "next/image";
import Link from "next/link";
import { useState } from "react";
import { useCart } from "@/context/CartContext";

interface ProductCardProps {
  product: Product;
}

export function ProductCard({ product }: ProductCardProps) {
  const [isHovered, setIsHovered] = useState(false);
  const { addToCart, isLoading } = useCart();

  // Check if product has variants
  const hasVariants = product.has_variants && product.variants && product.variants.length > 0;
  
  // Get default variant for price display if product has variants
  const defaultVariant = hasVariants 
    ? product.variants?.find(variant => variant.is_default)
    : null;
    
  // Determine the price to display
  const displayPrice = defaultVariant ? defaultVariant.price : product.price;
  
  // Format the price with currency
  const formattedPrice = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: product.currency || "USD",
  }).format(displayPrice);

  // Get compare price if available
  const comparePrice = defaultVariant?.compare_price;
  const formattedComparePrice = comparePrice 
    ? new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: product.currency || "USD",
      }).format(comparePrice)
    : null;

  // Determine the stock to display
  const displayStock = defaultVariant ? defaultVariant.stock : product.stock;

  const handleAddToCart = () => {
    // If product has variants, we should redirect to product detail page
    // instead of adding to cart directly
    if (hasVariants) {
      return;
    }
    
    addToCart(product.id, 1);
  };

  return (
    <div
      className="bg-white rounded-lg shadow-md overflow-hidden transition-all duration-300 h-full flex flex-col"
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <Link href={`/products/${product.id}`}>
        <div className="relative h-56 overflow-hidden">
          <Image
            src={"/" + product.images[0] || "/placeholder.png"}
            alt={product.name}
            className={`object-cover w-full transition-all duration-500 ${
              isHovered ? "scale-110" : "scale-100"
            }`}
            width={500}
            height={300}
          />
          {displayStock < 10 && displayStock > 0 && (
            <div className="absolute top-2 right-2 bg-orange-500 text-white text-xs px-2 py-1 rounded-md">
              Low Stock
            </div>
          )}
          {displayStock === 0 && (
            <div className="absolute top-2 right-2 bg-red-500 text-white text-xs px-2 py-1 rounded-md">
              Out of Stock
            </div>
          )}
          {hasVariants && (
            <div className="absolute bottom-2 left-2 bg-blue-600 text-white text-xs px-2 py-1 rounded-md">
              Options Available
            </div>
          )}
        </div>
      </Link>

      <div className="p-4 flex-grow flex flex-col">
        <Link href={`/products/${product.id}`} className="flex-grow">
          <h3 className="text-lg font-semibold text-gray-800 hover:text-blue-600">
            {product.name}
          </h3>
          <p className="text-gray-500 text-sm mt-1 line-clamp-2">
            {product.description}
          </p>
        </Link>

        <div className="mt-4">
          <div className="flex items-center justify-between">
            <div>
              <span className={`font-bold text-lg ${formattedComparePrice ? 'text-red-600' : ''}`}>
                {formattedPrice}
              </span>
              {formattedComparePrice && (
                <span className="ml-2 text-gray-500 text-sm line-through">
                  {formattedComparePrice}
                </span>
              )}
            </div>
            <div className="flex items-center">
              {product.rating && (
                <div className="flex items-center">
                  <svg
                    className="w-4 h-4 text-yellow-500 mr-1"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"></path>
                  </svg>
                  <span className="text-sm text-gray-600">
                    {product.rating}
                  </span>
                </div>
              )}
            </div>
          </div>

          {hasVariants ? (
            <Link href={`/products/${product.id}`} className="block w-full">
              <button
                className={`mt-3 w-full py-2 px-4 rounded-md text-white font-medium bg-blue-600 hover:bg-blue-700`}
              >
                View Options
              </button>
            </Link>
          ) : (
            <button
              onClick={handleAddToCart}
              disabled={displayStock === 0 || isLoading}
              className={`mt-3 w-full py-2 px-4 rounded-md text-white font-medium ${
                displayStock === 0
                  ? "bg-gray-300 cursor-not-allowed"
                  : isLoading
                  ? "bg-blue-400 cursor-not-allowed"
                  : "bg-blue-600 hover:bg-blue-700"
              }`}
            >
              {displayStock === 0
                ? "Out of Stock"
                : isLoading
                ? "Adding..."
                : "Add to Cart"}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
