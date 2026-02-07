# Root and Health Endpoints

## Overview
Basic endpoints for API information and health checking.

## Endpoints

### Root Endpoint

#### Get API Information
```
GET /
```

**Response (200 OK):**
```json
{
  "message": "Welcome to BelajarGolang API",
  "version": "1.0.0",
  "endpoints": [
    "GET    /              - API information",
    "GET    /health        - Health check",
    "GET    /categories    - Get all categories",
    "POST   /categories    - Create a category",
    "GET    /categories/{id} - Get a category by ID",
    "GET    /categories?name={name} - Search categories by name",
    "PUT    /categories/{id} - Update a category",
    "DELETE /categories/{id} - Delete a category",
    "GET    /products      - Get all products",
    "POST   /products      - Create a product",
    "GET    /products/{id} - Get a product by ID",
    "GET    /products?name={name} - Search products by name",
    "GET    /products?category_id={id} - Get products by category",
    "PUT    /products/{id} - Update a product",
    "DELETE /products/{id} - Delete a product"
  ]
}
```

### Health Check Endpoint

#### Check Service Health
```
GET /health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2026-02-07T09:20:40.123456789+07:00",
  "service": "BelajarGolang API"
}
```

**Response (405 Method Not Allowed):**
```
Method not allowed
```
Only GET method is supported.

## Features

- **Root endpoint**: Provides API information and available endpoints
- **Health check**: Simple health check for monitoring and load balancers
- **Version information**: Returns current API version
- **Timestamp**: Health check includes server timestamp
- **Method validation**: Health endpoint only accepts GET requests

## Use Cases

### Root Endpoint
- **API Discovery**: New users can discover available endpoints
- **Documentation**: Quick reference for available operations
- **Version Check**: Verify API version

### Health Check
- **Monitoring**: Used by monitoring tools (Prometheus, Datadog, etc.)
- **Load Balancers**: Health checks for load balancer configuration
- **Kubernetes**: Liveness and readiness probes
- **CI/CD**: Verify service is running after deployment

## Examples

### Using curl

**Get API Information:**
```bash
curl http://localhost:8080/
```

**Health Check:**
```bash
curl http://localhost:8080/health
```

### Using httpie

**Get API Information:**
```bash
http GET localhost:8080/
```

**Health Check:**
```bash
http GET localhost:8080/health
```

## Implementation Details

### Root Handler
- Returns API metadata and endpoint list
- Only responds to exact `/` path
- Returns 404 for any other path

### Health Handler
- Returns simple health status
- Includes current timestamp
- Only accepts GET method
- Returns 405 for other HTTP methods

## Testing

Run the tests with:
```bash
go test ./handlers -v -run "Test.*Handler"
```

All tests pass successfully, ensuring the endpoints work as expected.
