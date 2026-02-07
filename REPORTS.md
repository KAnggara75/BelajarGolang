# Sales Report Feature

## Overview
The sales report feature provides comprehensive analytics on transactions, revenue, and product performance. It supports both daily reports and custom period reports.

## Features

- ✅ **Daily Reports**: Get today's sales summary
- ✅ **Period Reports**: Get sales summary for any date range
- ✅ **Revenue Tracking**: Total revenue in the period
- ✅ **Transaction Count**: Number of transactions
- ✅ **Top Product**: Best-selling product with quantity sold
- ✅ **Date Validation**: Ensures valid date formats and ranges

## API Endpoints

### 1. Get Today's Report

**Endpoint**: `GET /api/report/hari-ini`

**Description**: Returns sales report for today (current date).

**Success Response (200 OK)**:
```json
{
  "total_revenue": 45000,
  "total_transaksi": 5,
  "produk_terlaris": {
    "nama": "Indomie Goreng",
    "qty_terjual": 12
  }
}
```

**Response with No Transactions**:
```json
{
  "total_revenue": 0,
  "total_transaksi": 0,
  "produk_terlaris": null
}
```

### 2. Get Period Report

**Endpoint**: `GET /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD`

**Description**: Returns sales report for a specific date range.

**Query Parameters**:
- `start_date` (required): Start date in YYYY-MM-DD format (e.g., 2026-01-01)
- `end_date` (required): End date in YYYY-MM-DD format (e.g., 2026-02-01)

**Success Response (200 OK)**:
```json
{
  "total_revenue": 150000,
  "total_transaksi": 15,
  "produk_terlaris": {
    "nama": "Indomie Goreng",
    "qty_terjual": 45
  }
}
```

**Error Responses**:

- **400 Bad Request** - Missing parameters:
```json
{
  "success": false,
  "message": "Missing start_date or end_date query parameters"
}
```

- **400 Bad Request** - Invalid date format:
```json
{
  "success": false,
  "message": "Invalid start_date format. Use YYYY-MM-DD"
}
```

- **400 Bad Request** - Invalid date range:
```json
{
  "success": false,
  "message": "end_date must be after start_date"
}
```

## Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `total_revenue` | integer | Total revenue in cents (e.g., 45000 = $450.00) |
| `total_transaksi` | integer | Number of transactions in the period |
| `produk_terlaris` | object/null | Best-selling product (null if no transactions) |
| `produk_terlaris.nama` | string | Product name |
| `produk_terlaris.qty_terjual` | integer | Total quantity sold |

## Examples

### Get Today's Report

```bash
curl http://localhost:8080/api/report/hari-ini
```

**Response**:
```json
{
  "total_revenue": 45000,
  "total_transaksi": 5,
  "produk_terlaris": {
    "nama": "Indomie Goreng",
    "qty_terjual": 12
  }
}
```

### Get January 2026 Report

```bash
curl "http://localhost:8080/api/report?start_date=2026-01-01&end_date=2026-01-31"
```

**Response**:
```json
{
  "total_revenue": 150000,
  "total_transaksi": 15,
  "produk_terlaris": {
    "nama": "MacBook Pro M3",
    "qty_terjual": 8
  }
}
```

### Get Year-to-Date Report

```bash
curl "http://localhost:8080/api/report?start_date=2026-01-01&end_date=2026-12-31"
```

### Get Last 7 Days Report

```bash
# Using date command to calculate dates
START_DATE=$(date -v-7d +%Y-%m-%d)
END_DATE=$(date +%Y-%m-%d)
curl "http://localhost:8080/api/report?start_date=$START_DATE&end_date=$END_DATE"
```

## Implementation Details

### Repository Layer

The report repository uses SQL aggregation queries:

**Revenue and Transaction Count**:
```sql
SELECT 
    COALESCE(SUM(total_amount), 0) as total_revenue,
    COUNT(*) as total_transaksi
FROM transactions
WHERE created_at >= $1 AND created_at <= $2
```

**Top Product**:
```sql
SELECT 
    p.name,
    SUM(td.quantity) as qty_terjual
FROM transaction_details td
JOIN transactions t ON td.transaction_id = t.id
JOIN products p ON td.product_id = p.id
WHERE t.created_at >= $1 AND t.created_at <= $2
GROUP BY p.id, p.name
ORDER BY qty_terjual DESC
LIMIT 1
```

