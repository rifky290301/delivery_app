package models

import "time"

type Order struct {
	Id        int       `json:"id"`
	BuyerId   int       `json:"buyer_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` //'pending','paid','shipped','delivered','cancelled'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt NullTime  `json:"updated_at"`
	Buyer     User      `json:"buyer"`
}
