package store

import (
	"testing"

	"github.com/KAnggara75/BelajarGolang/models"
)

// TestNewCategoryStore tests the store initialization
func TestNewCategoryStore(t *testing.T) {
	s := NewCategoryStore()

	if s == nil {
		t.Fatal("Expected non-nil store")
	}

	if s.categories == nil {
		t.Error("Expected categories map to be initialized")
	}

	if s.nextID != 1 {
		t.Errorf("Expected nextID to be 1, got %d", s.nextID)
	}
}

// TestGetAll_Empty tests GetAll on empty store
func TestGetAll_Empty(t *testing.T) {
	s := NewCategoryStore()

	result := s.GetAll()

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d items", len(result))
	}
}

// TestGetAll_WithData tests GetAll with data
func TestGetAll_WithData(t *testing.T) {
	s := NewCategoryStore()

	_, _ = s.Create(models.Category{Name: "Category 1"})
	_, _ = s.Create(models.Category{Name: "Category 2"})

	result := s.GetAll()

	if len(result) != 2 {
		t.Errorf("Expected 2 items, got %d", len(result))
	}
}

// TestCreate_Success tests successful category creation
func TestCreate_Success(t *testing.T) {
	s := NewCategoryStore()

	cat := models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
	}

	created, err := s.Create(cat)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("Expected ID 1, got %d", created.ID)
	}

	if created.Name != "Electronics" {
		t.Errorf("Expected name 'Electronics', got '%s'", created.Name)
	}

	if created.Description != "Electronic devices" {
		t.Errorf("Expected description 'Electronic devices', got '%s'", created.Description)
	}

	// Verify nextID incremented
	if s.nextID != 2 {
		t.Errorf("Expected nextID to be 2, got %d", s.nextID)
	}
}

// TestCreate_IncrementingIDs tests that IDs are assigned incrementally
func TestCreate_IncrementingIDs(t *testing.T) {
	s := NewCategoryStore()

	cat1, _ := s.Create(models.Category{Name: "Cat 1"})
	cat2, _ := s.Create(models.Category{Name: "Cat 2"})
	cat3, _ := s.Create(models.Category{Name: "Cat 3"})

	if cat1.ID != 1 {
		t.Errorf("Expected first ID to be 1, got %d", cat1.ID)
	}
	if cat2.ID != 2 {
		t.Errorf("Expected second ID to be 2, got %d", cat2.ID)
	}
	if cat3.ID != 3 {
		t.Errorf("Expected third ID to be 3, got %d", cat3.ID)
	}
}

// TestCreate_DuplicateName tests that duplicate names are rejected
func TestCreate_DuplicateName(t *testing.T) {
	s := NewCategoryStore()

	_, _ = s.Create(models.Category{Name: "Electronics"})
	_, err := s.Create(models.Category{Name: "Electronics"})

	if err == nil {
		t.Fatal("Expected error for duplicate name")
	}

	if err != ErrNameExists {
		t.Errorf("Expected ErrNameExists, got %v", err)
	}
}

// TestGetByID_Success tests GetByID with valid ID
func TestGetByID_Success(t *testing.T) {
	s := NewCategoryStore()

	original, _ := s.Create(models.Category{
		Name:        "Electronics",
		Description: "Electronic devices",
	})

	retrieved, err := s.GetByID(original.ID)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if retrieved.ID != original.ID {
		t.Errorf("Expected ID %d, got %d", original.ID, retrieved.ID)
	}

	if retrieved.Name != original.Name {
		t.Errorf("Expected name '%s', got '%s'", original.Name, retrieved.Name)
	}

	if retrieved.Description != original.Description {
		t.Errorf("Expected description '%s', got '%s'", original.Description, retrieved.Description)
	}
}

