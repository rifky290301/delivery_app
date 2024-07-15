package models

import "time"

type Rating struct {
	Id        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	BuyerID   int       `json:"buyer_id"`
	ShopID    int       `json:"shop_id"`
	Rating    int       `json:"rating"`
	Feedback  string    `json:"feedback"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt NullTime  `json:"updated_at"`
}
