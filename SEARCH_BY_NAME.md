# Search by Name Feature

## Overview
Added search by name functionality for both categories and products endpoints with **partial, case-insensitive matching** using PostgreSQL's `ILIKE` operator. Returns **all matching items** as an array.

## Endpoints

### Categories

#### Search by Name (Partial Match)
```
GET /categories?name=Elect
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Electronics",
      "description": "Electronic devices and gadgets"
    }
  ]
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Category not found"
}
```

**Example with URL encoding:**
```
GET /categories?name=Food
```
This will match "Food & Beverages"

### Products

#### Search by Name (Partial Match)
```
GET /products?name=MacBook
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": 2,
      "name": "MacBook Pro M3",
      "price": 2499.99,
      "stock": 25,
      "category_id": 1,
      "category": {
        "id": 1,
        "name": "Electronics",
        "description": "Electronic devices"
      }
    }
  ]
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Product not found"
}
```

## Features

- **Partial matching**: Search term can appear anywhere in the name (e.g., "Book" will match "MacBook Pro M3")
- **Case-insensitive search**: Uses PostgreSQL's `ILIKE` operator (e.g., "macbook" will match "MacBook")
- **URL encoding support**: Special characters and spaces should be URL-encoded
- **404 on not found**: Returns HTTP 404 status code when no matches are found
- **Category inclusion**: Product search results include the associated category information
- **All matches returned**: Returns an array of all matching items, not just the first one
- **Empty array handling**: Returns 404 if the search yields no results

## Search Examples

### Categories
- `?name=Elect` → matches "Electronics"
- `?name=book` → matches "Books" (case-insensitive)
- `?name=food` → matches "Food & Beverages"
- `?name=sport` → matches "Sports"

### Products
- `?name=MacBook` → matches "MacBook Pro M3"
- `?name=iphone` → matches "iPhone 15 Pro" (case-insensitive)
- `?name=Pro` → matches "iPhone 15 Pro", "MacBook Pro M3", "AirPods Pro", etc.
- `?name=Air` → matches "iPad Air", "AirPods Pro"

## Implementation Details

### Repository Layer
- Added `GetByName(ctx context.Context, name string)` method to both `CategoryRepository` and `ProductRepository` interfaces
- Returns `[]models.Category` and `[]models.Product` respectively (arrays, not single items)
- Implemented partial, case-insensitive search using PostgreSQL's `ILIKE '%' || $1 || '%'` pattern
- Returns all matching records ordered by ID

### Handler Layer
- Updated `ServeHTTP` methods to check for `name` query parameter
- Added `GetByName` handler methods for both categories and products
- Returns 404 status code when no items are found (empty array)
- Returns array of matching items in the response

## Testing

Run the tests with:
```bash
go test ./handlers -v -run "TestGet.*ByName"
```

All existing tests continue to pass, ensuring backward compatibility.

## SQL Query Pattern

The search uses the following SQL pattern:
```sql
-- Categories
SELECT id, name, description 
FROM categories 
WHERE name ILIKE '%' || $1 || '%'

-- Products
SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0),
       c.id, c.name, c.description
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.name ILIKE '%' || $1 || '%'
```

The `ILIKE` operator provides case-insensitive matching, and the `'%' || $1 || '%'` pattern allows the search term to appear anywhere in the name.
