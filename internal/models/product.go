package models

import "time"

type Product struct {
	ID            string    `json:"id"`
	CategoryID    string    `json:"category_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Price         int       `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