### Date Handling

- **Today's Report**: Uses current date with time set to 00:00:00 - 23:59:59
- **Period Report**: Includes full days (start at 00:00:00, end at 23:59:59)
- **Timezone**: Uses server's local timezone
- **Date Format**: ISO 8601 format (YYYY-MM-DD)

### Handler Layer

- **Date Parsing**: Uses `time.Parse("2006-01-02", dateStr)`
- **Validation**: Checks date format and ensures end_date >= start_date
- **Error Handling**: Returns appropriate HTTP status codes and messages

## Use Cases

### Business Analytics
- **Daily Sales Tracking**: Monitor today's performance
- **Period Comparison**: Compare different time periods
- **Product Performance**: Identify best-selling products
- **Revenue Trends**: Track revenue over time

### Reporting
- **Daily Reports**: End-of-day sales summary
- **Weekly Reports**: Weekly performance review
- **Monthly Reports**: Monthly business analysis
- **Quarterly Reports**: Quarterly financial reporting

### Dashboard Integration
- **Real-time Metrics**: Display current day's performance
- **Historical Data**: Show trends and patterns
- **Product Analytics**: Track product popularity

## Testing

Run report tests:
```bash
go test ./handlers -v -run "Test.*Report"
```

All tests include:
- ✅ Today's report success
- ✅ Period report success
- ✅ Missing parameters validation
- ✅ Invalid date format handling
- ✅ Invalid date range handling
- ✅ Method not allowed
- ✅ Not found handling

## Performance Considerations

- **Indexes**: Ensure indexes on `transactions.created_at` for fast date filtering
- **Aggregation**: SQL aggregation is efficient for large datasets
- **Caching**: Consider caching daily reports (they don't change once the day is over)
- **Date Range**: Very large date ranges may be slow; consider pagination or limits

## Future Enhancements

Potential improvements:
- [ ] Add hourly breakdown for today's report
- [ ] Add product category breakdown
- [ ] Add revenue by payment method
- [ ] Add customer analytics (if customer data is added)
- [ ] Add export to CSV/PDF
- [ ] Add chart data for visualization
- [ ] Add comparison with previous period
- [ ] Add caching for historical reports
- [ ] Add real-time WebSocket updates for today's report

## Best Practices

1. **Always validate date inputs** before querying
2. **Use appropriate date ranges** to avoid performance issues
3. **Cache historical reports** since they don't change
4. **Index date columns** for better query performance
5. **Consider timezone** when displaying reports to users
6. **Handle edge cases** like no transactions gracefully
7. **Monitor query performance** for large datasets

## Error Handling

| Error | HTTP Status | Description |
|-------|-------------|-------------|
| Missing parameters | 400 | start_date or end_date not provided |
| Invalid date format | 400 | Date not in YYYY-MM-DD format |
| Invalid date range | 400 | end_date before start_date |
| Method not allowed | 405 | Non-GET request |
| Internal error | 500 | Database or server error |

## Integration Examples

### JavaScript/Fetch
```javascript
// Get today's report
fetch('http://localhost:8080/api/report/hari-ini')
  .then(res => res.json())
  .then(data => {
    console.log('Revenue:', data.total_revenue);
    console.log('Transactions:', data.total_transaksi);
    console.log('Top Product:', data.produk_terlaris?.nama);
  });

// Get period report
const startDate = '2026-01-01';
const endDate = '2026-01-31';
fetch(`http://localhost:8080/api/report?start_date=${startDate}&end_date=${endDate}`)
  .then(res => res.json())
  .then(data => console.log(data));
```

### Python/Requests
```python
import requests

# Get today's report
response = requests.get('http://localhost:8080/api/report/hari-ini')
data = response.json()
print(f"Revenue: {data['total_revenue']}")
print(f"Transactions: {data['total_transaksi']}")

# Get period report
params = {
    'start_date': '2026-01-01',
    'end_date': '2026-01-31'
}
response = requests.get('http://localhost:8080/api/report', params=params)
data = response.json()
print(data)
```
