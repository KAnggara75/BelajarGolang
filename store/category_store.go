package store

import (
	"errors"

	"github.com/KAnggara75/BelajarGolang/models"
)

var (
	ErrNotFound = errors.New("category not found")
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
func (s *CategoryStore) Create(cat models.Category) models.Category {
	cat.ID = s.nextID
	s.nextID++
	s.categories[cat.ID] = cat
	return cat
}
