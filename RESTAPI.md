### E-Commerce API Documentation

Below is a comprehensive list of all API endpoints in the e-commerce system, including their request and response bodies.

## Table of Contents

- [Authentication](#authentication)
- [Users](#users)
- [Products](#products)
- [Categories](#categories)
- [Cart](#cart)
- [Orders](#orders)
- [Payment](#payment)
- [Webhooks](#webhooks)

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```plaintext
Authorization: Bearer <token>
```

## Users

### Register User

```plaintext
POST /api/users/register
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response Body:**

```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user",
    "created_at": "2023-04-20T12:00:00Z",
    "updated_at": "2023-04-20T12:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Status Codes:**

- `201 Created`: User registered successfully
- `400 Bad Request`: Invalid request body or email already exists

### Login

```plaintext
POST /api/users/login
```

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response Body:**

```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "user",
    "created_at": "2023-04-20T12:00:00Z",
    "updated_at": "2023-04-20T12:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Status Codes:**

- `200 OK`: Login successful
- `401 Unauthorized`: Invalid email or password

### Get User Profile

```plaintext
GET /api/users/me
```

**Response Body:**

```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "user",
  "created_at": "2023-04-20T12:00:00Z",
  "updated_at": "2023-04-20T12:00:00Z"
}
```

**Status Codes:**

- `200 OK`: Profile retrieved successfully
- `401 Unauthorized`: Not authenticated

### Update User Profile

```plaintext
PUT /api/users/me
```

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Smith"
}
```

**Response Body:**

```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Smith",
  "role": "user",
  "created_at": "2023-04-20T12:00:00Z",
  "updated_at": "2023-04-20T12:30:00Z"
}
```

**Status Codes:**

- `200 OK`: Profile updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated

### Change Password

```plaintext
PUT /api/users/me/password
```

**Request Body:**

```json
{
  "current_password": "password123",
  "new_password": "newpassword123"
}
```

**Response Body:**

```json
{
  "message": "Password changed successfully"
}
```

**Status Codes:**

- `200 OK`: Password changed successfully
- `400 Bad Request`: Invalid request body or current password is incorrect
- `401 Unauthorized`: Not authenticated

### List Users (Admin Only)

```plaintext
GET /api/admin/users
```

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

**Response Body:**

```json
[
  {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Smith",
    "role": "user",
    "created_at": "2023-04-20T12:00:00Z",
    "updated_at": "2023-04-20T12:30:00Z"
  },
  {
    "id": 2,
    "email": "admin@example.com",
    "first_name": "Admin",
    "last_name": "User",
    "role": "admin",
    "created_at": "2023-04-19T10:00:00Z",
    "updated_at": "2023-04-19T10:00:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Users retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

## Products

### List Products

```plaintext
GET /api/products
```

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

**Response Body:**

```json
[
  {
    "id": 1,
    "name": "Smartphone",
    "description": "Latest smartphone model",
    "price": 999.99,
    "stock": 50,
    "category_id": 1,
    "seller_id": 2,
    "images": ["smartphone.jpg"],
    "has_variants": false,
    "created_at": "2023-04-15T10:00:00Z",
    "updated_at": "2023-04-15T10:00:00Z"
  },
  {
    "id": 2,
    "name": "Laptop",
    "description": "Powerful laptop for professionals",
    "price": 1499.99,
    "stock": 25,
    "category_id": 1,
    "seller_id": 2,
    "images": ["laptop.jpg"],
    "has_variants": true,
    "created_at": "2023-04-16T11:00:00Z",
    "updated_at": "2023-04-16T11:00:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Products retrieved successfully

### Get Product

```plaintext
GET /api/products/{id}
```

**Response Body:**

```json
{
  "id": 1,
  "name": "Smartphone",
  "description": "Latest smartphone model",
  "price": 999.99,
  "stock": 50,
  "category_id": 1,
  "seller_id": 2,
  "images": ["smartphone.jpg"],
  "has_variants": false,
  "created_at": "2023-04-15T10:00:00Z",
  "updated_at": "2023-04-15T10:00:00Z"
}
```

**Status Codes:**

- `200 OK`: Product retrieved successfully
- `404 Not Found`: Product not found

### Search Products

```plaintext
GET /api/products/search
```

**Query Parameters:**

- `q` (optional): Search query
- `category` (optional): Category ID
- `min_price` (optional): Minimum price
- `max_price` (optional): Maximum price
- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

**Response Body:**

```json
[
  {
    "id": 1,
    "name": "Smartphone",
    "description": "Latest smartphone model",
    "price": 999.99,
    "stock": 50,
    "category_id": 1,
    "seller_id": 2,
    "images": ["smartphone.jpg"],
    "has_variants": false,
    "created_at": "2023-04-15T10:00:00Z",
    "updated_at": "2023-04-15T10:00:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Search results retrieved successfully

### Create Product (Seller Only)

```plaintext
POST /api/products
```

**Request Body:**

```json
{
  "name": "New Product",
  "description": "Product description",
  "price": 199.99,
  "stock": 100,
  "category_id": 1,
  "images": ["product.jpg"],
  "has_variants": false
}
```

**Response Body:**

```json
{
  "id": 3,
  "name": "New Product",
  "description": "Product description",
  "price": 199.99,
  "stock": 100,
  "category_id": 1,
  "seller_id": 2,
  "images": ["product.jpg"],
  "has_variants": false,
  "created_at": "2023-04-20T14:00:00Z",
  "updated_at": "2023-04-20T14:00:00Z"
}
```

**Status Codes:**

- `201 Created`: Product created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not a seller)

### Update Product (Seller Only)

```plaintext
PUT /api/products/{id}
```

**Request Body:**

```json
{
  "name": "Updated Product",
  "description": "Updated description",
  "price": 249.99,
  "stock": 75,
  "category_id": 1,
  "images": ["updated-product.jpg"]
}
```

**Response Body:**

```json
{
  "id": 3,
  "name": "Updated Product",
  "description": "Updated description",
  "price": 249.99,
  "stock": 75,
  "category_id": 1,
  "seller_id": 2,
  "images": ["updated-product.jpg"],
  "has_variants": false,
  "created_at": "2023-04-20T14:00:00Z",
  "updated_at": "2023-04-20T14:30:00Z"
}
```

**Status Codes:**

- `200 OK`: Product updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product not found

### Delete Product (Seller Only)

```plaintext
DELETE /api/products/{id}
```

**Status Codes:**

- `204 No Content`: Product deleted successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product not found

### List Seller Products (Seller Only)

```plaintext
GET /api/products/seller
```

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

**Response Body:**

```json
[
  {
    "id": 3,
    "name": "Updated Product",
    "description": "Updated description",
    "price": 249.99,
    "stock": 75,
    "category_id": 1,
    "seller_id": 2,
    "images": ["updated-product.jpg"],
    "has_variants": false,
    "created_at": "2023-04-20T14:00:00Z",
    "updated_at": "2023-04-20T14:30:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Products retrieved successfully
- `401 Unauthorized`: Not authenticated

### Add Product Variant (Seller Only)

```plaintext
POST /api/products/{productId}/variants
```

**Request Body:**

```json
{
  "sku": "PROD-RED-M",
  "price": 29.99,
  "compare_price": 39.99,
  "stock": 10,
  "attributes": [
    { "name": "Color", "value": "Red" },
    { "name": "Size", "value": "Medium" }
  ],
  "images": ["red-shirt.jpg"],
  "is_default": true
}
```

**Response Body:**

```json
{
  "id": 1,
  "product_id": 2,
  "sku": "PROD-RED-M",
  "price": 29.99,
  "compare_price": 39.99,
  "stock": 10,
  "attributes": [
    { "name": "Color", "value": "Red" },
    { "name": "Size", "value": "Medium" }
  ],
  "images": ["red-shirt.jpg"],
  "is_default": true,
  "created_at": "2023-04-20T15:00:00Z",
  "updated_at": "2023-04-20T15:00:00Z"
}
```

**Status Codes:**

- `201 Created`: Variant created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product not found

### Update Product Variant (Seller Only)

```plaintext
PUT /api/products/{productId}/variants/{variantId}
```

**Request Body:**

```json
{
  "sku": "PROD-RED-M",
  "price": 24.99,
  "compare_price": 34.99,
  "stock": 15,
  "attributes": [
    { "name": "Color", "value": "Red" },
    { "name": "Size", "value": "Medium" }
  ],
  "images": ["red-shirt-updated.jpg"],
  "is_default": true
}
```

**Response Body:**

```json
{
  "id": 1,
  "product_id": 2,
  "sku": "PROD-RED-M",
  "price": 24.99,
  "compare_price": 34.99,
  "stock": 15,
  "attributes": [
    { "name": "Color", "value": "Red" },
    { "name": "Size", "value": "Medium" }
  ],
  "images": ["red-shirt-updated.jpg"],
  "is_default": true,
  "created_at": "2023-04-20T15:00:00Z",
  "updated_at": "2023-04-20T15:30:00Z"
}
```

**Status Codes:**

- `200 OK`: Variant updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product or variant not found

### Delete Product Variant (Seller Only)

```plaintext
DELETE /api/products/{productId}/variants/{variantId}
```

**Status Codes:**

- `204 No Content`: Variant deleted successfully
- `400 Bad Request`: Cannot delete the only variant of a product
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product or variant not found

## Categories

### List Categories

```plaintext
GET /api/categories
```

**Response Body:**

```json
[
  {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and accessories",
    "parent_id": null,
    "created_at": "2023-04-10T10:00:00Z",
    "updated_at": "2023-04-10T10:00:00Z"
  },
  {
    "id": 2,
    "name": "Smartphones",
    "description": "Mobile phones and accessories",
    "parent_id": 1,
    "created_at": "2023-04-10T10:05:00Z",
    "updated_at": "2023-04-10T10:05:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Categories retrieved successfully

## Cart

### Get Cart

```plaintext
GET /api/cart
```

**Response Body:**

```json
{
  "id": 1,
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "cart_id": 1,
      "product_id": 1,
      "quantity": 2,
      "created_at": "2023-04-20T16:00:00Z",
      "updated_at": "2023-04-20T16:00:00Z"
    }
  ],
  "created_at": "2023-04-20T15:45:00Z",
  "updated_at": "2023-04-20T16:00:00Z"
}
```

**Status Codes:**

- `200 OK`: Cart retrieved successfully
- `401 Unauthorized`: Not authenticated
- `404 Not Found`: Cart not found

### Add to Cart

```plaintext
POST /api/cart/items
```

**Request Body:**

```json
{
  "product_id": 2,
  "variant_id": 5,
  "quantity": 1
}
```

**Response Body:**

```json
{
  "id": 1,
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "cart_id": 1,
      "product_id": 1,
      "product_variant_id": 3,
      "quantity": 2,
      "created_at": "2023-04-20T16:00:00Z",
      "updated_at": "2023-04-20T16:00:00Z"
    },
    {
      "id": 2,
      "cart_id": 1,
      "product_id": 2,
      "product_variant_id": 5,
      "quantity": 1,
      "created_at": "2023-04-20T16:15:00Z",
      "updated_at": "2023-04-20T16:15:00Z"
    }
  ],
  "created_at": "2023-04-20T15:00:00Z",
  "updated_at": "2023-04-20T16:15:00Z"
}
```

**Status Codes:**

- `200 OK`: Item added to cart successfully
- `400 Bad Request`: Invalid request body, product not found, or insufficient stock
- `401 Unauthorized`: Not authenticated (for user cart operations)

### Update Cart Item

```plaintext
PUT /api/cart/items/{productId}
```

**Request Body:**

```json
{
  "quantity": 3,
  "variant_id": 5
}
```

**Response Body:**

```json
{
  "id": 1,
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "cart_id": 1,
      "product_id": 1,
      "product_variant_id": 3,
      "quantity": 2,
      "created_at": "2023-04-20T16:00:00Z",
      "updated_at": "2023-04-20T16:00:00Z"
    },
    {
      "id": 2,
      "cart_id": 1,
      "product_id": 2,
      "product_variant_id": 5,
      "quantity": 3,
      "created_at": "2023-04-20T16:15:00Z",
      "updated_at": "2023-04-20T16:20:00Z"
    }
  ],
  "created_at": "2023-04-20T15:00:00Z",
  "updated_at": "2023-04-20T16:20:00Z"
}
```

**Status Codes:**

- `200 OK`: Cart item updated successfully
- `400 Bad Request`: Invalid request body, product not found, or insufficient stock
- `401 Unauthorized`: Not authenticated (for user cart operations)

### Remove from Cart

```plaintext
DELETE /api/cart/items/{productId}?variantId={variantId}
```

**Response Body:**

```json
{
  "id": 1,
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "cart_id": 1,
      "product_id": 1,
      "product_variant_id": 3,
      "quantity": 2,
      "created_at": "2023-04-20T16:00:00Z",
      "updated_at": "2023-04-20T16:00:00Z"
    }
  ],
  "created_at": "2023-04-20T15:00:00Z",
  "updated_at": "2023-04-20T16:25:00Z"
}
```

**Status Codes:**

- `200 OK`: Cart item removed successfully
- `400 Bad Request`: Product not found in cart
- `401 Unauthorized`: Not authenticated (for user cart operations)

### Clear Cart

```plaintext
DELETE /api/cart
```

**Response Body:**

```json
{
  "id": 1,
  "user_id": 1,
  "items": [],
  "created_at": "2023-04-20T15:45:00Z",
  "updated_at": "2023-04-20T17:00:00Z"
}
```

**Status Codes:**

- `200 OK`: Cart cleared successfully
- `401 Unauthorized`: Not authenticated

```plaintext
POST /api/orders
```

**Request Body:**

```json
{
  "shipping_addr": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  },
  "billing_addr": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  }
}
```

**Response Body:**

```json
{
  "id": 1,
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 2,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99
    }
  ],
  "total_amount": 1499.99,
  "status": "pending",
  "shipping_address": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  },
  "billing_address": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  },
  "payment_id": "",
  "payment_provider": "",
  "tracking_code": "",
  "created_at": "2023-04-20T17:30:00Z",
  "updated_at": "2023-04-20T17:30:00Z",
  "completed_at": null
}
```

**Status Codes:**

- `201 Created`: Order created successfully
- `400 Bad Request`: Invalid request body, cart is empty, or insufficient stock
- `401 Unauthorized`: Not authenticated

### Get Order

```plaintext
GET /api/orders/{id}
```

**Response Body:**

```json
{
  "id": 1,
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 2,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99
    }
  ],
  "total_amount": 1499.99,
  "status": "pending",
  "shipping_address": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  },
  "billing_address": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  },
  "payment_id": "",
  "payment_provider": "",
  "tracking_code": "",
  "created_at": "2023-04-20T17:30:00Z",
  "updated_at": "2023-04-20T17:30:00Z",
  "completed_at": null
}
```

**Status Codes:**

- `200 OK`: Order retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the owner of this order)
- `404 Not Found`: Order not found

### List Orders

```plaintext
GET /api/orders
```

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

**Response Body:**

```json
[
  {
    "id": 1,
    "user_id": 1,
    "items": [
      {
        "id": 1,
        "order_id": 1,
        "product_id": 2,
        "quantity": 1,
        "price": 1499.99,
        "subtotal": 1499.99
      }
    ],
    "total_amount": 1499.99,
    "status": "pending",
    "shipping_address": {
      "street": "123 Main St",
      "city": "Anytown",
      "state": "CA",
      "postal_code": "12345",
      "country": "USA"
    },
    "billing_address": {
      "street": "123 Main St",
      "city": "Anytown",
      "state": "CA",
      "postal_code": "12345",
      "country": "USA"
    },
    "payment_id": "",
    "payment_provider": "",
    "tracking_code": "",
    "created_at": "2023-04-20T17:30:00Z",
    "updated_at": "2023-04-20T17:30:00Z",
    "completed_at": null
  }
]
```

**Status Codes:**

- `200 OK`: Orders retrieved successfully
- `401 Unauthorized`: Not authenticated

### Process Payment

```plaintext
POST /api/orders/{id}/payment
```

**Request Body:**

```json
{
  "payment_method": "credit_card",
  "payment_provider": "stripe",
  "card_details": {
    "card_number": "4242424242424242",
    "expiry_month": 12,
    "expiry_year": 2025,
    "cvv": "123",
    "cardholder_name": "John Doe",
    "token": "tok_visa"
  },
  "customer_email": "user@example.com"
}
```

**Response Body:**

```json
{
  "id": 1,
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 2,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99
    }
  ],
  "total_amount": 1499.99,
  "status": "paid",
  "shipping_address": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  },
  "billing_address": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  },
  "payment_id": "pi_3MkCrjKZ6o8QJAcJ0KjkLNZt",
  "payment_provider": "stripe",
  "tracking_code": "",
  "created_at": "2023-04-20T17:30:00Z",
  "updated_at": "2023-04-20T18:00:00Z",
  "completed_at": null
}
```

**Status Codes:**

- `200 OK`: Payment processed successfully
- `400 Bad Request`: Invalid request body, payment failed, or order already paid
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the owner of this order)
- `404 Not Found`: Order not found

### List All Orders (Admin Only)

```plaintext
GET /api/admin/orders
```

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)
- `status` (optional): Filter by order status

**Response Body:**

```json
[
  {
    "id": 1,
    "user_id": 1,
    "items": [
      {
        "id": 1,
        "order_id": 1,
        "product_id": 2,
        "quantity": 1,
        "price": 1499.99,
        "subtotal": 1499.99
      }
    ],
    "total_amount": 1499.99,
    "status": "paid",
    "shipping_address": {
      "street": "123 Main St",
      "city": "Anytown",
      "state": "CA",
      "postal_code": "12345",
      "country": "USA"
    },
    "billing_address": {
      "street": "123 Main St",
      "city": "Anytown",
      "state": "CA",
      "postal_code": "12345",
      "country": "USA"
    },
    "payment_id": "pi_3MkCrjKZ6o8QJAcJ0KjkLNZt",
    "payment_provider": "stripe",
    "tracking_code": "",
    "created_at": "2023-04-20T17:30:00Z",
    "updated_at": "2023-04-20T18:00:00Z",
    "completed_at": null
  }
]
```

**Status Codes:**

- `200 OK`: Orders retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

### Update Order Status (Admin Only)

```plaintext
PUT /api/admin/orders/{id}/status
```

**Request Body:**

```json
{
  "status": "shipped"
}
```

**Response Body:**

```json
{
  "id": 1,
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "order_id": 1,
      "product_id": 2,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99
    }
  ],
  "total_amount": 1499.99,
  "status": "shipped",
  "shipping_address": {  1499.99
    }
  ],
  "total_amount": 1499.99,
  "status": "shipped",
  "shipping_address": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  },
  "billing_address": {
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "postal_code": "12345",
    "country": "USA"
  },
  "payment_id": "pi_3MkCrjKZ6o8QJAcJ0KjkLNZt",
  "payment_provider": "stripe",
  "tracking_code": "TRACK123456",
  "created_at": "2023-04-20T17:30:00Z",
  "updated_at": "2023-04-20T19:00:00Z",
  "completed_at": null
}
```

**Status Codes:**

- `200 OK`: Order status updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Order not found

## Payment

### Get Available Payment Providers

```plaintext
GET /api/payment/providers
```

**Response Body:**

```json
[
  {
    "type": "stripe",
    "name": "Stripe",
    "description": "Pay with credit or debit card",
    "icon_url": "/assets/images/stripe-logo.png",
    "methods": ["credit_card"],
    "enabled": true
  },
  {
    "type": "paypal",
    "name": "PayPal",
    "description": "Pay with your PayPal account",
    "icon_url": "/assets/images/paypal-logo.png",
    "methods": ["paypal"],
    "enabled": true
  },
  {
    "type": "mock",
    "name": "Test Payment",
    "description": "For testing purposes only",
    "methods": ["credit_card", "paypal", "bank_transfer"],
    "enabled": true
  }
]
```

**Status Codes:**

- `200 OK`: Payment providers retrieved successfully

## Webhooks

### Stripe Webhook

```plaintext
POST /api/webhooks/stripe
```

**Request Headers:**

- `Stripe-Signature`: Webhook signature from Stripe

**Request Body:**

- Stripe event object (varies based on event type)

**Response Body:**

```json
{
  "status": "success"
}
```

**Status Codes:**

- `200 OK`: Webhook processed successfully
- `400 Bad Request`: Invalid webhook signature or event

This completes the comprehensive API documentation for the e-commerce system, including all endpoints with their request and response bodies.
