# Cart API Examples

This document provides example request bodies for the shopping cart system API endpoints.

## Public Cart Endpoints

### Get Cart

```plaintext
GET /api/guest/cart
```

Retrieves the current guest cart. Uses a cookie to identify the session.

Example response:

```json
{
  "id": 42,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:15:45Z",
  "user_id": null,
  "items": [
    {
      "id": 15,
      "product_id": "00000000-0000-0000-0000-000000000003",
      "variant_id": null,
      "quantity": 1,
      "unit_price": 19.99,
      "total_price": 19.99
    },
    {
      "id": 16,
      "product_id": "00000000-0000-0000-0000-000000000005",
      "variant_id": "00000000-0000-0000-0000-00000000000a",
      "quantity": 2,
      "unit_price": 99.99,
      "total_price": 199.98
    }
  ],
  "total_amount": 219.97,
  "currency": "USD",
  "discount_amount": 0
}
```

**Status Codes:**

- `200 OK`: Cart retrieved successfully
- `500 Internal Server Error`: Failed to retrieve cart

### Add Item to Cart

```plaintext
POST /api/guest/cart/items
```

Adds a product or variant to the guest cart.

**Request Body:**

```json
{
  "product_id": "00000000-0000-0000-0000-000000000005",
  "variant_id": "00000000-0000-0000-0000-00000000000a",
  "quantity": 2
}
```

Example response:

