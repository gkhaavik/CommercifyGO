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
  "email": "customer@example.com",
  "full_name": "John Smith",
  "phone_number": "+1234567890",
  "shipping_address": {
    "street_address": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "street_address": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US" 
  },
  "shipping_method_id": 3
}
```

Example response:

```json
{
  "id": 10,
  "order_number": "ORD-20230625-000010",
  "user_id": 0,
  "items": [
    {
      "id": 15,
      "order_id": 10,
      "product_id": 3,
      "quantity": 1,
      "price": 19.99,
      "subtotal": 19.99,
      "weight": 0.2
    },
    {
      "id": 16,
      "order_id": 10,
      "product_id": 5,
      "variant_id": 10,
      "quantity": 2,
      "price": 99.99,
      "subtotal": 199.98,
      "weight": 1.6
    }
  ],
  "subtotal": 219.97,
  "discount_code": null,
  "discount_amount": 0,
  "shipping_cost": 7.99,
  "final_amount": 227.96,
  "status": "pending",
  "shipping_address": {
    "street_address": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "street_address": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US" 
  },
  "email": "customer@example.com",
  "full_name": "John Smith",
  "phone_number": "+1234567890",
  "shipping_method_id": 3,
  "shipping_method_name": "Standard Shipping",
  "payment_id": null,
  "payment_provider": null,
  "is_guest_order": true,
  "requires_action": false,
  "action_url": "",
  "created_at": "2023-06-25T15:30:45Z",
  "updated_at": "2023-06-25T15:30:45Z"
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
  "customer_email": "customer@example.com"
}
```

Example response:

```json
{
  "id": 10,
  "order_number": "ORD-20230625-000010",
  "user_id": 0,
  "items": [
    {
      "id": 15,
      "order_id": 10,
      "product_id": 3,
      "quantity": 1,
      "price": 19.99,
      "subtotal": 19.99,
      "weight": 0.2
    },
    {
      "id": 16,
      "order_id": 10,
      "product_id": 5,
      "variant_id": 10,
      "quantity": 2,
      "price": 99.99,
      "subtotal": 199.98,
      "weight": 1.6
    }
  ],
  "subtotal": 219.97,
  "discount_code": null,
  "discount_amount": 0,
  "shipping_cost": 7.99,
  "final_amount": 227.96,
  "status": "paid",
  "shipping_address": {
    "street_address": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "street_address": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US" 
  },
  "email": "customer@example.com",
  "full_name": "John Smith",
  "phone_number": "+1234567890",
  "shipping_method_id": 3,
  "shipping_method_name": "Standard Shipping",
  "payment_id": "pi_3NJQDLGSwq9VmN8I0bmUrvYx",
  "payment_provider": "stripe",
  "is_guest_order": true,
  "requires_action": false,
  "action_url": "",
  "created_at": "2023-06-25T15:30:45Z",
  "updated_at": "2023-06-25T15:35:20Z"
}
```

Alternative response (when 3D Secure authentication is required):

```json
{
  "id": 10,
  "order_number": "ORD-20230625-000010",
  "user_id": 0,
  "items": [
    {
      "id": 15,
      "order_id": 10,
      "product_id": 3,
      "quantity": 1,
      "price": 19.99,
      "subtotal": 19.99,
      "weight": 0.2
    },
    {
      "id": 16,
      "order_id": 10,
      "product_id": 5,
      "variant_id": 10,
      "quantity": 2,
      "price": 99.99,
      "subtotal": 199.98,
      "weight": 1.6
    }
  ],
  "subtotal": 219.97,
  "discount_code": null,
  "discount_amount": 0,
  "shipping_cost": 7.99,
  "final_amount": 227.96,
  "status": "pending_action",
  "shipping_address": {
    "street_address": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US"
  },
  "billing_address": {
    "street_address": "123 Main St",
    "city": "San Francisco",
    "state": "CA",
    "postal_code": "94105",
    "country": "US" 
  },
  "email": "customer@example.com",
  "full_name": "John Smith",
  "phone_number": "+1234567890",
  "shipping_method_id": 3,
  "shipping_method_name": "Standard Shipping",
  "payment_id": "pi_3NJQDLGSwq9VmN8I0bmUrvYx",
  "payment_provider": "stripe",
  "is_guest_order": true,
  "requires_action": true,
  "action_url": "https://hooks.stripe.com/3d_secure_2_eap/begin_test/src_1NJQDLGSwq9VmN8I0OOVbLwE/src_client_secret_CG9LMEyAnFQw9OdPvRD0NCmz",
  "created_at": "2023-06-25T15:30:45Z",
  "updated_at": "2023-06-25T15:35:20Z"
}
```

**Status Codes:**

- `200 OK`: Payment processed or requires further action
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
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "billing_address": {
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "shipping_method_id": 2,
  "phone_number": "+1987654321"
}
```

Example response:

