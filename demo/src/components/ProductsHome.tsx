"use client";

import { ProductCard } from "@/components/ProductCard";
import { getProducts, getProductCategories } from "@/actions/products";
import { Product, ProductCategory } from "@/types";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";

export default function ProductsHome() {
  const searchParams = useSearchParams();
  const categoryId = searchParams.get("category");

  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<ProductCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [sortOption, setSortOption] = useState("featured");
  const [priceRange, setPriceRange] = useState({ min: "", max: "" });

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);

      try {
        // Use the server actions instead of direct API calls
        const params: { category?: string } = {};
        if (categoryId) {
          params.category = categoryId;
        }

        const [productsData, categoriesData] = await Promise.all([
          getProducts(params),
          getProductCategories(),
        ]);

        setProducts(productsData);
        setCategories(categoriesData);
      } catch (err) {
        console.error("Failed to fetch data:", err);
        setError("Failed to load products. Please try again later.");
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [categoryId]);

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    setSortOption(value);

    // Sort products based on selected option
    const sortedProducts = [...products];

    switch (value) {
      case "newest":
        sortedProducts.sort(
          (a, b) =>
            new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        );
        break;
      case "priceAsc":
        sortedProducts.sort((a, b) => a.price - b.price);
        break;
      case "priceDesc":
        sortedProducts.sort((a, b) => b.price - a.price);
        break;
      case "rating":
        sortedProducts.sort((a, b) => (b.rating || 0) - (a.rating || 0));
        break;
      // Default is "featured" which uses the server's order
    }

    setProducts(sortedProducts);
  };

  const handlePriceRangeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPriceRange((prev) => ({ ...prev, [name]: value }));
  };

  const applyFilters = () => {
    setLoading(true);

    // Filter products by price range
    let filteredProducts = [...products];

    const minPrice = priceRange.min ? parseFloat(priceRange.min) : null;
    const maxPrice = priceRange.max ? parseFloat(priceRange.max) : null;

    if (minPrice !== null || maxPrice !== null) {
      filteredProducts = filteredProducts.filter((product) => {
        if (minPrice !== null && maxPrice !== null) {
          return product.price >= minPrice && product.price <= maxPrice;
        } else if (minPrice !== null) {
          return product.price >= minPrice;
        } else if (maxPrice !== null) {
          return product.price <= maxPrice;
        }
        return true;
      });
    }

    setProducts(filteredProducts);
    setLoading(false);
  };

  return (
    <div className="bg-gray-50 min-h-screen">
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-800">Products</h1>
          <p className="text-gray-600 mt-2">
            Browse our catalog of premium products
          </p>
        </div>

        <div className="flex flex-col md:flex-row gap-6">
          {/* Category Sidebar */}
          <div className="w-full md:w-64 bg-white p-4 rounded-lg shadow-md h-fit">
            <h2 className="font-semibold text-lg mb-4">Categories</h2>
            <ul className="space-y-2">
              <li>
                <Link
                  href="/products"
                  className="block py-1 px-2 rounded hover:bg-blue-50 hover:text-blue-600 font-medium"
                >
                  All Products
                </Link>
              </li>
              {categories.map((category) => (
                <li key={category.id}>
                  <Link
                    href={`/products?category=${category.id}`}
                    className={`block py-1 px-2 rounded hover:bg-blue-50 hover:text-blue-600 ${
                      categoryId === category.id
                        ? "bg-blue-50 text-blue-600"
                        : ""
                    }`}
                  >
                    {category.name}
                  </Link>
                </li>
              ))}
            </ul>

            <div className="mt-8">
              <h2 className="font-semibold text-lg mb-4">Filters</h2>
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium mb-2">Price Range</h3>
                  <div className="flex items-center gap-2">
                    <input
                      type="number"
                      name="min"
                      placeholder="Min"
                      value={priceRange.min}
                      onChange={handlePriceRangeChange}
                      className="border rounded p-1 text-sm w-full"
                    />
                    <span>-</span>
                    <input
                      type="number"
                      name="max"
                      placeholder="Max"
                      value={priceRange.max}
                      onChange={handlePriceRangeChange}
                      className="border rounded p-1 text-sm w-full"
                    />
                  </div>
                </div>

                <div>
                  <h3 className="text-sm font-medium mb-2">Rating</h3>
                  <div className="flex items-center gap-1">
                    {[1, 2, 3, 4, 5].map((star) => (
                      <button
                        key={star}
                        className="text-gray-300 hover:text-yellow-500"
                      >
                        <svg
                          className="w-5 h-5"
                          fill="currentColor"
                          viewBox="0 0 20 20"
                          xmlns="http://www.w3.org/2000/svg"
                        >
                          <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"></path>
                        </svg>
                      </button>
                    ))}
                  </div>
                </div>

                <button
                  onClick={applyFilters}
                  className="w-full bg-blue-600 text-white rounded-md py-2 text-sm font-medium hover:bg-blue-700"
                >
                  Apply Filters
                </button>
              </div>
            </div>
          </div>

          {/* Products Grid */}
          <div className="flex-1">
            {error && (
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
            )}

            <div className="flex justify-between items-center mb-6">
              <div>
                <span className="text-gray-600">
                  {loading
                    ? "Loading products..."
                    : `${products.length} products`}
                </span>
              </div>
              <div>
                <select
                  className="border rounded-md p-2 bg-white"
                  value={sortOption}
                  onChange={handleSortChange}
                >
                  <option value="featured">Featured</option>
                  <option value="newest">Newest</option>
                  <option value="priceAsc">Price: Low to High</option>
                  <option value="priceDesc">Price: High to Low</option>
                  <option value="rating">Top Rated</option>
                </select>
              </div>
            </div>

            {loading ? (
              <div className="flex justify-center items-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
              </div>
            ) : products.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500 text-lg">No products found.</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                {products.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
