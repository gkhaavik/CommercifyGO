# Checkout API Examples

This document provides examples of using the Checkout API for Commercify.

## Public Checkout Endpoints

### Get Guest Checkout

```plaintext
GET /api/guest/checkout
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 3,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:30:22Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "shipping_method_id": 2,
  "shipping_method": {
    "id": 2,
    "name": "Express Shipping",
    "description": "Delivered within 2-3 business days",
    "cost": 14.99
  },
  "payment_provider": "stripe",
  "total_amount": 49.98,
  "shipping_cost": 14.99,
  "total_weight": 0.6,
  "customer_details": {
    "email": "customer@example.com",
    "phone": "+1-555-123-4567",
    "full_name": "John Doe"
  },
  "currency": "USD",
  "discount_code": "SUMMER2025",
  "discount_amount": 5.00,
  "final_amount": 59.97,
  "applied_discount": {
    "id": 3,
    "code": "SUMMER2025",
    "type": "basket",
    "method": "fixed",
    "value": 5.00,
    "amount": 5.00
  },
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:45:18Z",
  "last_activity_at": "2025-05-22T15:45:18Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Checkout retrieved successfully
- `404 Not Found`: Checkout not found

### Add Item to Guest Checkout

```plaintext
POST /api/guest/checkout/items
```

**Request Body:**
```json
{
  "product_id": 12,
  "variant_id": 3,
  "quantity": 2
}
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 3,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:30:22Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "total_amount": 49.98,
  "shipping_cost": 0,
  "total_weight": 0.6,
  "customer_details": {
    "email": "",
    "phone": "",
    "full_name": ""
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 49.98,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:30:22Z",
  "last_activity_at": "2025-05-22T15:30:22Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Item added successfully
- `400 Bad Request`: Invalid request body
- `404 Not Found`: Product or variant not found

### Update Item in Guest Checkout

```plaintext
PUT /api/guest/checkout/items/{productId}
```

**Request Body:**
```json
{
  "quantity": 3,
  "variant_id": 4
}
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 4,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, L",
      "sku": "TS-BLK-L",
      "price": 24.99,
      "quantity": 3,
      "weight": 0.3,
      "subtotal": 74.97,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:35:18Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "total_amount": 74.97,
  "shipping_cost": 0,
  "total_weight": 0.9,
  "customer_details": {
    "email": "",
    "phone": "",
    "full_name": ""
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 74.97,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:35:18Z",
  "last_activity_at": "2025-05-22T15:35:18Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Item updated successfully
- `400 Bad Request`: Invalid request body
- `404 Not Found`: Item, product, or variant not found

### Remove Item from Guest Checkout

```plaintext
DELETE /api/guest/checkout/items/{productId}
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [],
  "status": "active",
  "shipping_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "total_amount": 0,
  "shipping_cost": 0,
  "total_weight": 0,
  "customer_details": {
    "email": "",
    "phone": "",
    "full_name": ""
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 0,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:40:10Z",
  "last_activity_at": "2025-05-22T15:40:10Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Item removed successfully
- `404 Not Found`: Item not found

### Clear Guest Checkout

```plaintext
DELETE /api/guest/checkout
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [],
  "status": "active",
  "shipping_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "total_amount": 0,
  "shipping_cost": 0,
  "total_weight": 0,
  "customer_details": {
    "email": "",
    "phone": "",
    "full_name": ""
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 0,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:41:15Z",
  "last_activity_at": "2025-05-22T15:41:15Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Checkout cleared successfully

### Set Shipping Address for Guest Checkout

```plaintext
PUT /api/guest/checkout/shipping-address
```

**Request Body:**
```json
{
  "address_line1": "123 Main St",
  "address_line2": "Apt 4B",
  "city": "San Francisco",
  "state": "CA",
  "postal_code": "94105",
  "country": "US"
}
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 3,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:30:22Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "",
    "address_line2": "",
    "city": "",
    "state": "",
    "postal_code": "",
    "country": ""
  },
  "total_amount": 49.98,
  "shipping_cost": 0,
  "total_weight": 0.6,
  "customer_details": {
    "email": "",
    "phone": "",
    "full_name": ""
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 49.98,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:42:30Z",
  "last_activity_at": "2025-05-22T15:42:30Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Shipping address set successfully
- `400 Bad Request`: Invalid address data

### Set Billing Address for Guest Checkout

```plaintext
PUT /api/guest/checkout/billing-address
```

**Request Body:**
```json
{
  "address_line1": "123 Main St",
  "address_line2": "Apt 4B",
  "city": "San Francisco",
  "state": "CA",
  "postal_code": "94105",
  "country": "US"
}
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 3,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:30:22Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "total_amount": 49.98,
  "shipping_cost": 0,
  "total_weight": 0.6,
  "customer_details": {
    "email": "",
    "phone": "",
    "full_name": ""
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 49.98,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:43:15Z",
  "last_activity_at": "2025-05-22T15:43:15Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Billing address set successfully
- `400 Bad Request`: Invalid address data

### Set Customer Details for Guest Checkout

```plaintext
PUT /api/guest/checkout/customer-details
```

**Request Body:**
```json
{
  "email": "customer@example.com",
  "phone": "+1-555-123-4567",
  "full_name": "John Doe"
}
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 3,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:30:22Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "total_amount": 49.98,
  "shipping_cost": 0,
  "total_weight": 0.6,
  "customer_details": {
    "email": "customer@example.com",
    "phone": "+1-555-123-4567",
    "full_name": "John Doe"
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 49.98,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:44:05Z",
  "last_activity_at": "2025-05-22T15:44:05Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Customer details set successfully
- `400 Bad Request`: Invalid customer data

### Set Shipping Method for Guest Checkout

```plaintext
PUT /api/guest/checkout/shipping-method
```

**Request Body:**
```json
{
  "shipping_method_id": 2
}
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 3,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:30:22Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "shipping_method_id": 2,
  "shipping_method": {
    "id": 2,
    "name": "Express Shipping",
    "description": "Delivered within 2-3 business days",
    "cost": 14.99
  },
  "total_amount": 49.98,
  "shipping_cost": 14.99,
  "total_weight": 0.6,
  "customer_details": {
    "email": "customer@example.com",
    "phone": "+1-555-123-4567",
    "full_name": "John Doe"
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 64.97,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:44:45Z",
  "last_activity_at": "2025-05-22T15:44:45Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Shipping method set successfully
- `400 Bad Request`: Invalid shipping method ID
- `404 Not Found`: Shipping method not found

### Apply Discount to Guest Checkout

```plaintext
POST /api/guest/checkout/discount
```

**Request Body:**
```json
{
  "discount_code": "SUMMER2025"
}
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 3,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:30:22Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "shipping_method_id": 2,
  "shipping_method": {
    "id": 2,
    "name": "Express Shipping",
    "description": "Delivered within 2-3 business days",
    "cost": 14.99
  },
  "total_amount": 49.98,
  "shipping_cost": 14.99,
  "total_weight": 0.6,
  "customer_details": {
    "email": "customer@example.com",
    "phone": "+1-555-123-4567",
    "full_name": "John Doe"
  },
  "currency": "USD",
  "discount_code": "SUMMER2025",
  "discount_amount": 5.00,
  "final_amount": 59.97,
  "applied_discount": {
    "id": 3,
    "code": "SUMMER2025",
    "type": "basket",
    "method": "fixed",
    "value": 5.00,
    "amount": 5.00
  },
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:45:18Z",
  "last_activity_at": "2025-05-22T15:45:18Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Discount applied successfully
- `400 Bad Request`: Invalid discount code
- `404 Not Found`: Discount not found
- `409 Conflict`: Discount cannot be applied (e.g., minimum order value not met)

### Remove Discount from Guest Checkout

```plaintext
DELETE /api/guest/checkout/discount
```

**Response Body:**
```json
{
  "id": 5,
  "session_id": "7f98a2c4-8f9d-4c3e-9f9a-3f4f5c6d7e8a",
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 3,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:30:22Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "shipping_method_id": 2,
  "shipping_method": {
    "id": 2,
    "name": "Express Shipping",
    "description": "Delivered within 2-3 business days",
    "cost": 14.99
  },
  "total_amount": 49.98,
  "shipping_cost": 14.99,
  "total_weight": 0.6,
  "customer_details": {
    "email": "customer@example.com",
    "phone": "+1-555-123-4567",
    "full_name": "John Doe"
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 64.97,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:46:05Z",
  "last_activity_at": "2025-05-22T15:46:05Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Discount removed successfully
- `400 Bad Request`: No discount applied

### Convert Guest Checkout to Order

```plaintext
POST /api/guest/checkout/to-order
```

**Response Body:**
```json
{
  "id": 123,
  "user_id": 0,
  "order_number": "ORD-20250522-000123",
  "items": [
    {
      "id": 187,
      "product_id": 12,
      "variant_id": 3,
      "name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "subtotal": 49.98
    }
  ],
  "status": "pending",
  "total_amount": 49.98,
  "final_amount": 64.97,
  "currency": "USD",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "payment_details": {
    "provider": "",
    "method": "",
    "id": "",
    "status": "pending",
    "captured": false,
    "refunded": false
  },
  "shipping_details": {
    "method_id": 2,
    "method": "Express Shipping",
    "cost": 14.99
  },
  "customer": {
    "email": "customer@example.com",
    "phone": "+1-555-123-4567",
    "full_name": "John Doe"
  },
  "created_at": "2025-05-22T15:47:10Z",
  "updated_at": "2025-05-22T15:47:10Z"
}
```

**Status Codes:**
- `201 Created`: Order created successfully
- `400 Bad Request`: Missing required fields (customer details, shipping address, etc.)
- `409 Conflict`: No items in checkout or other validation errors

## Authenticated Checkout Endpoints

All endpoints above are also available for authenticated users by replacing `/guest/checkout` with `/checkout`. Additionally:

### Convert Guest Checkout to User Checkout

```plaintext
POST /api/checkout/convert
```

**Response Body:**
```json
{
  "id": 5,
  "user_id": 42,
  "items": [
    {
      "id": 7,
      "product_id": 12,
      "variant_id": 3,
      "product_name": "Premium T-shirt",
      "variant_name": "Black, XL",
      "sku": "TS-BLK-XL",
      "price": 24.99,
      "quantity": 2,
      "weight": 0.3,
      "subtotal": 49.98,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:30:22Z"
    }
  ],
  "status": "active",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "shipping_method_id": 2,
  "shipping_method": {
    "id": 2,
    "name": "Express Shipping",
    "description": "Delivered within 2-3 business days",
    "cost": 14.99
  },
  "total_amount": 49.98,
  "shipping_cost": 14.99,
  "total_weight": 0.6,
  "customer_details": {
    "email": "customer@example.com",
    "phone": "+1-555-123-4567",
    "full_name": "John Doe"
  },
  "currency": "USD",
  "discount_amount": 0,
  "final_amount": 64.97,
  "created_at": "2025-05-22T15:30:22Z",
  "updated_at": "2025-05-22T15:50:20Z",
  "last_activity_at": "2025-05-22T15:50:20Z",
  "expires_at": "2025-05-23T15:30:22Z"
}
```

**Status Codes:**
- `200 OK`: Checkout converted successfully
- `401 Unauthorized`: User not authenticated
- `404 Not Found`: No guest checkout found for this session

## Admin Checkout Endpoints

### List Checkouts (Admin Only)

```plaintext
GET /api/admin/checkouts
```

**Query Parameters:**
- `page` (optional): Page number (defaults to 1)
- `page_size` (optional): Number of items per page (defaults to 10)
- `status` (optional): Filter by status (active, completed, abandoned, expired)

**Response Body:**
```json
{
  "success": true,
  "data": [
    {
      "id": 5,
      "user_id": 42,
      "items": [
        {
          "id": 7,
          "product_id": 12,
          "variant_id": 3,
          "product_name": "Premium T-shirt",
          "variant_name": "Black, XL",
          "sku": "TS-BLK-XL",
          "price": 24.99,
          "quantity": 2,
          "weight": 0.3,
          "subtotal": 49.98,
          "created_at": "2025-05-22T15:30:22Z",
          "updated_at": "2025-05-22T15:30:22Z"
        }
      ],
      "status": "active",
      "shipping_address": {
        "address_line1": "123 Main St",
        "address_line2": "Apt 4B",
        "city": "San Francisco",
        "state": "CA",
        "postal_code": "94105",
        "country": "US"
      },
      "billing_address": {
        "address_line1": "123 Main St",
        "address_line2": "Apt 4B",
        "city": "San Francisco",
        "state": "CA",
        "postal_code": "94105",
        "country": "US"
      },
      "shipping_method_id": 2,
      "shipping_method": {
        "id": 2,
        "name": "Express Shipping",
        "description": "Delivered within 2-3 business days",
        "cost": 14.99
      },
      "total_amount": 49.98,
      "shipping_cost": 14.99,
      "total_weight": 0.6,
      "customer_details": {
        "email": "customer@example.com",
        "phone": "+1-555-123-4567",
        "full_name": "John Doe"
      },
      "currency": "USD",
      "discount_amount": 0,
      "final_amount": 64.97,
      "created_at": "2025-05-22T15:30:22Z",
      "updated_at": "2025-05-22T15:50:20Z",
      "last_activity_at": "2025-05-22T15:50:20Z",
      "expires_at": "2025-05-23T15:30:22Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 1
  }
}
```

**Status Codes:**
- `200 OK`: Checkouts retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)

## Example Workflow

### Guest Checkout Flow

1. Guest adds items to their checkout
   - `POST /api/guest/checkout/items`
2. Guest provides shipping address
   - `PUT /api/guest/checkout/shipping-address`
3. Guest provides billing address
   - `PUT /api/guest/checkout/billing-address`
4. Guest provides customer details
   - `PUT /api/guest/checkout/customer-details`
5. Guest selects shipping method
   - `PUT /api/guest/checkout/shipping-method`
6. Guest applies discount code (optional)
   - `POST /api/guest/checkout/discount`
7. Guest converts checkout to order
   - `POST /api/guest/checkout/to-order`
8. Guest processes payment for the order
   - `POST /api/guest/orders/{id}/payment`

### Authenticated User Checkout Flow

1. User logs in
   - `POST /api/auth/signin`
2. User converts guest checkout to user checkout (optional)
   - `POST /api/checkout/convert`
3. User adds items to their checkout
   - `POST /api/checkout/items`
4. User provides shipping address
   - `PUT /api/checkout/shipping-address`
5. User provides billing address
   - `PUT /api/checkout/billing-address`
6. User provides customer details
   - `PUT /api/checkout/customer-details`
7. User selects shipping method
   - `PUT /api/checkout/shipping-method`
8. User applies discount code (optional)
   - `POST /api/checkout/discount`
9. User converts checkout to order
   - `POST /api/checkout/to-order`
10. User processes payment for the order
    - `POST /api/orders/{id}/payment`
