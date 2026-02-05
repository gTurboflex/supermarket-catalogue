package models

import "time"

type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	Stock         int       `json:"stock"`
	Image         string    `json:"image,omitempty"`
	CategoryID    int       `json:"category_id"`
	OwnerID       int       `json:"owner_id,omitempty"`
	Barcode       string    `json:"barcode,omitempty"`
	Unit          string    `json:"unit,omitempty"`
	UnitPrice     float64   `json:"unit_price,omitempty"`
	LastUpdated   time.Time `json:"last_updated,omitempty"`
	SupermarketID int       `json:"supermarket_id,omitempty"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	OwnerID     int    `json:"owner_id,omitempty"`
}

type Supermarket struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address,omitempty"`
	OwnerID   int       `json:"owner_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
