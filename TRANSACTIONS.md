# Transaction Feature

## Overview
The transaction feature provides a complete checkout system with atomic operations to prevent race conditions. It handles product stock management and transaction recording with full ACID compliance.

## Features

- ✅ **Atomic Checkout**: All operations in a transaction succeed or fail together
- ✅ **Race Condition Prevention**: Uses database row-level locking (`FOR UPDATE`)
- ✅ **Stock Management**: Automatically decrements product stock on checkout
- ✅ **Transaction History**: Complete record of all transactions with details
- ✅ **Validation**: Comprehensive validation for quantities and stock availability
- ✅ **Error Handling**: Clear error messages for all failure scenarios

## Database Schema

### Transactions Table
```sql
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Transaction Details Table
```sql
CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    subtotal INT NOT NULL
);
```

## API Endpoints

### 1. Checkout (Create Transaction)

**Endpoint**: `POST /transactions/checkout`

**Description**: Process a checkout request, decrement product stock, and create a transaction record.

**Request Body**:
```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 3,
      "quantity": 1
    }
  ]
}
```

**Success Response (201 Created)**:
```json
{
  "success": true,
  "message": "Transaction created successfully",
  "data": {
    "id": 1,
    "total_amount": 224997,
    "created_at": "2026-02-07T09:31:36+07:00",
    "details": [
      {
        "id": 1,
        "transaction_id": 1,
        "product_id": 1,
        "product_name": "iPhone 15 Pro",
        "quantity": 2,
        "subtotal": 199998
      },
      {
        "id": 2,
        "transaction_id": 1,
        "product_id": 3,
        "product_name": "AirPods Pro",
        "quantity": 1,
        "subtotal": 24999
      }
    ]
  }
}
```

**Error Responses**:

- **400 Bad Request** - Empty checkout items:
```json
{
  "success": false,
  "message": "Checkout items cannot be empty"
}
```

- **400 Bad Request** - Invalid quantity:
```json
{
  "success": false,
  "message": "Invalid quantity: must be greater than 0"
}
```

- **400 Bad Request** - Insufficient stock:
```json
{
  "success": false,
  "message": "insufficient stock: product iPhone 15 Pro has only 5 in stock, requested 10"
}
```

- **404 Not Found** - Product not found:
```json
{
  "success": false,
  "message": "One or more products not found"
}
```

### 2. Get All Transactions

**Endpoint**: `GET /transactions`

**Description**: Retrieve all transactions with their details.

**Success Response (200 OK)**:
```json
{
  "success": true,
  "message": "Transactions retrieved successfully",
  "data": [
    {
      "id": 1,
      "total_amount": 224997,
      "created_at": "2026-02-07T09:31:36+07:00",
      "details": [
        {
          "id": 1,
          "transaction_id": 1,
          "product_id": 1,
          "product_name": "iPhone 15 Pro",
          "quantity": 2,
          "subtotal": 199998
        }
      ]
    }
  ]
}
```

### 3. Get Transaction by ID

**Endpoint**: `GET /transactions/{id}`

**Description**: Retrieve a specific transaction by ID.

**Success Response (200 OK)**:
```json
{
  "success": true,
  "message": "Transaction retrieved successfully",
  "data": {
    "id": 1,
    "total_amount": 224997,
    "created_at": "2026-02-07T09:31:36+07:00",
    "details": [...]
  }
}
```

**Error Response (404 Not Found)**:
```json
{
  "success": false,
  "message": "Transaction not found"
}
```

## Race Condition Prevention

The checkout process uses PostgreSQL's row-level locking to prevent race conditions:

```go
// Lock the product row for update
query := `SELECT id, name, price, stock FROM products WHERE id = $1 FOR UPDATE`
```

**How it works**:
1. Transaction begins
2. Product rows are locked with `FOR UPDATE`
3. Stock is checked
4. Stock is decremented
5. Transaction record is created
6. Transaction details are created
7. Transaction is committed (or rolled back on error)

**Benefits**:
- Multiple concurrent checkouts won't cause overselling
- Stock is always accurate
- All operations are atomic (all succeed or all fail)

## Price Handling

Prices are stored as integers (cents) to avoid floating-point precision issues:
- Database stores: `999.99` → `99999` cents
- API returns: `99999` cents
- Frontend should divide by 100 to display: `$999.99`

## Examples

### Checkout with curl

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

### Get all transactions

```bash
curl http://localhost:8080/transactions
```

### Get specific transaction

```bash
curl http://localhost:8080/transactions/1
```

## Testing

Run transaction tests:
```bash
go test ./handlers -v -run TestCheckout
go test ./handlers -v -run "Test.*Transaction"
```

All tests include:
- ✅ Successful checkout
- ✅ Empty items validation
- ✅ Invalid quantity validation
- ✅ Insufficient stock handling
- ✅ Product not found handling
- ✅ Transaction retrieval

## Error Handling

The transaction feature includes comprehensive error handling:

| Error | HTTP Status | Description |
|-------|-------------|-------------|
| `ErrEmptyCheckout` | 400 | No items in checkout request |
| `ErrInvalidQuantity` | 400 | Quantity is 0 or negative |
| `ErrInsufficientStock` | 400 | Not enough stock available |
| `ErrProductNotFoundInTx` | 404 | Product doesn't exist |
| `ErrTransactionNotFound` | 404 | Transaction doesn't exist |

## Implementation Details

### Repository Layer
- **Database Transactions**: Uses `pgx.Tx` for atomic operations
- **Row Locking**: `FOR UPDATE` prevents concurrent modifications
- **Rollback on Error**: Automatic rollback if any step fails
- **Commit on Success**: All changes committed together

### Handler Layer
- **Request Validation**: Validates checkout request structure
- **Error Translation**: Converts repository errors to HTTP responses
- **Response Formatting**: Consistent JSON response format

## Best Practices

1. **Always use transactions** for operations that modify multiple rows
2. **Lock rows** when reading data that will be modified
3. **Validate input** before starting database transaction
4. **Handle all error cases** explicitly
5. **Test concurrent scenarios** to ensure race condition prevention works

## Performance Considerations

- Row-level locking is efficient for low-to-medium concurrency
- For high concurrency, consider optimistic locking or queue-based processing
- Indexes on `transaction_id` and `product_id` improve query performance
- Keep transactions short to minimize lock contention

## Future Enhancements

Potential improvements:
- [ ] Add transaction cancellation/refund
- [ ] Add payment integration
- [ ] Add transaction status (pending, completed, cancelled)
- [ ] Add discount/coupon support
- [ ] Add transaction export (CSV, PDF)
- [ ] Add analytics and reporting
