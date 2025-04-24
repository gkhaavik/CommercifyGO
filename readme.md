# Commercify Backend API

A robust, scalable e-commerce backend API built with Go, following clean architecture principles and best practices.

## Features

- **User Management**: Registration, authentication, profile management
- **Product Management**: CRUD operations, categories, variants, search
- **Shopping Cart**: Add, update, remove items
- **Order Processing**: Create orders, payment processing, order status tracking
- **Payment Integration**: Support for multiple payment providers (Stripe, PayPal, etc.)
- **Email Notifications**: Order confirmations, status updates

## Technology Stack

- **Language**: Go 1.20+
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Payment Processing**: Stripe, PayPal (configurable)
- **Email**: SMTP integration

## Project Structure

The project follows clean architecture principles with clear separation of concerns:

```
├── cmd/ # Application entry points
│ ├── api/ # API server
│ ├── migrate/ # Database migration tool
│ └── seed/ # Database seeding tool
├── config/ # Configuration
├── internal/ # Internal packages
│ ├── api/ # API layer (handlers, middleware, server)
│ ├── application/ # Application layer (use cases)
│ ├── domain/ # Domain layer (entities, repositories interfaces)
│ └── infrastructure/ # Infrastructure layer (repositories implementation, services)
├── migrations/ # Database migrations
├── templates/ # Email templates
└── testutil/ # Testing utilities
```

## Setup and Installation

### Prerequisites

- Go 1.20+
- PostgreSQL 15 (Only tested on v15)
- Docker (optional)

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Server

SERVER_PORT=8080
SERVER_READ_TIMEOUT=15
SERVER_WRITE_TIMEOUT=15

# Database

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ecommerce
DB_SSL_MODE=disable

# Authentication

AUTH_JWT_SECRET=your-secret-key
AUTH_TOKEN_DURATION=24

# Email

EMAIL_SMTP_HOST=smtp.example.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=your-username
EMAIL_SMTP_PASSWORD=your-password
EMAIL_FROM_ADDRESS=noreply@example.com
EMAIL_FROM_NAME=E-Commerce Store
EMAIL_ADMIN_ADDRESS=admin@example.com
EMAIL_ENABLED=false

# Payment - Stripe

STRIPE_SECRET_KEY=your-stripe-secret-key
STRIPE_PUBLIC_KEY=your-stripe-public-key
STRIPE_WEBHOOK_SECRET=your-stripe-webhook-secret
STRIPE_PAYMENT_DESCRIPTION=E-Commerce Store Purchase
STRIPE_ENABLED=false

# Payment - PayPal

PAYPAL_CLIENT_ID=your-paypal-client-id
PAYPAL_CLIENT_SECRET=your-paypal-client-secret
PAYPAL_SANDBOX=true
PAYPAL_ENABLED=false
```

### Database Setup

1. Create a PostgreSQL user (optional):

```bash
createuser -s newuser
```

2. Create a PostgreSQL database:

```bash
createdb -U newuser commercify
```

3. Run migrations:

```bash
go run cmd/migrate/main.go -up
```

4. Seed the database with sample data (optional):

```bash
go run cmd/seed/main.go -all
```

### Running the Application

# Build the application

```
go build -o commercify cmd/api/main.go
```

# Run the application

```bash
./commercify
```

Or simply:

```bash
go run cmd/api/main.go
```

## API Documentation

### Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

### API Endpoints

#### Users

- `POST /api/users/register` - Register a new user
- `POST /api/users/login` - Login and get JWT token
- `GET /api/users/me` - Get current user profile
- `PUT /api/users/me` - Update user profile
- `PUT /api/users/me/password` - Change password

#### Products

- `GET /api/products` - List products
- `GET /api/products/{id}` - Get product details
- `GET /api/products/search` - Search products
- `POST /api/products` - Create product (seller only)
- `PUT /api/products/{id}` - Update product (seller only)
- `DELETE /api/products/{id}` - Delete product (seller only)

#### Product Variants

- `POST /api/products/{productId}/variants` - Add variant (seller only)
- `PUT /api/products/{productId}/variants/{variantId}` - Update variant (seller only)
- `DELETE /api/products/{productId}/variants/{variantId}` - Delete variant (seller only)

#### Categories

- `GET /api/categories` - List categories

#### Cart

- `GET /api/cart` - Get user's cart
- `POST /api/cart/items` - Add item to cart
- `PUT /api/cart/items/{productId}` - Update cart item
- `DELETE /api/cart/items/{productId}` - Remove item from cart
- `DELETE /api/cart` - Clear cart

#### Orders

- `POST /api/orders` - Create order from cart
- `GET /api/orders/{id}` - Get order details
- `GET /api/orders` - List user's orders
- `POST /api/orders/{id}/payment` - Process payment for order

#### Payment

- `GET /api/payment/providers` - Get available payment providers

#### Admin

- `GET /api/admin/users` - List all users (admin only)
- `GET /api/admin/orders` - List all orders (admin only)
- `PUT /api/admin/orders/{id}/status` - Update order status (admin only)

### Webhooks

- `POST /api/webhooks/stripe` - Stripe webhook endpoint

## Database Schema

The database consists of the following main tables:

- `users` - User accounts
- `categories` - Product categories
- `products` - Products information
- `product_variants` - Product variants
- `carts` - Shopping carts
- `cart_items` - Items in shopping carts
- `orders` - Customer orders
- `order_items` - Items in orders

## Development

### Running Tests

# Run all tests

```bash
go test ./...
```

# Run tests with coverage

```bash
go test -cover ./...
```

### Adding Migrations

# Create a new migration

Install the [golang-migrate/migrate](https://github.com/golang-migrate/migrate) tool

### Homebrew

```bash
brew install migrate
```

using the cli tool to generate migration files, then you can use the following command, which creates both the files in the right format:

```bash
migrate create -ext sql -dir migrations -seq add_friendly_numbers
```

otherwise you can create them manually wher `migrations` is the migrations folder, the `sequence` is the 6 digits in front and `migration_name` is a short description

```bash
touch migrations/[sequence]_[migration_name].up.sql
touch migrations/[sequence]_[migration_name].down.sql
```

## Multi-Provider Payment System

The application supports multiple payment providers through a flexible payment service architecture:

- **Stripe**: Credit card payments
- **PayPal**: PayPal account payments
- **Mock**: Test payment provider for development

Payment providers can be enabled/disabled through configuration, and new providers can be added by implementing the `PaymentService` interface.

## User-Friendly Identifiers

The system uses user-friendly identifiers for better readability:

- **Order Numbers**: Format `ORD-YYYYMMDD-000001` (date-based with sequential numbering)
- **Product Numbers**: Format `PROD-000001` (sequential numbering)

These identifiers make it easier to reference orders and products in the UI and customer communications.
