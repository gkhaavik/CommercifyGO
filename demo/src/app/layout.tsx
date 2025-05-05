import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
import Link from "next/link";
import { CartProvider } from "@/context/CartContext";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Commercify Demo",
  description: "Demo application showcasing Commercify e-commerce features",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <CartProvider>
          <nav className="bg-white shadow-md">
            <div className="max-w-7xl mx-auto px-4">
              <div className="flex justify-between h-16">
                <div className="flex">
                  <div className="flex-shrink-0 flex items-center">
                    <Link href="/" className="text-xl font-bold text-blue-600">
                      Commercify
                    </Link>
                  </div>
                  <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                    <Link
                      href="/products"
                      className="border-transparent text-gray-500 hover:border-blue-500 hover:text-blue-500 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                    >
                      Products
                    </Link>
                    <Link
                      href="/cart"
                      className="border-transparent text-gray-500 hover:border-blue-500 hover:text-blue-500 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                    >
                      Cart
                    </Link>
                    <Link
                      href="/orders"
                      className="border-transparent text-gray-500 hover:border-blue-500 hover:text-blue-500 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                    >
                      Orders
                    </Link>
                  </div>
                </div>
                <div className="hidden sm:ml-6 sm:flex sm:items-center">
                  <Link
                    href="/auth/login"
                    className="border-transparent text-gray-500 hover:border-blue-500 hover:text-blue-500 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
                  >
                    Login
                  </Link>
                </div>
              </div>
            </div>
          </nav>
          {children}
          <footer className="bg-gray-800 text-white py-8">
            <div className="max-w-7xl mx-auto px-4">
              <div className="flex flex-col md:flex-row justify-between">
                <div className="mb-4 md:mb-0">
                  <h3 className="text-lg font-bold">Commercify Demo</h3>
                  <p className="text-gray-300 text-sm mt-2">
                    A demonstration of Commercify&apos;s e-commerce features
                  </p>
                </div>
                <div>
                  <h4 className="font-semibold mb-2">Demo Features</h4>
                  <ul className="text-sm text-gray-300">
                    <li>
                      <Link href="/products" className="hover:text-white">
                        Products
                      </Link>
                    </li>
                    <li>
                      <Link href="/cart" className="hover:text-white">
                        Cart
                      </Link>
                    </li>
                    <li>
                      <Link href="/auth/login" className="hover:text-white">
                        Authentication
                      </Link>
                    </li>
                    <li>
                      <Link href="/checkout" className="hover:text-white">
                        Checkout
                      </Link>
                    </li>
                  </ul>
                </div>
              </div>
              <div className="mt-8 border-t border-gray-700 pt-4 text-sm text-gray-400">
                <p>
                  &copy; {new Date().getFullYear()} Commercify Demo. All rights
                  reserved.
                </p>
              </div>
            </div>
          </footer>
        </CartProvider>
      </body>
    </html>
  );
}