```json
{
  "id": 12,
  "order_number": "ORD-20230626-000012",
  "user_id": 5,
  "items": [
    {
      "id": 20,
      "order_id": 12,
      "product_id": 1,
      "quantity": 1,
      "price": 999.99,
      "subtotal": 999.99,
      "weight": 0.35
    },
    {
      "id": 21,
      "order_id": 12,
      "product_id": 2,
      "variant_id": 1,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99,
      "weight": 2.1
    }
  ],
  "subtotal": 2499.98,
  "discount_code": null,
  "discount_amount": 0,
  "shipping_cost": 14.99,
  "final_amount": 2514.97,
  "status": "pending",
  "shipping_address": {
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "billing_address": {
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "email": "user@example.com",
  "full_name": "Sarah Johnson",
  "phone_number": "+1987654321",
  "shipping_method_id": 2,
  "shipping_method_name": "Express Shipping",
  "payment_id": null,
  "payment_provider": null,
  "is_guest_order": false,
  "requires_action": false,
  "action_url": "",
  "created_at": "2023-06-26T10:15:30Z",
  "updated_at": "2023-06-26T10:15:30Z"
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
  "id": 12,
  "order_number": "ORD-20230626-000012",
  "user_id": 5,
  "items": [
    {
      "id": 20,
      "order_id": 12,
      "product_id": 1,
      "quantity": 1,
      "price": 999.99,
      "subtotal": 999.99,
      "weight": 0.35
    },
    {
      "id": 21,
      "order_id": 12,
      "product_id": 2,
      "variant_id": 1,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99,
      "weight": 2.1
    }
  ],
  "subtotal": 2499.98,
  "discount_code": null,
  "discount_amount": 0,
  "shipping_cost": 14.99,
  "final_amount": 2514.97,
  "status": "pending_action",
  "shipping_address": {
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "billing_address": {
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "email": "user@example.com",
  "full_name": "Sarah Johnson",
  "phone_number": "+1987654321",
  "shipping_method_id": 2,
  "shipping_method_name": "Express Shipping",
  "payment_id": "mp-123456789",
  "payment_provider": "mobilepay",
  "is_guest_order": false,
  "requires_action": true,
  "action_url": "https://api.mobilepay.dk/v1/payments/mp-123456789/authorize",
  "created_at": "2023-06-26T10:15:30Z",
  "updated_at": "2023-06-26T10:18:45Z"
}
```

**Status Codes:**

