package repository

import (
	"context"
	"errors"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound         = errors.New("product not found")
	ErrProductNameExists       = errors.New("product name already exists")
	ErrProductCategoryNotFound = errors.New("category not found")
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (models.Product, error)
	GetByCategory(ctx context.Context, categoryID int) ([]models.Product, error)
	Create(ctx context.Context, product models.Product) (models.Product, error)
	Update(ctx context.Context, id int, product models.Product) (models.Product, error)
	Delete(ctx context.Context, id int) error
	CategoryExists(ctx context.Context, categoryID int) (bool, error)
}

// productRepository implements ProductRepository using PostgreSQL
type productRepository struct {
	db *pgx.Conn
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *pgx.Conn) ProductRepository {
	return &productRepository{db: db}
}

// GetAll returns all products from the database with their category
func (r *productRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0), c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var catIDFromJoin *int
		var catName, catDesc *string

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID,
			&catIDFromJoin, &catName, &catDesc); err != nil {
			return nil, err
		}

		// Attach category if exists
		if catIDFromJoin != nil && catName != nil {
			p.Category = &models.Category{
				ID:   *catIDFromJoin,
				Name: *catName,
			}
			if catDesc != nil {
				p.Category.Description = *catDesc
			}
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice instead of nil
	if products == nil {
		products = []models.Product{}
	}

	return products, nil
}

// GetByID returns a product by its ID with category
func (r *productRepository) GetByID(ctx context.Context, id int) (models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0),
			   c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	var p models.Product
	var catID *int
	var catName, catDesc *string

	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID,
		&catID, &catName, &catDesc)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, ErrProductNotFound
		}
		return models.Product{}, err
	}

	// Attach category if exists
	if catID != nil && catName != nil {
		p.Category = &models.Category{
			ID:   *catID,
			Name: *catName,
		}
		if catDesc != nil {
			p.Category.Description = *catDesc
		}
	}

	return p, nil
}

// GetByCategory returns all products for a specific category
func (r *productRepository) GetByCategory(ctx context.Context, categoryID int) ([]models.Product, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, COALESCE(p.category_id, 0),
			   c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.category_id = $1
		ORDER BY p.id
	`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var catID *int
		var catName, catDesc *string

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID,
			&catID, &catName, &catDesc); err != nil {
			return nil, err
		}

		// Attach category if exists
		if catID != nil && catName != nil {
			p.Category = &models.Category{
				ID:   *catID,
				Name: *catName,
			}
			if catDesc != nil {
				p.Category.Description = *catDesc
			}
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice instead of nil
	if products == nil {
		products = []models.Product{}
	}

	return products, nil
}

// CategoryExists checks if a category with the given ID exists
func (r *productRepository) CategoryExists(ctx context.Context, categoryID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`
	err := r.db.QueryRow(ctx, query, categoryID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Create adds a new product to the database
func (r *productRepository) Create(ctx context.Context, product models.Product) (models.Product, error) {
	// Check if name already exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM products WHERE name = $1)`
	if err := r.db.QueryRow(ctx, checkQuery, product.Name).Scan(&exists); err != nil {
		return models.Product{}, err
	}
	if exists {
		return models.Product{}, ErrProductNameExists
	}

	// Check if category exists (if specified)
	if product.CategoryID > 0 {
		catExists, err := r.CategoryExists(ctx, product.CategoryID)
		if err != nil {
			return models.Product{}, err
		}
		if !catExists {
			return models.Product{}, ErrProductCategoryNotFound
		}
	}

	// Insert the new product
	var query string
	var err error

	if product.CategoryID > 0 {
		query = `INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id`
		err = r.db.QueryRow(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	} else {
		query = `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`
		err = r.db.QueryRow(ctx, query, product.Name, product.Price, product.Stock).Scan(&product.ID)
	}

	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

// Update updates an existing product
func (r *productRepository) Update(ctx context.Context, id int, product models.Product) (models.Product, error) {
	// Check if category exists (if specified)
	if product.CategoryID > 0 {
		catExists, err := r.CategoryExists(ctx, product.CategoryID)
		if err != nil {
			return models.Product{}, err
		}
		if !catExists {
			return models.Product{}, ErrProductCategoryNotFound
		}
	}

	var query string
	var updated models.Product
	var err error

	if product.CategoryID > 0 {
		query = `UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5 
				 RETURNING id, name, price, stock, COALESCE(category_id, 0)`
		err = r.db.QueryRow(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID, id).
			Scan(&updated.ID, &updated.Name, &updated.Price, &updated.Stock, &updated.CategoryID)
	} else {
		query = `UPDATE products SET name = $1, price = $2, stock = $3, category_id = NULL WHERE id = $4 
				 RETURNING id, name, price, stock, COALESCE(category_id, 0)`
		err = r.db.QueryRow(ctx, query, product.Name, product.Price, product.Stock, id).
			Scan(&updated.ID, &updated.Name, &updated.Price, &updated.Stock, &updated.CategoryID)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, ErrProductNotFound
		}
		return models.Product{}, err
	}

	return updated, nil
}

// Delete removes a product by its ID
func (r *productRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}
