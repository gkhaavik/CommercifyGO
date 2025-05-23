# Migrating from Cart to Checkout

## Overview

As of May 23, 2025, the cart system has been deprecated in favor of the new checkout system. This document provides guidance for migrating from the cart API to the checkout API.

## Migration Strategy

### Database Changes

A migration (`000021_disable_cart_tables.up.sql`) has been added that:

1. Archives existing cart data to `cart_archive` and `cart_items_archive` tables
2. Adds triggers to prevent inserts and updates to cart tables
3. Creates views (`legacy_carts` and `legacy_cart_items`) that provide read-only access to cart data
4. Adds comments to the tables indicating they are deprecated

### API Changes

The cart API endpoints have been deprecated. The following table shows the mapping from cart endpoints to checkout endpoints:

| Cart Endpoint | Checkout Equivalent |
|---------------|---------------------|
| `GET /api/guest/cart` | `GET /api/guest/checkout` |
| `POST /api/guest/cart/items` | `POST /api/guest/checkout/items` |
| `PUT /api/guest/cart/items/{productId}` | `PUT /api/guest/checkout/items/{productId}` |
| `DELETE /api/guest/cart/items/{productId}` | `DELETE /api/guest/checkout/items/{productId}` |
| `DELETE /api/guest/cart` | `DELETE /api/guest/checkout` |

### Client Changes

The `CommercifyClient` TypeScript class has been updated to mark cart methods as deprecated. Use the following equivalents instead:

| Deprecated Method | Replacement Method |
|------------------|-------------------|
| `getCart()` | `getCheckout()` |
| `addToCart()` | `addToCheckout()` |
| `updateCartItem()` | `updateCheckoutItem()` |
| `removeCartItem()` | `removeFromCheckout()` |
| `clearCart()` | `clearCheckout()` |

## Advantages of the Checkout System

The new checkout system offers several improvements over the cart system:

1. **Unified Guest and User Experience**: The same API endpoints are used for both guest and authenticated users
2. **Extended Functionality**: Support for shipping methods, discount codes, and customer details
3. **Better Data Integrity**: Improved validation and error handling
4. **Performance Improvements**: Optimized database queries and caching
5. **Multi-currency Support**: Built-in support for multiple currencies

## Examples

### Example: Adding an Item to Checkout

```typescript
// Old cart approach
const cartResponse = await client.addToCart({
  productId: 123,
  variantId: 456,
  quantity: 2
});

// New checkout approach
const checkoutResponse = await client.addToCheckout({
  productId: 123,
  variantId: 456,
  quantity: 2
});
```

### Example: Retrieving Checkout Contents

```typescript
// Old cart approach
const cartResponse = await client.getCart();

// New checkout approach
const checkoutResponse = await client.getCheckout();
```

## Timeline

- **May 23, 2025**: Cart API deprecated (current)
- **August 1, 2025**: Cart API endpoints will return 410 Gone
- **January 1, 2026**: Cart tables will be removed from the database

For more information on the checkout API, please refer to the [Checkout API Examples](./checkout_api_examples.md) document.
