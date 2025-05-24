### Complete Checkout

```plaintext
POST /api/checkout/complete
```

**Description:**
Converts the current checkout to an order. The checkout is identified by the checkout_session_id cookie.

**Request Body:**

```json
{
  "user_id": 123, // Optional - only needed to link to a user account
}
```

**Response Body:**

```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": 42,
    "order_number": "ORD-20250524-00042",
    "user_id": 123,
    "status": "pending",
    "total_amount": 149.95,
    "final_amount": 139.95,
    "currency": "USD",
    "items": [
      {
        "id": 85,
        "order_id": 42,
        "product_id": 12,
        "quantity": 1,
        "unit_price": 149.95,
        "total_price": 149.95,
        "created_at": "2025-05-24T10:20:29.143443+02:00",
        "updated_at": "2025-05-24T10:20:29.143443+02:00"
      }
    ],
    "shipping_address": {
      "address_line1": "123 Main St",
      "address_line2": "Apt 4B",
      "city": "New York",
      "state": "NY",
      "postal_code": "10001",
      "country": "US"
    },
    "billing_address": {
      "address_line1": "123 Main St",
      "address_line2": "Apt 4B",
      "city": "New York",
      "state": "NY", 
      "postal_code": "10001",
      "country": "US"
    },
    "payment_details": {
      "provider": "stripe",
      "method": "credit_card",
      "id": "",
      "status": "",
      "captured": false,
      "refunded": false
    },
    "shipping_details": {
      "method_id": 1,
      "method": "Standard Shipping",
      "cost": 10.00
    },
    "discount_details": {
      "code": "SUMMER2025",
      "amount": 20.00
    },
    "customer": {
      "email": "customer@example.com",
      "phone": "+1234567890",
      "full_name": "John Doe"
    },
    "checkout_session_id": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
    "created_at": "2025-05-24T10:20:29.143443+02:00",
    "updated_at": "2025-05-24T10:20:29.143443+02:00"
  }
}
```

**Status Codes:**

- `201 Created`: Order created successfully
- `400 Bad Request`: Invalid request or checkout has no items
- `401 Unauthorized`: Not authenticated (if trying to associate with a user)
- `404 Not Found`: No checkout found with this session ID
- `500 Internal Server Error`: Server error

### Checkout Flow

1. A checkout session ID is automatically created and stored as a cookie (`checkout_session_id`) when the user interacts with the checkout system.
2. All checkout operations (adding/removing items, setting shipping address, etc.) are tied to this checkout session.
3. The checkout can be associated with a user account if the user is authenticated.
4. When the checkout is completed, an order is created, and the checkout session ID is stored with the order.
