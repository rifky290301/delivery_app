package models

import (
	"time"
)

type Product struct {
	Id          int        `json:"id"`
	ShopID      int        `json:"shop_id"`
	Name        string     `json:"product_name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   NullString `json:"updated_at"`
	// UpdatedAt   time.Time `json:"updated_at"`
}
