package models

import (
	"time"
)

type OrderItem struct {
	OrderItemID int       `json:"order_item_id"`
	OrderID     int       `json:"order_id"`
	ProductID   int       `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   NullTime  `json:"updated_at"`
	Order       Order     `json:"order"`
	Product     Product   `json:"product"`
}
