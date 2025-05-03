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
  "session_id": "0ed9a755-48a9-4d6c-b544-df3f325d37b2",
  "user_id": null,
  "items": [
    {
      "id": 15,
      "cart_id": 42,
      "product_id": 3,
      "variant_id": 0,
      "quantity": 1,
      "price": 19.99,
      "subtotal": 19.99
    },
    {
      "id": 16,
      "cart_id": 42,
      "product_id": 5,
      "variant_id": 10,
      "quantity": 2,
      "price": 99.99,
      "subtotal": 199.98
    }
  ],
  "total_items": 3,
  "subtotal": 219.97,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:15:45Z"
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
  "product_id": 5,
  "variant_id": 10,
  "quantity": 2
}
```

Example response:

```json
{
  "id": 42,
  "session_id": "0ed9a755-48a9-4d6c-b544-df3f325d37b2",
  "user_id": null,
  "items": [
    {
      "id": 15,
      "cart_id": 42,
      "product_id": 3,
      "variant_id": 0,
      "quantity": 1,
      "price": 19.99,
      "subtotal": 19.99
    },
    {
      "id": 16,
      "cart_id": 42,
      "product_id": 5,
      "variant_id": 10,
      "quantity": 2,
      "price": 99.99,
      "subtotal": 199.98
    }
  ],
  "total_items": 3,
  "subtotal": 219.97,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:15:45Z"
}
```

**Status Codes:**

- `200 OK`: Item added successfully
- `400 Bad Request`: Invalid request body, product not found, or insufficient stock
- `500 Internal Server Error`: Failed to add item

### Update Cart Item

```plaintext
PUT /api/guest/cart/items/{productId}
```

Updates the quantity of a product in the guest cart.

**Request Body:**

```json
{
  "quantity": 3,
  "variant_id": 10
}
```

Example response:

```json
{
  "id": 42,
  "session_id": "0ed9a755-48a9-4d6c-b544-df3f325d37b2",
  "user_id": null,
  "items": [
    {
      "id": 15,
      "cart_id": 42,
      "product_id": 3,
      "variant_id": 0,
      "quantity": 1,
      "price": 19.99,
      "subtotal": 19.99
    },
    {
      "id": 16,
      "cart_id": 42,
      "product_id": 5,
      "variant_id": 10,
      "quantity": 3,
      "price": 99.99,
      "subtotal": 299.97
    }
  ],
  "total_items": 4,
  "subtotal": 319.96,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:25:30Z"
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
  "session_id": "0ed9a755-48a9-4d6c-b544-df3f325d37b2",
  "user_id": null,
  "items": [
    {
      "id": 15,
      "cart_id": 42,
      "product_id": 3,
      "variant_id": 0,
      "quantity": 1,
      "price": 19.99,
      "subtotal": 19.99
    }
  ],
  "total_items": 1,
  "subtotal": 19.99,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:35:12Z"
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
  "session_id": "0ed9a755-48a9-4d6c-b544-df3f325d37b2",
  "user_id": null,
  "items": [],
  "total_items": 0,
  "subtotal": 0.00,
  "created_at": "2023-06-15T10:30:22Z",
  "updated_at": "2023-06-15T11:45:20Z"
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
  "session_id": null,
  "user_id": 12,
  "items": [
    {
      "id": 55,
      "cart_id": 34,
      "product_id": 1,
      "variant_id": 0,
      "quantity": 1,
      "price": 999.99,
      "subtotal": 999.99
    },
    {
      "id": 56,
      "cart_id": 34,
      "product_id": 2,
      "variant_id": 1,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99
    }
  ],
  "total_items": 2,
  "subtotal": 2499.98,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T09:45:22Z"
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
  "product_id": 2,
  "variant_id": 2,
  "quantity": 1
}
```

Example response:

```json
{
  "id": 34,
  "session_id": null,
  "user_id": 12,
  "items": [
    {
      "id": 55,
      "cart_id": 34,
      "product_id": 1,
      "variant_id": 0,
      "quantity": 1,
      "price": 999.99,
      "subtotal": 999.99
    },
    {
      "id": 56,
      "cart_id": 34,
      "product_id": 2,
      "variant_id": 1,
      "quantity": 1,
      "price": 1499.99,
      "subtotal": 1499.99
    },
    {
      "id": 57,
      "cart_id": 34,
      "product_id": 2,
      "variant_id": 2,
      "quantity": 1,
      "price": 1799.99,
      "subtotal": 1799.99
    }
  ],
  "total_items": 3,
  "subtotal": 4299.97,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T10:05:15Z"
}
```

**Status Codes:**

- `200 OK`: Item added successfully
- `400 Bad Request`: Invalid request body, product not found, or insufficient stock
- `401 Unauthorized`: User not authenticated
- `500 Internal Server Error`: Failed to add item

### Update User Cart Item

```plaintext
PUT /api/cart/items/{productId}
```

Updates the quantity of a product in the authenticated user's cart.

**Request Body:**

```json
{
  "quantity": 2,
  "variant_id": 1
}
```

Example response:

```json
{
  "id": 34,
  "session_id": null,
  "user_id": 12,
  "items": [
    {
      "id": 55,
      "cart_id": 34,
      "product_id": 1,
      "variant_id": 0,
      "quantity": 1,
      "price": 999.99,
      "subtotal": 999.99
    },
    {
      "id": 56,
      "cart_id": 34,
      "product_id": 2,
      "variant_id": 1,
      "quantity": 2,
      "price": 1499.99,
      "subtotal": 2999.98
    },
    {
      "id": 57,
      "cart_id": 34,
      "product_id": 2,
      "variant_id": 2,
      "quantity": 1,
      "price": 1799.99,
      "subtotal": 1799.99
    }
  ],
  "total_items": 4,
  "subtotal": 5799.96,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T10:15:45Z"
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
  "session_id": null,
  "user_id": 12,
  "items": [
    {
      "id": 55,
      "cart_id": 34,
      "product_id": 1,
      "variant_id": 0,
      "quantity": 1,
      "price": 999.99,
      "subtotal": 999.99
    },
    {
      "id": 57,
      "cart_id": 34,
      "product_id": 2,
      "variant_id": 2,
      "quantity": 1,
      "price": 1799.99,
      "subtotal": 1799.99
    }
  ],
  "total_items": 2,
  "subtotal": 2799.98,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T10:25:30Z"
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
  "session_id": null,
  "user_id": 12,
  "items": [],
  "total_items": 0,
  "subtotal": 0.00,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T10:35:15Z"
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
  "session_id": null,
  "user_id": 12,
  "items": [
    {
      "id": 55,
      "cart_id": 34,
      "product_id": 1,
      "variant_id": 0,
      "quantity": 1,
      "price": 999.99,
      "subtotal": 999.99
    },
    {
      "id": 56,
      "cart_id": 34,
      "product_id": 5,
      "variant_id": 10,
      "quantity": 2,
      "price": 99.99,
      "subtotal": 199.98
    }
  ],
  "total_items": 3,
  "subtotal": 1199.97,
  "created_at": "2023-06-16T09:15:30Z",
  "updated_at": "2023-06-16T11:05:22Z"
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
