package repository

import (
	"context"
	"errors"

	"github.com/KAnggara75/BelajarGolang/models"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var (
	ErrNotFound   = errors.New("category not found")
	ErrNameExists = errors.New("category name already exists")
)

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (models.Category, error)
	GetByName(ctx context.Context, name string) ([]models.Category, error)
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
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "CategoryRepository.GetAll")
	defer span.End()

	query := `SELECT id, name, description FROM categories ORDER BY id`
	span.SetAttributes(attribute.String("db.query", query))

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Return empty slice instead of nil
	if categories == nil {
		categories = []models.Category{}
	}

	span.SetAttributes(attribute.Int("result.count", len(categories)))
	return categories, nil
}

// GetByID returns a category by its ID
func (r *categoryRepository) GetByID(ctx context.Context, id int) (models.Category, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "CategoryRepository.GetByID")
	defer span.End()

	span.SetAttributes(attribute.Int("category.id", id))

	query := `SELECT id, name, description FROM categories WHERE id = $1`
	span.SetAttributes(attribute.String("db.query", query))

	var cat models.Category
	err := r.db.QueryRow(ctx, query, id).Scan(&cat.ID, &cat.Name, &cat.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetStatus(codes.Error, "category not found")
			return models.Category{}, ErrNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return models.Category{}, err
	}

	return cat, nil
}

// GetByName returns all categories matching the name (case-insensitive partial match)
func (r *categoryRepository) GetByName(ctx context.Context, name string) ([]models.Category, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "CategoryRepository.GetByName")
	defer span.End()

	span.SetAttributes(attribute.String("category.name", name))

	query := `SELECT id, name, description FROM categories WHERE name ILIKE '%' || $1 || '%' ORDER BY id`
	span.SetAttributes(attribute.String("db.query", query))

	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Return empty slice instead of nil
	if categories == nil {
		categories = []models.Category{}
	}

	span.SetAttributes(attribute.Int("result.count", len(categories)))
	return categories, nil
}

// Create adds a new category to the database
func (r *categoryRepository) Create(ctx context.Context, cat models.Category) (models.Category, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "CategoryRepository.Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("category.name", cat.Name),
		attribute.String("category.description", cat.Description),
	)

	// Check if name already exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1)`
	if err := r.db.QueryRow(ctx, checkQuery, cat.Name).Scan(&exists); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return models.Category{}, err
	}
	if exists {
		span.SetStatus(codes.Error, "category name already exists")
		return models.Category{}, ErrNameExists
	}

	// Insert the new category
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	span.SetAttributes(attribute.String("db.query", query))
	err := r.db.QueryRow(ctx, query, cat.Name, cat.Description).Scan(&cat.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return models.Category{}, err
	}

	span.SetAttributes(attribute.Int("category.id", cat.ID))
	return cat, nil
}

// Update updates an existing category
func (r *categoryRepository) Update(ctx context.Context, id int, cat models.Category) (models.Category, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "CategoryRepository.Update")
	defer span.End()

	span.SetAttributes(
		attribute.Int("category.id", id),
		attribute.String("category.name", cat.Name),
		attribute.String("category.description", cat.Description),
	)

	query := `UPDATE categories SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description`
	span.SetAttributes(attribute.String("db.query", query))

	var updated models.Category
	err := r.db.QueryRow(ctx, query, cat.Name, cat.Description, id).Scan(&updated.ID, &updated.Name, &updated.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetStatus(codes.Error, "category not found")
			return models.Category{}, ErrNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return models.Category{}, err
	}

	return updated, nil
}

// Delete removes a category by its ID
func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "CategoryRepository.Delete")
	defer span.End()

	span.SetAttributes(attribute.Int("category.id", id))

	query := `DELETE FROM categories WHERE id = $1`
	span.SetAttributes(attribute.String("db.query", query))

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if result.RowsAffected() == 0 {
		span.SetStatus(codes.Error, "category not found")
		return ErrNotFound
	}

	span.SetAttributes(attribute.Int64("rows.affected", result.RowsAffected()))
	return nil
}
