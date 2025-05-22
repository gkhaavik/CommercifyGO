# Order API Examples

This document provides example request bodies for the order system API endpoints.

## Public Order Endpoints

### Create Guest Order

```plaintext
POST /api/guest/orders
```

Create an order as a guest user.

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Smith",
  "phone_number": "+1234567890",
  "shipping_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "USA"
  },
  "billing_address": {
    "address_line1": "123 Main St",
    "address_line2": "Apt 4B",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "USA"
  },
  "shipping_method_id": 1
}
```

Example response:

```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "order_number": "ORD_xxxxx",
    "session_id": "",
    "status": "pending",
    "total_amount": 227.96,
    "currency": "USD",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "order_id": "550e8400-e29b-41d4-a716-446655440000",
        "product_id": "550e8400-e29b-41d4-a716-446655440002",
        "name": "Product Name",
        "sku": "PROD-001",
        "quantity": 1,
        "unit_price": 199.99,
        "total_price": 199.99
      }
    ],
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
    "payment_method": "credit_card",
    "payment_status": "pending",
    "shipping_method": "standard",
    "shipping_cost": 7.99,
    "discount_amount": 0,
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
  }
}
```

**Status Codes:**

- `201 Created`: Order created successfully
- `400 Bad Request`: Invalid request body or empty cart
- `500 Internal Server Error`: Failed to create order

### Process Payment for Guest Order

```plaintext
POST /api/guest/orders/{id}/payment
```

Process payment for a guest order.

**Request Body:**

```json
{
  "payment_method": "credit_card",
  "payment_provider": "stripe",
  "card_details": {
    "card_number": "4242424242424242",
    "expiry_month": 12,
    "expiry_year": 2025,
    "cvc": "123",
    "card_holder_name": "John Smith"
  },
  "phone_number": "+1234567890"
}
```

Example response:

```json
{
  "success": true,
  "message": "Payment processed successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "00000000-0000-0000-0000-000000000000",
    "status": "paid",
    "total_amount": 227.96,
    "currency": "USD",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "order_id": "550e8400-e29b-41d4-a716-446655440000",
        "product_id": "550e8400-e29b-41d4-a716-446655440002",
        "name": "Product Name",
        "sku": "PROD-001",
        "quantity": 1,
        "unit_price": 199.99,
        "total_price": 199.99
      }
    ],
    "shipping_address": {
      "first_name": "John",
      "last_name": "Smith",
      "address_line1": "123 Main St",
      "address_line2": "Apt 4B",
      "city": "San Francisco",
      "state": "CA",
      "postal_code": "94105",
      "country": "US",
      "phone_number": "+1234567890"
    },
    "billing_address": {
      "first_name": "John",
      "last_name": "Smith",
      "address_line1": "123 Main St",
      "address_line2": "Apt 4B",
      "city": "San Francisco",
      "state": "CA",
      "postal_code": "94105",
      "country": "US",
      "phone_number": "+1234567890"
    },
    "payment_method": "credit_card",
    "payment_status": "paid",
    "shipping_method": "standard",
    "shipping_cost": 7.99,
    "tax_amount": 19.98,
    "discount_amount": 0,
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:05:00Z"
  }
}
```

**Status Codes:**

- `200 OK`: Payment processed successfully
- `400 Bad Request`: Invalid payment details or order already paid
- `401 Unauthorized`: Invalid session for guest order
- `404 Not Found`: Order not found
- `500 Internal Server Error`: Payment processing failed

## Authenticated Order Endpoints

### Create Order for Authenticated User

```plaintext
POST /api/orders
```

Create an order for the authenticated user.

**Request Body:**

```json
{
  "shipping_address": {
    "first_name": "Sarah",
    "last_name": "Johnson",
    "address_line1": "456 Oak Avenue",
    "address_line2": "Suite 100",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US",
    "phone_number": "+1987654321"
  },
  "billing_address": {
    "first_name": "Sarah",
    "last_name": "Johnson",
    "address_line1": "456 Oak Avenue",
    "address_line2": "Suite 100",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US",
    "phone_number": "+1987654321"
  },
  "shipping_method": "express",
  "payment_method": "wallet"
}
```

Example response:

```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440004",
    "status": "pending",
    "total_amount": 2514.97,
    "currency": "USD",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440005",
        "order_id": "550e8400-e29b-41d4-a716-446655440003",
        "product_id": "550e8400-e29b-41d4-a716-446655440006",
        "name": "Premium Product",
        "sku": "PROD-002",
        "quantity": 1,
        "unit_price": 2499.99,
        "total_price": 2499.99
      }
    ],
    "shipping_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "billing_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "payment_method": "wallet",
    "payment_status": "pending",
    "shipping_method": "express",
    "shipping_cost": 14.99,
    "tax_amount": 0,
    "discount_amount": 0,
    "created_at": "2024-03-20T11:00:00Z",
    "updated_at": "2024-03-20T11:00:00Z"
  }
}
```

**Status Codes:**

- `201 Created`: Order created successfully
- `400 Bad Request`: Invalid request body or empty cart
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to create order

### Process Payment for User Order

```plaintext
POST /api/orders/{id}/payment
```

Process payment for an authenticated user's order.

**Request Body:**

```json
{
  "payment_method": "wallet",
  "payment_provider": "mobilepay",
  "phone_number": "+1987654321"
}
```

Example response:

```json
{
  "success": true,
  "message": "Payment processed successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440004",
    "status": "paid",
    "total_amount": 2514.97,
    "currency": "USD",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440005",
        "order_id": "550e8400-e29b-41d4-a716-446655440003",
        "product_id": "550e8400-e29b-41d4-a716-446655440006",
        "name": "Premium Product",
        "sku": "PROD-002",
        "quantity": 1,
        "unit_price": 2499.99,
        "total_price": 2499.99
      }
    ],
    "shipping_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "billing_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "payment_method": "wallet",
    "payment_status": "paid",
    "shipping_method": "express",
    "shipping_cost": 14.99,
    "tax_amount": 0,
    "discount_amount": 0,
    "created_at": "2024-03-20T11:00:00Z",
    "updated_at": "2024-03-20T11:05:00Z"
  }
}
```

**Status Codes:**

- `200 OK`: Payment processed successfully
- `400 Bad Request`: Invalid payment details or order already paid
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized for this order
- `404 Not Found`: Order not found
- `500 Internal Server Error`: Payment processing failed

### Get Order

```plaintext
GET /api/orders/{id}
```

Retrieve a specific order for the authenticated user.

Example response:

```json
{
  "success": true,
  "message": "Order retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440004",
    "status": "paid",
    "total_amount": 2514.97,
    "currency": "USD",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440005",
        "order_id": "550e8400-e29b-41d4-a716-446655440003",
        "product_id": "550e8400-e29b-41d4-a716-446655440006",
        "name": "Premium Product",
        "sku": "PROD-002",
        "quantity": 1,
        "unit_price": 2499.99,
        "total_price": 2499.99
      }
    ],
    "shipping_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "billing_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "payment_method": "wallet",
    "payment_status": "paid",
    "shipping_method": "express",
    "shipping_cost": 14.99,
    "tax_amount": 0,
    "discount_amount": 0,
    "created_at": "2024-03-20T11:00:00Z",
    "updated_at": "2024-03-20T11:05:00Z"
  }
}
```

**Status Codes:**

- `200 OK`: Order retrieved successfully
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized for this order
- `404 Not Found`: Order not found
- `500 Internal Server Error`: Failed to retrieve order

### List User Orders

```plaintext
GET /api/orders
```

List all orders for the authenticated user.

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

Example response:

```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "user_id": "550e8400-e29b-41d4-a716-446655440004",
      "status": "paid",
      "total_amount": 2514.97,
      "currency": "USD",
      "payment_method": "wallet",
      "payment_status": "paid",
      "shipping_method": "express",
      "shipping_cost": 14.99,
      "tax_amount": 0,
      "discount_amount": 0,
      "created_at": "2024-03-20T11:00:00Z",
      "updated_at": "2024-03-20T11:05:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "offset": 0,
    "limit": 10
  }
}
```

**Status Codes:**

- `200 OK`: Orders retrieved successfully
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to retrieve orders

## Admin Order Endpoints

### List All Orders

```plaintext
GET /api/admin/orders
```

List all orders in the system (admin only).

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)
- `status` (optional): Filter by order status

Example response:

```json
{
  "success": true,
  "message": "Orders retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "user_id": "550e8400-e29b-41d4-a716-446655440004",
      "status": "paid",
      "total_amount": 2514.97,
      "currency": "USD",
      "payment_method": "wallet",
      "payment_status": "paid",
      "shipping_method": "express",
      "shipping_cost": 14.99,
      "tax_amount": 0,
      "discount_amount": 0,
      "created_at": "2024-03-20T11:00:00Z",
      "updated_at": "2024-03-20T11:05:00Z"
    }
  ],
  "pagination": {
    "total": 1,
    "offset": 0,
    "limit": 10
  }
}
```

**Status Codes:**

- `200 OK`: Orders retrieved successfully
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `500 Internal Server Error`: Failed to retrieve orders

### Update Order Status

```plaintext
PUT /api/admin/orders/{id}/status
```

Update an order's status (admin only).

**Request Body:**

```json
{
  "status": "shipped"
}
```

Example response:

```json
{
  "success": true,
  "message": "Order status updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440004",
    "status": "shipped",
    "total_amount": 2514.97,
    "currency": "USD",
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440005",
        "order_id": "550e8400-e29b-41d4-a716-446655440003",
        "product_id": "550e8400-e29b-41d4-a716-446655440006",
        "name": "Premium Product",
        "sku": "PROD-002",
        "quantity": 1,
        "unit_price": 2499.99,
        "total_price": 2499.99
      }
    ],
    "shipping_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "billing_address": {
      "first_name": "Sarah",
      "last_name": "Johnson",
      "address_line1": "456 Oak Avenue",
      "address_line2": "Suite 100",
      "city": "Seattle",
      "state": "WA",
      "postal_code": "98101",
      "country": "US",
      "phone_number": "+1987654321"
    },
    "payment_method": "wallet",
    "payment_status": "paid",
    "shipping_method": "express",
    "shipping_cost": 14.99,
    "tax_amount": 0,
    "discount_amount": 0,
    "created_at": "2024-03-20T11:00:00Z",
    "updated_at": "2024-03-20T14:30:00Z"
  }
}
```

**Status Codes:**

- `200 OK`: Order status updated successfully
- `400 Bad Request`: Invalid order status
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `404 Not Found`: Order not found
- `500 Internal Server Error`: Failed to update order status

## Example Workflow

### Guest Checkout Flow

1. Guest adds items to their cart
2. Guest provides shipping information and selects shipping method
3. System creates an order with `POST /api/guest/orders`
4. Guest provides payment details with `POST /api/guest/orders/{id}/payment`
5. Payment is processed and order status is updated to "paid"

### Authenticated User Checkout Flow

1. User adds items to their cart
2. User provides shipping information and selects shipping method
3. System creates an order with `POST /api/orders`
4. User provides payment details with `POST /api/orders/{id}/payment`
5. Payment is processed and order status is updated to "paid"

### Order Fulfillment Flow (Admin)

1. Admin views orders with `GET /api/admin/orders`
2. Admin processes the order (picking, packing)
3. Admin updates order status to "shipped" with `PUT /api/admin/orders/{id}/status`
4. System sends shipping confirmation email to customer
5. When delivery is confirmed, admin updates status to "delivered"
