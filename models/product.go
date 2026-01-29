package models

// Product represents a product entity for API responses
type Product struct {
	ID         int       `json:"-"`
	Name       string    `json:"name"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID int       `json:"-"`
	Category   *Category `json:"category,omitempty"`
}

// ProductInput is used for API input to accept category_id
type ProductInput struct {
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	CategoryID int     `json:"category_id,omitempty"`
}

// ToProduct converts a ProductInput to a Product
func (r *ProductInput) ToProduct() Product {
	return Product{
		Name:       r.Name,
		Price:      r.Price,
		Stock:      r.Stock,
		CategoryID: r.CategoryID,
	}
}
