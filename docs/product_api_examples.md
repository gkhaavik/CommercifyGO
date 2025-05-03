# Product API Examples

This document provides example request bodies for the product system API endpoints.

## Public Product Endpoints

### List Products

`GET /api/products`

List all products with pagination.

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

Example response:

```json
[
  {
    "id": 1,
    "product_number": "PROD-000001",
    "name": "Smartphone",
    "description": "Latest smartphone model",
    "price": 999.99,
    "stock": 50,
    "weight": 0.35,
    "category_id": 1,
    "seller_id": 2,
    "images": ["smartphone.jpg"],
    "has_variants": false,
    "created_at": "2023-04-15T10:00:00Z",
    "updated_at": "2023-04-15T10:00:00Z"
  },
  {
    "id": 2,
    "product_number": "PROD-000002",
    "name": "Laptop",
    "description": "Powerful laptop for professionals",
    "price": 1499.99,
    "stock": 25,
    "weight": 2.1,
    "category_id": 1,
    "seller_id": 2,
    "images": ["laptop.jpg"],
    "has_variants": true,
    "variants": [
      {
        "id": 1,
        "product_id": 2,
        "sku": "LAPT-8GB-256",
        "price": 1499.99,
        "compare_price": 1599.99,
        "stock": 10,
        "weight": 2.1,
        "attributes": {
          "ram": "8GB",
          "storage": "256GB",
          "color": "Silver"
        },
        "images": ["laptop_silver.jpg"],
        "is_default": true,
        "created_at": "2023-04-15T10:00:00Z",
        "updated_at": "2023-04-15T10:00:00Z"
      },
      {
        "id": 2,
        "product_id": 2,
        "sku": "LAPT-16GB-512",
        "price": 1799.99,
        "compare_price": 1899.99,
        "stock": 15,
        "weight": 2.1,
        "attributes": {
          "ram": "16GB",
          "storage": "512GB",
          "color": "Space Gray"
        },
        "images": ["laptop_gray.jpg"],
        "is_default": false,
        "created_at": "2023-04-15T10:00:00Z",
        "updated_at": "2023-04-15T10:00:00Z"
      }
    ],
    "created_at": "2023-04-16T11:00:00Z",
    "updated_at": "2023-04-16T11:00:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Products retrieved successfully

### Get Product

`GET /api/products/{id}`

Get details of a specific product.

Example response:

```json
{
  "id": 2,
  "product_number": "PROD-000002",
  "name": "Laptop",
  "description": "Powerful laptop for professionals",
  "price": 1499.99,
  "stock": 25,
  "weight": 2.1,
  "category_id": 1,
  "seller_id": 2,
  "images": ["laptop.jpg"],
  "has_variants": true,
  "variants": [
    {
      "id": 1,
      "product_id": 2,
      "sku": "LAPT-8GB-256",
      "price": 1499.99,
      "compare_price": 1599.99,
      "stock": 10,
      "weight": 2.1,
      "attributes": {
        "ram": "8GB",
        "storage": "256GB",
        "color": "Silver"
      },
      "images": ["laptop_silver.jpg"],
      "is_default": true,
      "created_at": "2023-04-15T10:00:00Z",
      "updated_at": "2023-04-15T10:00:00Z"
    },
    {
      "id": 2,
      "product_id": 2,
      "sku": "LAPT-16GB-512",
      "price": 1799.99,
      "compare_price": 1899.99,
      "stock": 15,
      "weight": 2.1,
      "attributes": {
        "ram": "16GB",
        "storage": "512GB",
        "color": "Space Gray"
      },
      "images": ["laptop_gray.jpg"],
      "is_default": false,
      "created_at": "2023-04-15T10:00:00Z",
      "updated_at": "2023-04-15T10:00:00Z"
    }
  ],
  "created_at": "2023-04-16T11:00:00Z",
  "updated_at": "2023-04-16T11:00:00Z"
}
```

**Status Codes:**

- `200 OK`: Product retrieved successfully
- `404 Not Found`: Product not found

### Search Products

`GET /api/products/search`

Search products based on various criteria.

**Query Parameters:**

- `q` (optional): Search query
- `category` (optional): Category ID
- `min_price` (optional): Minimum price
- `max_price` (optional): Maximum price
- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

Example response:

```json
[
  {
    "id": 3,
    "product_number": "PROD-000003",
    "name": "T-Shirt",
    "description": "Cotton t-shirt for everyday wear",
    "price": 19.99,
    "stock": 150,
    "weight": 0.2,
    "category_id": 3,
    "seller_id": 3,
    "images": ["tshirt.jpg"],
    "has_variants": true,
    "variants": [
      {
        "id": 5,
        "product_id": 3,
        "sku": "TS-BLU-M",
        "price": 19.99,
        "compare_price": 24.99,
        "stock": 50,
        "weight": 0.2,
        "attributes": {
          "color": "Blue",
          "size": "M"
        },
        "images": ["tshirt_blue.jpg"],
        "is_default": true,
        "created_at": "2023-04-20T10:00:00Z",
        "updated_at": "2023-04-20T10:00:00Z"
      }
    ],
    "created_at": "2023-04-20T10:00:00Z",
    "updated_at": "2023-04-20T10:00:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Search results retrieved successfully

### List Categories

`GET /api/categories`

List all product categories.

Example response:

```json
[
  {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and gadgets",
    "parent_id": null,
    "created_at": "2023-04-10T09:00:00Z",
    "updated_at": "2023-04-10T09:00:00Z"
  },
  {
    "id": 2,
    "name": "Smartphones",
    "description": "Mobile phones and smartphones",
    "parent_id": 1,
    "created_at": "2023-04-10T09:05:00Z",
    "updated_at": "2023-04-10T09:05:00Z"
  },
  {
    "id": 3,
    "name": "Clothing",
    "description": "Apparel and fashion items",
    "parent_id": null,
    "created_at": "2023-04-10T09:10:00Z",
    "updated_at": "2023-04-10T09:10:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Categories retrieved successfully

## Seller Product Endpoints

### Create Product

`POST /api/products`

Create a new product (seller only).

```json
{
  "name": "New Product",
  "description": "Product description",
  "price": 199.99,
  "stock": 100,
  "weight": 1.5,
  "category_id": 1,
  "images": ["product.jpg"],
  "has_variants": false
}
```

Example response:

```json
{
  "id": 4,
  "product_number": "PROD-000004",
  "name": "New Product",
  "description": "Product description",
  "price": 199.99,
  "stock": 100,
  "weight": 1.5,
  "category_id": 1,
  "seller_id": 2,
  "images": ["product.jpg"],
  "has_variants": false,
  "created_at": "2023-04-25T14:00:00Z",
  "updated_at": "2023-04-25T14:00:00Z"
}
```

**Status Codes:**

- `201 Created`: Product created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not a seller)

### Update Product

`PUT /api/products/{id}`

Update an existing product (seller only).

```json
{
  "name": "Updated Product",
  "description": "Updated product description",
  "price": 249.99,
  "stock": 75,
  "weight": 1.6,
  "category_id": 1,
  "images": ["updated-product.jpg"]
}
```

Example response:

```json
{
  "id": 4,
  "product_number": "PROD-000004",
  "name": "Updated Product",
  "description": "Updated product description",
  "price": 249.99,
  "stock": 75,
  "weight": 1.6,
  "category_id": 1,
  "seller_id": 2,
  "images": ["updated-product.jpg"],
  "has_variants": false,
  "created_at": "2023-04-25T14:00:00Z",
  "updated_at": "2023-04-25T14:30:00Z"
}
```

**Status Codes:**

- `200 OK`: Product updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product not found

### Delete Product

`DELETE /api/products/{id}`

Delete a product (seller only).

**Status Codes:**

- `204 No Content`: Product deleted successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product not found

### List Seller Products

`GET /api/products/seller`

List all products for the authenticated seller.

**Query Parameters:**

- `offset` (optional): Pagination offset (default: 0)
- `limit` (optional): Pagination limit (default: 10)

Example response:

```json
[
  {
    "id": 4,
    "product_number": "PROD-000004",
    "name": "Updated Product",
    "description": "Updated product description",
    "price": 249.99,
    "stock": 75,
    "weight": 1.6,
    "category_id": 1,
    "seller_id": 2,
    "images": ["updated-product.jpg"],
    "has_variants": false,
    "created_at": "2023-04-25T14:00:00Z",
    "updated_at": "2023-04-25T14:30:00Z"
  },
  {
    "id": 5,
    "product_number": "PROD-000005",
    "name": "Another Product",
    "description": "Another product description",
    "price": 99.99,
    "stock": 50,
    "weight": 0.8,
    "category_id": 2,
    "seller_id": 2,
    "images": ["another-product.jpg"],
    "has_variants": true,
    "variants": [
      {
        "id": 10,
        "product_id": 5,
        "sku": "AP-RED",
        "price": 99.99,
        "compare_price": 119.99,
        "stock": 25,
        "weight": 0.8,
        "attributes": {
          "color": "Red"
        },
        "images": ["another-product-red.jpg"],
        "is_default": true,
        "created_at": "2023-04-26T10:00:00Z",
        "updated_at": "2023-04-26T10:00:00Z"
      }
    ],
    "created_at": "2023-04-26T10:00:00Z",
    "updated_at": "2023-04-26T10:00:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Products retrieved successfully
- `401 Unauthorized`: Not authenticated

## Product Variant Endpoints

### Add Product Variant

`POST /api/products/{productId}/variants`

Add a variant to a product (seller only).

```json
{
  "sku": "PROD-RED-M",
  "price": 29.99,
  "compare_price": 39.99,
  "stock": 10,
  "weight": 0.3,
  "attributes": {
    "color": "Red",
    "size": "Medium"
  },
  "images": ["red-shirt.jpg"],
  "is_default": true
}
```

Example response:

```json
{
  "id": 11,
  "product_id": 3,
  "sku": "PROD-RED-M",
  "price": 29.99,
  "compare_price": 39.99,
  "stock": 10,
  "weight": 0.3,
  "attributes": {
    "color": "Red",
    "size": "Medium"
  },
  "images": ["red-shirt.jpg"],
  "is_default": true,
  "created_at": "2023-04-28T15:00:00Z",
  "updated_at": "2023-04-28T15:00:00Z"
}
```

**Status Codes:**

- `201 Created`: Variant created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product not found

### Update Product Variant

`PUT /api/products/{productId}/variants/{variantId}`

Update a product variant (seller only).

```json
{
  "sku": "PROD-RED-M",
  "price": 24.99,
  "compare_price": 34.99,
  "stock": 15,
  "weight": 0.3,
  "attributes": {
    "color": "Red",
    "size": "Medium"
  },
  "images": ["red-shirt-updated.jpg"],
  "is_default": true
}
```

Example response:

```json
{
  "id": 11,
  "product_id": 3,
  "sku": "PROD-RED-M",
  "price": 24.99,
  "compare_price": 34.99,
  "stock": 15,
  "weight": 0.3,
  "attributes": {
    "color": "Red",
    "size": "Medium"
  },
  "images": ["red-shirt-updated.jpg"],
  "is_default": true,
  "created_at": "2023-04-28T15:00:00Z",
  "updated_at": "2023-04-28T15:30:00Z"
}
```

**Status Codes:**

- `200 OK`: Variant updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product or variant not found

### Delete Product Variant

`DELETE /api/products/{productId}/variants/{variantId}`

Delete a product variant (seller only).

**Status Codes:**

- `204 No Content`: Variant deleted successfully
- `400 Bad Request`: Cannot delete the only variant of a product
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not the seller of this product)
- `404 Not Found`: Product or variant not found

## Example Workflow

### Product Management Flow (Seller)

1. Seller creates a base product through the seller interface
2. If the product has variants, seller adds variants with different attributes (color, size, etc.)
3. Seller can update product information or variant details as needed
4. Seller can manage inventory levels for products and variants
5. Seller can deactivate or delete products when they're no longer available

### Product Shopping Flow (Customer)

1. Customers browse products by category or use the search function
2. Customers can view detailed product information including available variants
3. When adding to cart, customers select specific variants if the product has them
4. Products and variants are displayed with current inventory levels
5. Out-of-stock products or variants can be marked as unavailable for purchase
