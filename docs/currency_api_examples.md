# Currency API Examples

This document provides example request bodies for the currency system API endpoints.

## Public Currency Endpoints

### Get Enabled Currencies

```plaintext
GET /api/currencies
```

Retrieves all currently enabled currencies.

Example response:

```json
[
  {
    "code": "USD",
    "name": "US Dollar",
    "symbol": "$",
    "precision": 2,
    "exchange_rate": 1.0,
    "is_default": true,
    "is_enabled": true,
    "formatted_name": "US Dollar ($)"
  },
  {
    "code": "EUR",
    "name": "Euro",
    "symbol": "€",
    "precision": 2,
    "exchange_rate": 0.85,
    "is_default": false,
    "is_enabled": true,
    "formatted_name": "Euro (€)"
  },
  {
    "code": "GBP",
    "name": "British Pound",
    "symbol": "£",
    "precision": 2,
    "exchange_rate": 0.75,
    "is_default": false,
    "is_enabled": true,
    "formatted_name": "British Pound (£)"
  },
  {
    "code": "JPY",
    "name": "Japanese Yen",
    "symbol": "¥",
    "precision": 0,
    "exchange_rate": 110.0,
    "is_default": false,
    "is_enabled": true,
    "formatted_name": "Japanese Yen (¥)"
  }
]
```

**Status Codes:**

- `200 OK`: Currencies retrieved successfully
- `500 Internal Server Error`: Failed to retrieve currencies

### Get Default Currency

```plaintext
GET /api/currencies/default
```

Retrieves the store's default currency.

Example response:

```json
{
  "code": "USD",
  "name": "US Dollar",
  "symbol": "$",
  "precision": 2,
  "exchange_rate": 1.0,
  "is_default": true,
  "is_enabled": true,
  "formatted_name": "US Dollar ($)"
}
```

**Status Codes:**

- `200 OK`: Default currency retrieved successfully
- `500 Internal Server Error`: Failed to retrieve default currency

### Get Currency by Code

```plaintext
GET /api/currencies/{code}
```

Retrieves a specific currency by its ISO code.

Example response:

```json
{
  "code": "EUR",
  "name": "Euro",
  "symbol": "€",
  "precision": 2,
  "exchange_rate": 0.85,
  "is_default": false,
  "is_enabled": true,
  "formatted_name": "Euro (€)"
}
```

**Status Codes:**

- `200 OK`: Currency retrieved successfully
- `404 Not Found`: Currency not found
- `500 Internal Server Error`: Failed to retrieve currency

### Convert Currency

```plaintext
POST /api/currencies/convert
```

Converts an amount from one currency to another.

**Request Body:**

```json
{
  "amount": 100.00,
  "from_currency": "USD",
  "to_currency": "EUR"
}
```

Example response:

```json
{
  "original_amount": 100.00,
  "converted_amount": 85.00,
  "from_currency": "USD",
  "to_currency": "EUR",
  "rate": 0.85,
  "formatted_original": "$100.00",
  "formatted_converted": "€85.00"
}
```

**Status Codes:**

- `200 OK`: Conversion performed successfully
- `400 Bad Request`: Invalid request or unsupported currency
- `500 Internal Server Error`: Failed to perform conversion

## Admin Currency Endpoints

### Get All Currencies

```plaintext
GET /api/admin/currencies
```

Retrieves all currencies in the system, including disabled ones (admin only).

Example response:

```json
[
  {
    "code": "USD",
    "name": "US Dollar",
    "symbol": "$",
    "precision": 2,
    "exchange_rate": 1.0,
    "is_default": true,
    "is_enabled": true,
    "formatted_name": "US Dollar ($)"
  },
  {
    "code": "EUR",
    "name": "Euro",
    "symbol": "€",
    "precision": 2,
    "exchange_rate": 0.85,
    "is_default": false,
    "is_enabled": true,
    "formatted_name": "Euro (€)"
  },
  {
    "code": "GBP",
    "name": "British Pound",
    "symbol": "£",
    "precision": 2,
    "exchange_rate": 0.75,
    "is_default": false,
    "is_enabled": true,
    "formatted_name": "British Pound (£)"
  },
  {
    "code": "CAD",
    "name": "Canadian Dollar",
    "symbol": "$",
    "precision": 2,
    "exchange_rate": 1.25,
    "is_default": false,
    "is_enabled": false,
    "formatted_name": "Canadian Dollar ($)"
  }
]
```

**Status Codes:**

- `200 OK`: Currencies retrieved successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `500 Internal Server Error`: Failed to retrieve currencies

### Create Currency

```plaintext
POST /api/admin/currencies
```

