package repository

import (
	"context"
	"errors"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound   = errors.New("category not found")
	ErrNameExists = errors.New("category name already exists")
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (models.Category, error)
	Create(ctx context.Context, cat models.Category) (models.Category, error)
	Update(ctx context.Context, id int, cat models.Category) (models.Category, error)
	Delete(ctx context.Context, id int) error
}

// categoryRepository implements CategoryRepository using PostgreSQL
type categoryRepository struct {
	db *pgx.Conn
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *pgx.Conn) CategoryRepository {
	return &categoryRepository{db: db}
}

// GetAll returns all categories from the database
func (r *categoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	query := `SELECT id, name, description FROM categories ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice instead of nil
	if categories == nil {
		categories = []models.Category{}
	}

	return categories, nil
}

// GetByID returns a category by its ID
func (r *categoryRepository) GetByID(ctx context.Context, id int) (models.Category, error) {
	query := `SELECT id, name, description FROM categories WHERE id = $1`

	var cat models.Category
	err := r.db.QueryRow(ctx, query, id).Scan(&cat.ID, &cat.Name, &cat.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Category{}, ErrNotFound
		}
		return models.Category{}, err
	}

	return cat, nil
}

// Create adds a new category to the database
func (r *categoryRepository) Create(ctx context.Context, cat models.Category) (models.Category, error) {
	// Check if name already exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1)`
	if err := r.db.QueryRow(ctx, checkQuery, cat.Name).Scan(&exists); err != nil {
		return models.Category{}, err
	}
	if exists {
		return models.Category{}, ErrNameExists
	}

	// Insert the new category
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, cat.Name, cat.Description).Scan(&cat.ID)
	if err != nil {
		return models.Category{}, err
	}

	return cat, nil
}

// Update updates an existing category
func (r *categoryRepository) Update(ctx context.Context, id int, cat models.Category) (models.Category, error) {
	query := `UPDATE categories SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description`

	var updated models.Category
	err := r.db.QueryRow(ctx, query, cat.Name, cat.Description, id).Scan(&updated.ID, &updated.Name, &updated.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Category{}, ErrNotFound
		}
		return models.Category{}, err
	}

	return updated, nil
}

// Delete removes a category by its ID
func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
