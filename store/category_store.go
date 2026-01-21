package store

import (
	"github.com/KAnggara75/BelajarGolang/models"
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