```json
{
  "id": 42,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:15:45Z",
  "user_id": null,
  "items": [
    {
      "id": 15,
      "product_id": "00000000-0000-0000-0000-000000000003",
      "variant_id": null,
      "quantity": 1,
      "unit_price": 19.99,
      "total_price": 19.99
    },
    {
      "id": 16,
      "product_id": "00000000-0000-0000-0000-000000000005",
      "variant_id": "00000000-0000-0000-0000-00000000000a",
      "quantity": 2,
      "unit_price": 99.99,
      "total_price": 199.98
    }
  ],
  "total_amount": 219.97,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Item added successfully
- `400 Bad Request`: Invalid request body, product not found, or insufficient stock
- `500 Internal Server Error`: Failed to add item

### Update Cart Item

```plaintext
PUT /api/guest/cart/items/{productId}?variantId={variantId}
```

Updates the quantity of a product in the guest cart.

**Request Body:**

```json
{
  "quantity": 3
}
```

Example response:

```json
{
  "id": 42,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:25:30Z",
  "user_id": null,
  "items": [
    {
      "id": 15,
      "product_id": "00000000-0000-0000-0000-000000000003",
      "variant_id": null,
      "quantity": 1,
      "unit_price": 19.99,
      "total_price": 19.99
    },
    {
      "id": 16,
      "product_id": "00000000-0000-0000-0000-000000000005",
      "variant_id": "00000000-0000-0000-0000-00000000000a",
      "quantity": 3,
      "unit_price": 99.99,
      "total_price": 299.97
    }
  ],
  "total_amount": 319.96,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Item updated successfully
- `400 Bad Request`: Invalid request body, product not found, or insufficient stock
- `500 Internal Server Error`: Failed to update item

### Remove Cart Item

```plaintext
DELETE /api/guest/cart/items/{productId}?variantId={variantId}
```

Removes a product or variant from the guest cart.

Example response:

```json
{
  "id": 42,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:35:12Z",
  "user_id": null,
  "items": [
    {
      "id": 15,
      "product_id": "00000000-0000-0000-0000-000000000003",
      "variant_id": null,
      "quantity": 1,
      "unit_price": 19.99,
      "total_price": 19.99
    }
  ],
  "total_amount": 19.99,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Item removed successfully
- `400 Bad Request`: Product not found in cart
- `500 Internal Server Error`: Failed to remove item

### Clear Cart

```plaintext
DELETE /api/guest/cart
```

Removes all items from the guest cart.

Example response:

```json
{
  "id": 42,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:45:20Z",
  "user_id": null,
  "items": [],
  "total_amount": 0,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Cart cleared successfully
- `500 Internal Server Error`: Failed to clear cart

## User Cart Endpoints (Authenticated)

### Get User Cart

```plaintext
GET /api/cart
```

Retrieves the current authenticated user's cart.

Example response:

```json
{
  "id": 34,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T09:45:22Z",
  "user_id": "00000000-0000-0000-0000-00000000000c",
  "items": [
    {
      "id": 55,
      "product_id": "00000000-0000-0000-0000-000000000001",
      "variant_id": null,
      "quantity": 1,
      "unit_price": 999.99,
      "total_price": 999.99
    },
    {
      "id": 56,
      "product_id": "00000000-0000-0000-0000-000000000002",
      "variant_id": "00000000-0000-0000-0000-000000000001",
      "quantity": 1,
      "unit_price": 1499.99,
      "total_price": 1499.99
    }
  ],
  "total_amount": 2499.98,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Cart retrieved successfully
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to retrieve cart

### Add Item to User Cart

```plaintext
POST /api/cart/items
```

Adds a product or variant to the authenticated user's cart.

**Request Body:**

```json
{
  "product_id": "00000000-0000-0000-0000-000000000002",
  "variant_id": "00000000-0000-0000-0000-000000000002",
  "quantity": 1
}
```

Example response:

```json
{
  "id": 34,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T10:05:15Z",
  "user_id": "00000000-0000-0000-0000-00000000000c",
  "items": [
    {
      "id": 55,
      "product_id": "00000000-0000-0000-0000-000000000001",
      "variant_id": null,
      "quantity": 1,
      "unit_price": 999.99,
      "total_price": 999.99
    },
    {
      "id": 56,
      "product_id": "00000000-0000-0000-0000-000000000002",
      "variant_id": "00000000-0000-0000-0000-000000000001",
      "quantity": 1,
      "unit_price": 1499.99,
      "total_price": 1499.99
    },
    {
      "id": 57,
      "product_id": "00000000-0000-0000-0000-000000000002",
      "variant_id": "00000000-0000-0000-0000-000000000002",
      "quantity": 1,
      "unit_price": 1799.99,
      "total_price": 1799.99
    }
  ],
  "total_amount": 4299.97,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Item added successfully
- `400 Bad Request`: Invalid request body, product not found, or insufficient stock
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to add item

### Update User Cart Item

```plaintext
PUT /api/cart/items/{productId}?variantId={variantId}
```

Updates the quantity of a product in the authenticated user's cart.

**Request Body:**

```json
{
  "quantity": 2
}
```

Example response:

```json
{
  "id": 34,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T10:15:45Z",
  "user_id": "00000000-0000-0000-0000-00000000000c",
  "items": [
    {
      "id": 55,
      "product_id": "00000000-0000-0000-0000-000000000001",
      "variant_id": null,
      "quantity": 1,
      "unit_price": 999.99,
      "total_price": 999.99
    },
    {
      "id": 56,
      "product_id": "00000000-0000-0000-0000-000000000002",
      "variant_id": "00000000-0000-0000-0000-000000000001",
      "quantity": 2,
      "unit_price": 1499.99,
      "total_price": 2999.98
    },
    {
      "id": 57,
      "product_id": "00000000-0000-0000-0000-000000000002",
      "variant_id": "00000000-0000-0000-0000-000000000002",
      "quantity": 1,
      "unit_price": 1799.99,
      "total_price": 1799.99
    }
  ],
  "total_amount": 5799.96,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Item updated successfully
- `400 Bad Request`: Invalid request body, product not found, or insufficient stock
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to update item

### Remove User Cart Item

```plaintext
DELETE /api/cart/items/{productId}?variantId={variantId}
```

Removes a product or variant from the authenticated user's cart.

Example response:

```json
{
  "id": 34,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T10:25:30Z",
  "user_id": "00000000-0000-0000-0000-00000000000c",
  "items": [
    {
      "id": 55,
      "product_id": "00000000-0000-0000-0000-000000000001",
      "variant_id": null,
      "quantity": 1,
      "unit_price": 999.99,
      "total_price": 999.99
    },
    {
      "id": 57,
      "product_id": "00000000-0000-0000-0000-000000000002",
      "variant_id": "00000000-0000-0000-0000-000000000002",
      "quantity": 1,
      "unit_price": 1799.99,
      "total_price": 1799.99
    }
  ],
  "total_amount": 2799.98,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Item removed successfully
- `400 Bad Request`: Product not found in cart
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to remove item

### Clear User Cart

```plaintext
DELETE /api/cart
```

Removes all items from the authenticated user's cart.

Example response:

```json
{
  "id": 34,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T10:35:15Z",
  "user_id": "00000000-0000-0000-0000-00000000000c",
  "items": [],
  "total_amount": 0,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Cart cleared successfully
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to clear cart

## Session Conversion

### Convert Guest Cart to User Cart

```plaintext
POST /api/guest/cart/convert
```

Converts a guest cart to an authenticated user's cart after login. Requires authentication.

Example response:

```json
{
  "id": 34,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T11:05:22Z",
  "user_id": "00000000-0000-0000-0000-00000000000c",
  "items": [
    {
      "id": 55,
      "product_id": "00000000-0000-0000-0000-000000000001",
      "variant_id": null,
      "quantity": 1,
      "unit_price": 999.99,
      "total_price": 999.99
    },
    {
      "id": 56,
      "product_id": "00000000-0000-0000-0000-000000000005",
      "variant_id": "00000000-0000-0000-0000-00000000000a",
      "quantity": 2,
      "unit_price": 99.99,
      "total_price": 199.98
    }
  ],
  "total_amount": 1199.97,
  "currency": "USD"
}
```

**Status Codes:**

- `200 OK`: Cart converted successfully
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to convert cart

## Example Workflow

1. Guest user adds items to cart
2. When the guest registers or logs in, their guest cart is converted to a user cart
3. User continues shopping, updating quantities or removing items as needed
4. When ready to check out, the cart contents are used to create an order
5. After successful order creation, the cart is typically cleared
