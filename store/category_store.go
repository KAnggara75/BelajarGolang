package store

import (
	"errors"

	"github.com/KAnggara75/BelajarGolang/models"
)

var (
	ErrNotFound   = errors.New("category not found")
	ErrNameExists = errors.New("category name already exists")
)

// CategoryStore is an in-memory store for categories
type CategoryStore struct {
	categories map[int]models.Category
	nextID     int
}

// NewCategoryStore creates a new CategoryStore
func NewCategoryStore() *CategoryStore {
	return &CategoryStore{
		categories: make(map[int]models.Category),
		nextID:     1,
	}
}

// GetAll returns all categories
func (s *CategoryStore) GetAll() []models.Category {
	result := make([]models.Category, 0, len(s.categories))
	for _, cat := range s.categories {
		result = append(result, cat)
	}
	return result
}

// GetByID returns a category by ID
func (s *CategoryStore) GetByID(id int) (models.Category, error) {
	cat, exists := s.categories[id]
	if !exists {
		return models.Category{}, ErrNotFound
	}
	return cat, nil
}

// Create adds a new category
func (s *CategoryStore) Create(cat models.Category) (models.Category, error) {
	// Check if name already exists
	for _, existing := range s.categories {
		if existing.Name == cat.Name {
			return models.Category{}, ErrNameExists
		}
	}

	cat.ID = s.nextID
	s.nextID++
	s.categories[cat.ID] = cat
	return cat, nil
}

// Update updates an existing category
func (s *CategoryStore) Update(id int, cat models.Category) (models.Category, error) {
	if _, exists := s.categories[id]; !exists {
		return models.Category{}, ErrNotFound
	}

	cat.ID = id
	s.categories[id] = cat
	return cat, nil
}

// Delete removes a category by ID
func (s *CategoryStore) Delete(id int) error {
	if _, exists := s.categories[id]; !exists {
		return ErrNotFound
	}

	delete(s.categories, id)
	return nil
}

// SeedData initializes the store with sample data
func (s *CategoryStore) SeedData() {
	initialData := []models.Category{
		{Name: "Electronics", Description: "Electronic devices and gadgets"},
		{Name: "Clothing", Description: "Apparel and fashion items"},
		{Name: "Books", Description: "Books and reading materials"},
		{Name: "Food & Beverages", Description: "Food products and drinks"},
		{Name: "Sports", Description: "Sports equipment and accessories"},
	}

	for _, cat := range initialData {
		_, _ = s.Create(cat)
	}
}