Creates a new currency (admin only).

**Request Body:**

```json
{
  "code": "AUD",
  "name": "Australian Dollar",
  "symbol": "$",
  "precision": 2,
  "rate": 1.35,
  "is_default": false,
  "is_enabled": true
}
```

Example response:

```json
{
  "code": "AUD",
  "name": "Australian Dollar",
  "symbol": "$",
  "precision": 2,
  "exchange_rate": 1.35,
  "is_default": false,
  "is_enabled": true,
  "formatted_name": "Australian Dollar ($)"
}
```

**Status Codes:**

- `201 Created`: Currency created successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `409 Conflict`: Currency code already exists
- `500 Internal Server Error`: Failed to create currency

### Update Currency

```plaintext
PUT /api/admin/currencies/{code}
```

Updates an existing currency (admin only).

**Request Body:**

```json
{
  "name": "Australian Dollar",
  "symbol": "$",
  "precision": 2,
  "rate": 1.38,
  "is_default": false,
  "is_enabled": true
}
```

Example response:

```json
{
  "code": "AUD",
  "name": "Australian Dollar",
  "symbol": "$",
  "precision": 2,
  "exchange_rate": 1.38,
  "is_default": false,
  "is_enabled": true,
  "formatted_name": "Australian Dollar ($)"
}
```

**Status Codes:**

- `200 OK`: Currency updated successfully
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Currency not found
- `500 Internal Server Error`: Failed to update currency

### Delete Currency

```plaintext
DELETE /api/admin/currencies/{code}
```

Deletes a currency (admin only). The default currency cannot be deleted.

**Status Codes:**

- `200 OK`: Currency deleted successfully
- `400 Bad Request`: Cannot delete default currency
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Currency not found
- `500 Internal Server Error`: Failed to delete currency

### Set Default Currency

```plaintext
POST /api/admin/currencies/{code}/default
```

Sets a currency as the default for the store (admin only).

Example response:

```json
{
  "message": "Default currency updated successfully"
}
```

**Status Codes:**

- `200 OK`: Default currency set successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `404 Not Found`: Currency not found
- `500 Internal Server Error`: Failed to set default currency

### Update Exchange Rates

```plaintext
POST /api/admin/currencies/update-rates
```

Updates all exchange rates from the provider (admin only).

Example response:

```json
{
  "message": "Exchange rates updated successfully"
}
```

**Status Codes:**

- `200 OK`: Exchange rates updated successfully
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `500 Internal Server Error`: Failed to update exchange rates

### Get Exchange Rate History

```plaintext
GET /api/admin/currencies/rates?base_currency=USD&target_currency=EUR&limit=5
```

Retrieves exchange rate history for a currency pair (admin only).

**Query Parameters:**

- `base_currency` (required): Base currency code
- `target_currency` (required): Target currency code
- `limit` (optional): Maximum number of records to return (default: 10)

Example response:

```json
[
  {
    "base_currency": "USD",
    "target_currency": "EUR",
    "rate": 0.85,
    "date": "2023-06-28T10:00:00Z"
  },
  {
    "base_currency": "USD",
    "target_currency": "EUR",
    "rate": 0.84,
    "date": "2023-06-27T10:00:00Z"
  },
  {
    "base_currency": "USD",
    "target_currency": "EUR",
    "rate": 0.84,
    "date": "2023-06-26T10:00:00Z"
  },
  {
    "base_currency": "USD",
    "target_currency": "EUR",
    "rate": 0.83,
    "date": "2023-06-25T10:00:00Z"
  },
  {
    "base_currency": "USD",
    "target_currency": "EUR",
    "rate": 0.83,
    "date": "2023-06-24T10:00:00Z"
  }
]
```

**Status Codes:**

- `200 OK`: Exchange rate history retrieved successfully
- `400 Bad Request`: Missing required parameters
- `401 Unauthorized`: Not authenticated
- `403 Forbidden`: Not authorized (not an admin)
- `500 Internal Server Error`: Failed to retrieve exchange rate history

## Example Workflows

### Currency Configuration Flow (Admin)

1. Admin views all currencies in the system
2. Admin adds new currencies as needed for their store
3. Admin configures which currencies are enabled for customers to use
4. Admin sets the default store currency
5. Admin updates exchange rates regularly to reflect current market rates
6. Admin can view historical exchange rates for analysis

### Customer Currency Selection Flow

1. Customer views products with prices in the default currency
2. Customer can select an alternative currency from the enabled list
3. System converts all prices to the selected currency using current exchange rates
4. Customer can switch between currencies at any time during their shopping session
5. Final checkout and payment may occur in the default or selected currency depending on system configuration