- `200 OK`: Payment processed or requires further action
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
  "id": 12,
  "order_number": "ORD-20230626-000012",
  "user_id": 5,
  "items": [
    {
      "id": 20,
      "order_id": 12,
      "product_id": 1,
      "quantity": 1,
      "price": 999.99,
      "subtotal": 999.99,
      "weight": 0.35
    },
    {
      "id": 21,
      "order_id": 12,
      "product_id": 2,
      "variant_id": 1,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99,
      "weight": 2.1
    }
  ],
  "subtotal": 2499.98,
  "discount_code": null,
  "discount_amount": 0,
  "shipping_cost": 14.99,
  "final_amount": 2514.97,
  "status": "paid",
  "shipping_address": {
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "billing_address": {
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "email": "user@example.com",
  "full_name": "Sarah Johnson",
  "phone_number": "+1987654321",
  "shipping_method_id": 2,
  "shipping_method_name": "Express Shipping",
  "payment_id": "mp-123456789",
  "payment_provider": "mobilepay",
  "is_guest_order": false,
  "requires_action": false,
  "action_url": "",
  "created_at": "2023-06-26T10:15:30Z",
  "updated_at": "2023-06-26T10:25:12Z"
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
[
  {
    "id": 12,
    "order_number": "ORD-20230626-000012",
    "user_id": 5,
    "subtotal": 2499.98,
    "discount_code": null,
    "discount_amount": 0,
    "shipping_cost": 14.99,
    "final_amount": 2514.97,
    "status": "paid",
    "items_count": 2,
    "shipping_method_name": "Express Shipping",
    "payment_provider": "mobilepay",
    "created_at": "2023-06-26T10:15:30Z",
    "updated_at": "2023-06-26T10:25:12Z"
  },
  {
    "id": 9,
    "order_number": "ORD-20230620-000009",
    "user_id": 5,
    "subtotal": 39.98,
    "discount_code": "SUMMER2023",
    "discount_amount": 4.00,
    "shipping_cost": 5.99,
    "final_amount": 41.97,
    "status": "shipped",
    "items_count": 2,
    "shipping_method_name": "Standard Shipping",
    "payment_provider": "stripe",
    "created_at": "2023-06-20T14:22:15Z",
    "updated_at": "2023-06-22T09:30:00Z"
  }
]
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
[
  {
    "id": 12,
    "order_number": "ORD-20230626-000012",
    "user_id": 5,
    "email": "user@example.com",
    "full_name": "Sarah Johnson",
    "subtotal": 2499.98,
    "discount_amount": 0,
    "shipping_cost": 14.99,
    "final_amount": 2514.97,
    "status": "paid",
    "items_count": 2,
    "payment_provider": "mobilepay",
    "created_at": "2023-06-26T10:15:30Z",
    "updated_at": "2023-06-26T10:25:12Z"
  },
  {
    "id": 10,
    "order_number": "ORD-20230625-000010",
    "user_id": 0,
    "email": "customer@example.com",
    "full_name": "John Smith",
    "subtotal": 219.97,
    "discount_amount": 0,
    "shipping_cost": 7.99,
    "final_amount": 227.96,
    "status": "paid",
    "items_count": 2,
    "payment_provider": "stripe",
    "is_guest_order": true,
    "created_at": "2023-06-25T15:30:45Z",
    "updated_at": "2023-06-25T15:35:20Z"
  }
]
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
  "id": 12,
  "order_number": "ORD-20230626-000012",
  "user_id": 5,
  "items": [
    {
      "id": 20,
      "order_id": 12,
      "product_id": 1,
      "quantity": 1,
      "price": 999.99,
      "subtotal": 999.99,
      "weight": 0.35
    },
    {
      "id": 21,
      "order_id": 12,
      "product_id": 2,
      "variant_id": 1,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99,
      "weight": 2.1
    }
  ],
  "subtotal": 2499.98,
  "discount_code": null,
  "discount_amount": 0,
  "shipping_cost": 14.99,
  "final_amount": 2514.97,
  "status": "shipped",
  "shipping_address": {
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "billing_address": {
    "street_address": "456 Oak Avenue",
    "city": "Seattle",
    "state": "WA",
    "postal_code": "98101",
    "country": "US"
  },
  "email": "user@example.com",
  "full_name": "Sarah Johnson",
  "phone_number": "+1987654321",
  "shipping_method_id": 2,
  "shipping_method_name": "Express Shipping",
  "payment_id": "mp-123456789",
  "payment_provider": "mobilepay",
  "is_guest_order": false,
  "requires_action": false,
  "action_url": "",
  "created_at": "2023-06-26T10:15:30Z",
  "updated_at": "2023-06-26T14:30:15Z"
}
```

**Status Codes:**

- `200 OK`: Order status updated successfully
- `400 Bad Request`: Invalid order status
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `404 Not Found`: Order not found
- `500 Internal Server Error`: Failed to update order status

## Payment Management Endpoints (Admin Only)

### Capture Payment

```plaintext
POST /api/admin/payments/{paymentId}/capture
```

Capture a previously authorized payment (admin only).

**Request Body:**

```json
{
  "amount": 2514.97
}
```

**Status Codes:**

- `200 OK`: Payment captured successfully
- `400 Bad Request`: Invalid request or capture not allowed
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `404 Not Found`: Payment not found
- `500 Internal Server Error`: Failed to capture payment

### Cancel Payment

```plaintext
POST /api/admin/payments/{paymentId}/cancel
```

Cancel a payment that requires action but hasn't been completed (admin only).

**Status Codes:**

- `200 OK`: Payment cancelled successfully
- `400 Bad Request`: Payment cancellation not allowed
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `404 Not Found`: Payment not found
- `500 Internal Server Error`: Failed to cancel payment

### Refund Payment

```plaintext
POST /api/admin/payments/{paymentId}/refund
```

Refund a captured payment (admin only).

**Request Body:**

```json
{
  "amount": 2514.97
}
```

**Status Codes:**

- `200 OK`: Payment refunded successfully
- `400 Bad Request`: Invalid request or refund not allowed
- `401 Unauthorized`: User not authenticated
- `403 Forbidden`: User not authorized (not an admin)
- `404 Not Found`: Payment not found
- `500 Internal Server Error`: Failed to refund payment

## Example Workflow

### Guest Checkout Flow

1. Guest adds items to their cart
2. Guest provides shipping information and selects shipping method
3. System creates an order with `POST /api/guest/orders`
4. Guest provides payment details with `POST /api/guest/orders/{id}/payment`
5. If payment requires additional action (3D Secure, etc.), guest completes it
6. Payment is processed and order status is updated to "paid"

### Authenticated User Checkout Flow

1. User adds items to their cart
2. User provides shipping information and selects shipping method
3. System creates an order with `POST /api/orders`
4. User provides payment details with `POST /api/orders/{id}/payment`
5. If payment requires additional action, user completes it
6. Payment is processed and order status is updated to "paid"

### Order Fulfillment Flow (Admin)

1. Admin views orders with `GET /api/admin/orders`
2. Admin captures payment if necessary with `POST /api/admin/payments/{paymentId}/capture`
3. Admin processes the order (picking, packing)
4. Admin updates order status to "shipped" with `PUT /api/admin/orders/{id}/status`
5. System sends shipping confirmation email to customer
6. When delivery is confirmed, admin updates status to "delivered"
