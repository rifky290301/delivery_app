package models

import (
	"time"
)

type Shop struct {
	Id              int       `json:"id"`
	SellerID        int       `json:"user_id"`
	ShopName        string    `json:"shop_name"`
	ShopDescription string    `json:"shop_description"`
	ShopAddress     string    `json:"shop_address"`
	CreatedAt       time.Time `json:"created_at"`
	Seller          User      `json:"seller"`
}
