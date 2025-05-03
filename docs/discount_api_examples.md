# Discount API Examples

This document provides example request bodies for the discount system API endpoints.

## Public Discount Endpoints

### List Active Discounts

`GET /api/discounts`

List all currently active discounts.

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

Example response:

```json
[
  {
    "id": 1,
    "code": "SUMMER2023",
    "type": "basket",
    "method": "percentage",
    "value": 10.0,
    "min_order_value": 50.0,
    "max_discount_value": 20.0,
    "start_date": "2023-06-01T00:00:00Z",
    "end_date": "2023-08-31T23:59:59Z",
    "usage_limit": 1000,
    "current_usage": 243,
    "active": true
  },
  {
    "id": 2,
    "code": "WELCOME10",
    "type": "basket",
    "method": "fixed",
    "value": 10.0,
    "min_order_value": 0.0,
    "max_discount_value": 10.0,
    "start_date": "2023-01-01T00:00:00Z",
    "end_date": "2023-12-31T23:59:59Z",
    "usage_limit": 0,
    "current_usage": 567,
    "active": true
  }
]
```

### Apply Discount to Order

`POST /api/orders/{id}/discounts`

Apply a discount code to an existing order.

```json
{
  "discount_code": "SUMMER2023"
}
```

Example response:

```json
{
  "id": 5,
  "user_id": 1,
  "items": [
    {
      "id": 12,
      "order_id": 5,
      "product_id": 3,
      "quantity": 2,
      "price": 24.99,
      "subtotal": 49.98
    }
  ],
  "subtotal": 49.98,
  "discount_code": "SUMMER2023",
  "discount_amount": 5.00,
  "shipping_cost": 5.99,
  "total_amount": 50.97,
  "status": "pending",
  "created_at": "2023-06-15T14:22:15Z",
  "updated_at": "2023-06-15T14:23:05Z"
}
```

### Remove Discount from Order

`DELETE /api/orders/{id}/discounts`

Remove an applied discount from an order.

Example response:

```json
{
  "id": 5,
  "user_id": 1,
  "items": [
    {
      "id": 12,
      "order_id": 5,
      "product_id": 3,
      "quantity": 2,
      "price": 24.99,
      "subtotal": 49.98
    }
  ],
  "subtotal": 49.98,
  "discount_code": null,
  "discount_amount": 0.00,
  "shipping_cost": 5.99,
  "total_amount": 55.97,
  "status": "pending",
  "created_at": "2023-06-15T14:22:15Z",
  "updated_at": "2023-06-15T14:24:30Z"
}
```

## Admin Discount Endpoints

### Create Discount

`POST /api/admin/discounts`

Create a new discount.

```json
{
  "code": "SUMMER2023",
  "type": "basket",
  "method": "percentage",
  "value": 10.0,
  "min_order_value": 50.0,
  "max_discount_value": 20.0,
  "product_ids": [],
  "category_ids": [],
  "start_date": "2023-06-01T00:00:00Z",
  "end_date": "2023-08-31T23:59:59Z",
  "usage_limit": 1000,
  "active": true
}
```

Example response:

```json
{
  "id": 3,
  "code": "SUMMER2023",
  "type": "basket",
  "method": "percentage",
  "value": 10.0,
  "min_order_value": 50.0,
  "max_discount_value": 20.0,
  "start_date": "2023-06-01T00:00:00Z",
  "end_date": "2023-08-31T23:59:59Z",
  "usage_limit": 1000,
  "current_usage": 0,
  "active": true,
  "created_at": "2023-05-15T10:30:00Z",
  "updated_at": "2023-05-15T10:30:00Z"
}
```

### Update Discount

`PUT /api/admin/discounts/{id}`

Update an existing discount.

```json
{
  "code": "SUMMER2023",
  "type": "basket",
  "method": "percentage",
  "value": 15.0,
  "min_order_value": 40.0,
  "max_discount_value": 25.0,
  "product_ids": [],
  "category_ids": [],
  "start_date": "2023-06-01T00:00:00Z",
  "end_date": "2023-09-15T23:59:59Z",
  "usage_limit": 2000,
  "active": true
}
```

Example response:

```json
{
  "id": 3,
  "code": "SUMMER2023",
  "type": "basket",
  "method": "percentage",
  "value": 15.0,
  "min_order_value": 40.0,
  "max_discount_value": 25.0,
  "start_date": "2023-06-01T00:00:00Z",
  "end_date": "2023-09-15T23:59:59Z",
  "usage_limit": 2000,
  "current_usage": 0,
  "active": true,
  "created_at": "2023-05-15T10:30:00Z",
  "updated_at": "2023-05-16T09:45:22Z"
}
```

### Delete Discount

`DELETE /api/admin/discounts/{id}`

Delete an existing discount.

**Status Codes:**

- `204 No Content`: Discount deleted successfully
- `400 Bad Request`: Cannot delete discount in use by orders
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Discount not found

### List All Discounts

`GET /api/admin/discounts`

List all discounts (active and inactive).

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

Example response:

```json
[
  {
    "id": 1,
    "code": "SUMMER2023",
    "type": "basket",
    "method": "percentage",
    "value": 10.0,
    "min_order_value": 50.0,
    "max_discount_value": 20.0,
    "start_date": "2023-06-01T00:00:00Z",
    "end_date": "2023-08-31T23:59:59Z",
    "usage_limit": 1000,
    "current_usage": 243,
    "active": true,
    "created_at": "2023-05-01T10:00:00Z",
    "updated_at": "2023-06-22T15:34:17Z"
  },
  {
    "id": 2,
    "code": "WELCOME10",
    "type": "basket",
    "method": "fixed",
    "value": 10.0,
    "min_order_value": 0.0,
    "max_discount_value": 10.0,
    "start_date": "2023-01-01T00:00:00Z",
    "end_date": "2023-12-31T23:59:59Z",
    "usage_limit": 0,
    "current_usage": 567,
    "active": true,
    "created_at": "2022-12-15T09:00:00Z",
    "updated_at": "2023-06-20T11:22:05Z"
  },
  {
    "id": 3,
    "code": "FLASH50",
    "type": "product",
    "method": "percentage",
    "value": 50.0,
    "min_order_value": 0.0,
    "max_discount_value": 0.0,
    "start_date": "2023-04-01T00:00:00Z",
    "end_date": "2023-04-03T23:59:59Z",
    "usage_limit": 500,
    "current_usage": 500,
    "active": false,
    "created_at": "2023-03-25T16:45:00Z",
    "updated_at": "2023-04-03T23:59:59Z"
  }
]
```

## Example Workflow

1. Create a new discount through the admin interface
2. Users see active discounts when shopping
3. During checkout, users apply a discount code to their order
4. The system validates the discount and applies it to the order if valid
5. The discount usage counter increments after successful order completion
6. Admin can modify or deactivate discounts as needed
