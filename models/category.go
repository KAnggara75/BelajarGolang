package models

// Category represents a category entity
type Category struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
