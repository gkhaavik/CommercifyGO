import Link from 'next/link'

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center p-8">
      <div className="z-10 w-full max-w-5xl items-center justify-between font-mono text-sm">
        <h1 className="text-4xl font-bold text-center mb-8">Commercify Demo</h1>
        <p className="text-center mb-12">
          Explore the powerful features of Commercify e-commerce platform
        </p>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
          <FeatureCard 
            title="Product Catalog" 
            description="Browse through our extensive product catalog with filtering and search capabilities."
            link="/products"
          />
          <FeatureCard 
            title="Shopping Cart" 
            description="Add products to your cart and manage your shopping experience."
            link="/cart" 
          />
          <FeatureCard 
            title="User Authentication" 
            description="Sign up, log in, and manage your user profile."
            link="/auth/login" 
          />
          <FeatureCard 
            title="Checkout Process" 
            description="Seamless checkout experience with multiple payment options."
            link="/checkout" 
          />
          <FeatureCard 
            title="Order Management" 
            description="Track and manage your orders in one place."
            link="/orders" 
          />
          <FeatureCard 
            title="Admin Panel" 
            description="Manage products, orders, and customers (Admin only)."
            link="/admin" 
          />
        </div>
      </div>
    </main>
  )
}

function FeatureCard({ title, description, link }: { title: string; description: string; link: string }) {
  return (
    <Link href={link} className="block">
      <div className="border border-gray-300 rounded-lg p-6 h-full hover:border-blue-500 hover:shadow-md transition-all">
        <h2 className="text-2xl font-semibold mb-2">{title}</h2>
        <p className="text-gray-600">{description}</p>
      </div>
    </Link>
  )
}
