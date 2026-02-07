# BelajarGolang API

[![CI](https://github.com/KAnggara75/BelajarGolang/actions/workflows/ci.yml/badge.svg)](https://github.com/KAnggara75/BelajarGolang/actions/workflows/ci.yml)
[![Go Tests](https://github.com/KAnggara75/BelajarGolang/actions/workflows/test.yml/badge.svg)](https://github.com/KAnggara75/BelajarGolang/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/KAnggara75/BelajarGolang)](https://goreportcard.com/report/github.com/KAnggara75/BelajarGolang)

A RESTful API built with Go for managing categories and products.

## Features

- ✅ **CRUD Operations** for Categories and Products
- ✅ **Transaction System** with atomic checkout operations
- ✅ **Race Condition Prevention** using database row-level locking
- ✅ **Stock Management** with automatic inventory updates
- ✅ **Sales Reports** with daily and period analytics
- ✅ **Revenue Tracking** and best-selling product identification
- ✅ **Search by Name** with partial, case-insensitive matching
- ✅ **Filter Products** by category
- ✅ **Health Check** endpoint for monitoring
- ✅ **API Information** endpoint
- ✅ **PostgreSQL** database with migrations
- ✅ **Comprehensive Tests** with high coverage

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd BelajarGolang
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env and set your DATABASE_URL
```

3. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### General

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API information and available endpoints |
| GET | `/health` | Health check for monitoring |

### Categories

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/categories` | Get all categories |
| GET | `/categories?name={name}` | Search categories by name (partial match) |
| GET | `/categories/{id}` | Get a category by ID |
| POST | `/categories` | Create a new category |
| PUT | `/categories/{id}` | Update a category |
| DELETE | `/categories/{id}` | Delete a category |

### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/products` | Get all products |
| GET | `/products?name={name}` | Search products by name (partial match) |
| GET | `/products?category_id={id}` | Get products by category |
| GET | `/products/{id}` | Get a product by ID |
| POST | `/products` | Create a new product |
| PUT | `/products/{id}` | Update a product |
| DELETE | `/products/{id}` | Delete a product |

### Transactions

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/transactions/checkout` | Create a new transaction (checkout) |
| GET | `/transactions` | Get all transactions |
| GET | `/transactions/{id}` | Get a transaction by ID |

### Reports

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/report/hari-ini` | Get today's sales report |
| GET | `/api/report?start_date={date}&end_date={date}` | Get period sales report |

## Examples

### Get API Information
```bash
curl http://localhost:8080/
```

### Health Check
```bash
curl http://localhost:8080/health
```

### Get All Categories
```bash
curl http://localhost:8080/categories
```

### Search Categories by Name
```bash
curl http://localhost:8080/categories?name=Elect
```

### Get All Products
```bash
curl http://localhost:8080/products
```

### Search Products by Name
```bash
curl http://localhost:8080/products?name=MacBook
```

### Filter Products by Category
```bash
curl http://localhost:8080/products?category_id=1
```

### Create a Category
```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Books","description":"Books and magazines"}'
```

### Create a Product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name":"iPhone 15 Pro",
    "price":999.99,
    "stock":50,
    "category_id":1
  }'
```

### Checkout (Create Transaction)
```bash
curl -X POST http://localhost:8080/transactions/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 3, "quantity": 1}
    ]
  }'
```

### Get All Transactions
```bash
curl http://localhost:8080/transactions
```

### Get Transaction by ID
```bash
curl http://localhost:8080/transactions/1
```

### Get Today's Sales Report
```bash
curl http://localhost:8080/api/report/hari-ini
```

### Get Period Sales Report
```bash
curl "http://localhost:8080/api/report?start_date=2026-01-01&end_date=2026-02-01"
```

## Response Format

All API responses follow this format:

**Success Response:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Error description"
}
```

## Search Feature

The search endpoints support:
- **Partial matching**: "Mac" matches "MacBook Pro"
- **Case-insensitive**: "iphone" matches "iPhone 15 Pro"
- **Multiple results**: Returns all matching items as an array
- **404 on empty**: Returns 404 if no matches found

## Testing

Run all tests:
```bash
go test ./...
```

Run specific test package:
```bash
go test ./handlers -v
go test ./repository -v
```

Run specific test:
```bash
go test ./handlers -v -run TestRootHandler
```

Run tests with coverage:
```bash
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out
```

## CI/CD

This project uses GitHub Actions for continuous integration and testing:

- **CI Workflow** (`.github/workflows/ci.yml`): Runs on every push and pull request
  - Tests against Go 1.21 and 1.22
  - Runs all tests with race detection
  - Builds the application

- **Test Workflow** (`.github/workflows/test.yml`): Comprehensive testing
  - Runs all tests with coverage reporting
  - Uploads coverage to Codecov
  - Runs `go vet` for static analysis
  - Runs golangci-lint for code quality

### Running Linter Locally

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## Project Structure

```
BelajarGolang/
├── config/          # Configuration management
├── database/        # Database connection and migrations
├── handlers/        # HTTP request handlers
├── models/          # Data models
├── repository/      # Data access layer
├── main.go          # Application entry point
└── README.md        # This file
```

## Documentation

- [Search by Name Feature](SEARCH_BY_NAME.md)
- [Root and Health Endpoints](ROOT_AND_HEALTH.md)
- [Transaction System](TRANSACTIONS.md)
- [Sales Reports](REPORTS.md)
- [GitHub Actions CI/CD](GITHUB_ACTIONS.md)

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `PORT` | Server port | `:8080` |

## Database Schema

### Categories Table
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);
```

### Products Table
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    category_id INTEGER REFERENCES categories(id)
);
```

### Transactions Table
```sql
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Transaction Details Table
```sql
CREATE TABLE transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    subtotal INT NOT NULL
);
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
