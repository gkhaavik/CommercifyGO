import { Suspense } from "react";
import ProductsHome from "@/components/ProductsHome";

export default function ProductsPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <ProductsHome />
    </Suspense>
  );
}