// TestGetByID_NotFound tests GetByID with non-existent ID
func TestGetByID_NotFound(t *testing.T) {
	s := NewCategoryStore()

	_, err := s.GetByID(999)

	if err == nil {
		t.Fatal("Expected error for non-existent ID")
	}

	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestUpdate_Success tests successful category update
func TestUpdate_Success(t *testing.T) {
	s := NewCategoryStore()

	original, _ := s.Create(models.Category{
		Name:        "Original",
		Description: "Original description",
	})

	updated, err := s.Update(original.ID, models.Category{
		Name:        "Updated",
		Description: "Updated description",
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if updated.ID != original.ID {
		t.Errorf("Expected ID %d to be preserved, got %d", original.ID, updated.ID)
	}

	if updated.Name != "Updated" {
		t.Errorf("Expected name 'Updated', got '%s'", updated.Name)
	}

	if updated.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", updated.Description)
	}

	// Verify the update was persisted
	retrieved, _ := s.GetByID(original.ID)
	if retrieved.Name != "Updated" {
		t.Error("Update was not persisted")
	}
}

// TestUpdate_NotFound tests Update with non-existent ID
func TestUpdate_NotFound(t *testing.T) {
	s := NewCategoryStore()

	_, err := s.Update(999, models.Category{Name: "Test"})

	if err == nil {
		t.Fatal("Expected error for non-existent ID")
	}

	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestUpdate_PreservesID tests that Update sets the correct ID
func TestUpdate_PreservesID(t *testing.T) {
	s := NewCategoryStore()

	original, _ := s.Create(models.Category{Name: "Original"})

	// Try to update with a different ID in the category struct
	updated, _ := s.Update(original.ID, models.Category{
		ID:   999, // This should be ignored
		Name: "Updated",
	})

	if updated.ID != original.ID {
		t.Errorf("Expected ID %d to be preserved, got %d", original.ID, updated.ID)
	}
}

// TestDelete_Success tests successful category deletion
func TestDelete_Success(t *testing.T) {
	s := NewCategoryStore()

	created, _ := s.Create(models.Category{Name: "To Delete"})

	err := s.Delete(created.ID)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify deletion
	_, err = s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Error("Category should have been deleted")
	}
}

// TestDelete_NotFound tests Delete with non-existent ID
func TestDelete_NotFound(t *testing.T) {
	s := NewCategoryStore()

	err := s.Delete(999)

	if err == nil {
		t.Fatal("Expected error for non-existent ID")
	}

	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestDelete_DoesNotAffectOthers tests that deletion doesn't affect other categories
func TestDelete_DoesNotAffectOthers(t *testing.T) {
	s := NewCategoryStore()

	cat1, _ := s.Create(models.Category{Name: "Category 1"})
	cat2, _ := s.Create(models.Category{Name: "Category 2"})
	cat3, _ := s.Create(models.Category{Name: "Category 3"})

	err := s.Delete(cat2.ID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify cat1 and cat3 still exist
	if _, err := s.GetByID(cat1.ID); err != nil {
		t.Error("Category 1 should still exist")
	}

	if _, err := s.GetByID(cat3.ID); err != nil {
		t.Error("Category 3 should still exist")
	}

	// Verify correct count
	all := s.GetAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 categories remaining, got %d", len(all))
	}
}

// TestSeedData tests the SeedData function
func TestSeedData(t *testing.T) {
	s := NewCategoryStore()

	s.SeedData()

	categories := s.GetAll()

	if len(categories) != 5 {
		t.Errorf("Expected 5 seeded categories, got %d", len(categories))
	}

	// Verify that nextID is correctly updated after seeding
	if s.nextID != 6 {
		t.Errorf("Expected nextID to be 6 after seeding 5 items, got %d", s.nextID)
	}

	// Check that we can still create new categories
	newCat, err := s.Create(models.Category{Name: "New Category"})
	if err != nil {
		t.Fatalf("Failed to create category after seeding: %v", err)
	}

	if newCat.ID != 6 {
		t.Errorf("Expected new category ID to be 6, got %d", newCat.ID)
	}
}

// TestSeedData_ExpectedCategories tests that SeedData creates expected categories
func TestSeedData_ExpectedCategories(t *testing.T) {
	s := NewCategoryStore()
	s.SeedData()

	expectedNames := []string{
		"Electronics",
		"Clothing",
		"Books",
		"Food & Beverages",
		"Sports",
	}

	categories := s.GetAll()

	// Create a map of existing names for easier lookup
	existingNames := make(map[string]bool)
	for _, cat := range categories {
		existingNames[cat.Name] = true
	}

	for _, name := range expectedNames {
		if !existingNames[name] {
			t.Errorf("Expected category '%s' not found in seeded data", name)
		}
	}
}

// TestConcurrentCreatePrevention tests that duplicate names can't be created
func TestConcurrentCreatePrevention(t *testing.T) {
	s := NewCategoryStore()

	// Create first category
	_, err1 := s.Create(models.Category{Name: "Unique Name"})
	// Try to create duplicate
	_, err2 := s.Create(models.Category{Name: "Unique Name"})

	if err1 != nil {
		t.Fatalf("First create should succeed: %v", err1)
	}

	if err2 == nil {
		t.Error("Second create with same name should fail")
	}

	if err2 != ErrNameExists {
		t.Errorf("Expected ErrNameExists, got %v", err2)
	}
}
