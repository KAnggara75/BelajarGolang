package repository

import (
	"context"
	"errors"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrProductNameExists = errors.New("product name already exists")
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (models.Product, error)
	Create(ctx context.Context, product models.Product) (models.Product, error)
	Update(ctx context.Context, id int, product models.Product) (models.Product, error)
	Delete(ctx context.Context, id int) error
}

// productRepository implements ProductRepository using PostgreSQL
type productRepository struct {
	db *pgx.Conn
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *pgx.Conn) ProductRepository {
	return &productRepository{db: db}
}

// GetAll returns all products from the database
func (r *productRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `SELECT id, name, price, stock FROM products ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
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

// GetByID returns a product by its ID
func (r *productRepository) GetByID(ctx context.Context, id int) (models.Product, error) {
	query := `SELECT id, name, price, stock FROM products WHERE id = $1`

	var p models.Product
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, ErrProductNotFound
		}
		return models.Product{}, err
	}

	return p, nil
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

	// Insert the new product
	query := `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(ctx, query, product.Name, product.Price, product.Stock).Scan(&product.ID)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

// Update updates an existing product
func (r *productRepository) Update(ctx context.Context, id int, product models.Product) (models.Product, error) {
	query := `UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4 RETURNING id, name, price, stock`

	var updated models.Product
	err := r.db.QueryRow(ctx, query, product.Name, product.Price, product.Stock, id).Scan(&updated.ID, &updated.Name, &updated.Price, &updated.Stock)
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
