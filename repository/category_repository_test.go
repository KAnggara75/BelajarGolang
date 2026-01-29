package repository

import (
	"context"
	"testing"

	"github.com/KAnggara75/BelajarGolang/models"
)

// mockRepository is a simple in-memory implementation for testing
type mockRepository struct {
	categories map[int]models.Category
	nextID     int
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		categories: make(map[int]models.Category),
		nextID:     1,
	}
}

func (m *mockRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	result := make([]models.Category, 0, len(m.categories))
	for _, cat := range m.categories {
		result = append(result, cat)
	}
	return result, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int) (models.Category, error) {
	cat, exists := m.categories[id]
	if !exists {
		return models.Category{}, ErrNotFound
	}
	return cat, nil
}

func (m *mockRepository) Create(ctx context.Context, cat models.Category) (models.Category, error) {
	for _, existing := range m.categories {
		if existing.Name == cat.Name {
			return models.Category{}, ErrNameExists
		}
	}

	cat.ID = m.nextID
	m.nextID++
	m.categories[cat.ID] = cat
	return cat, nil
}

func (m *mockRepository) Update(ctx context.Context, id int, cat models.Category) (models.Category, error) {
	if _, exists := m.categories[id]; !exists {
		return models.Category{}, ErrNotFound
	}

	cat.ID = id
	m.categories[id] = cat
	return cat, nil
}

func (m *mockRepository) Delete(ctx context.Context, id int) error {
	if _, exists := m.categories[id]; !exists {
		return ErrNotFound
	}

	delete(m.categories, id)
	return nil
}

// TestMockRepository_GetAll tests GetAll functionality
func TestMockRepository_GetAll(t *testing.T) {
	repo := newMockRepository()
	ctx := context.Background()

	// Test empty
	result, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 categories, got %d", len(result))
	}

	// Add some categories
	_, _ = repo.Create(ctx, models.Category{Name: "Cat 1"})
	_, _ = repo.Create(ctx, models.Category{Name: "Cat 2"})

	result, err = repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(result))
	}
}

// TestMockRepository_Create tests Create functionality
func TestMockRepository_Create(t *testing.T) {
	repo := newMockRepository()
	ctx := context.Background()

	cat := models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
	}

	created, err := repo.Create(ctx, cat)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("Expected ID 1, got %d", created.ID)
	}
	if created.Name != "Electronics" {
		t.Errorf("Expected name 'Electronics', got '%s'", created.Name)
	}
}

// TestMockRepository_Create_DuplicateName tests duplicate name prevention
func TestMockRepository_Create_DuplicateName(t *testing.T) {
	repo := newMockRepository()
	ctx := context.Background()

	_, _ = repo.Create(ctx, models.Category{Name: "Electronics"})
	_, err := repo.Create(ctx, models.Category{Name: "Electronics"})

	if err == nil {
		t.Fatal("Expected error for duplicate name")
	}
	if err != ErrNameExists {
		t.Errorf("Expected ErrNameExists, got %v", err)
	}
}

// TestMockRepository_GetByID tests GetByID functionality
func TestMockRepository_GetByID(t *testing.T) {
	repo := newMockRepository()
	ctx := context.Background()

	created, _ := repo.Create(ctx, models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
	})

	retrieved, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, retrieved.ID)
	}
	if retrieved.Name != created.Name {
		t.Errorf("Expected name '%s', got '%s'", created.Name, retrieved.Name)
	}
}

// TestMockRepository_GetByID_NotFound tests GetByID with non-existent ID
func TestMockRepository_GetByID_NotFound(t *testing.T) {
	repo := newMockRepository()
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 999)
	if err == nil {
		t.Fatal("Expected error for non-existent ID")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestMockRepository_Update tests Update functionality
func TestMockRepository_Update(t *testing.T) {
	repo := newMockRepository()
	ctx := context.Background()

	original, _ := repo.Create(ctx, models.Category{
		Name:        "Original",
		Description: "Original description",
	})

	updated, err := repo.Update(ctx, original.ID, models.Category{
		Name:        "Updated",
		Description: "Updated description",
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if updated.ID != original.ID {
		t.Errorf("Expected ID %d, got %d", original.ID, updated.ID)
	}
	if updated.Name != "Updated" {
		t.Errorf("Expected name 'Updated', got '%s'", updated.Name)
	}
}

// TestMockRepository_Update_NotFound tests Update with non-existent ID
func TestMockRepository_Update_NotFound(t *testing.T) {
	repo := newMockRepository()
	ctx := context.Background()

	_, err := repo.Update(ctx, 999, models.Category{Name: "Test"})
	if err == nil {
		t.Fatal("Expected error for non-existent ID")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestMockRepository_Delete tests Delete functionality
func TestMockRepository_Delete(t *testing.T) {
	repo := newMockRepository()
	ctx := context.Background()

	created, _ := repo.Create(ctx, models.Category{Name: "To Delete"})

	err := repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	_, err = repo.GetByID(ctx, created.ID)
	if err != ErrNotFound {
		t.Error("Category should have been deleted")
	}
}

// TestMockRepository_Delete_NotFound tests Delete with non-existent ID
func TestMockRepository_Delete_NotFound(t *testing.T) {
	repo := newMockRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	if err == nil {
		t.Fatal("Expected error for non-existent ID")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestCategoryRepositoryInterface ensures mockRepository implements CategoryRepository
func TestCategoryRepositoryInterface(t *testing.T) {
	var _ CategoryRepository = (*mockRepository)(nil)
}
