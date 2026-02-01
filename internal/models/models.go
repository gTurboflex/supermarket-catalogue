package models

type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	Stock      int     `json:"stock"`
	Image      string  `json:"image,omitempty"`
	CategoryID int     `json:"category_id"`
	AdminID    int     `json:"admin_id,omitempty"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	AdminID     int    `json:"admin_id,omitempty"`
}