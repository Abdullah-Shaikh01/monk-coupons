package models

import "time"

type Coupon struct {
	ID                  int       `json:"id"` // optional in POST (ignored if 0)
	Type                string    `json:"type"` // required
	DiscountValue       *float64   `json:"discount_value"`
	DiscountType        string    `json:"discount_type"`
	RepetitionThreshold *int      `json:"repetition_threshold"`
	BuyQuantity            *int      `json:"buy_quantity"`
	GetQuantity            *int      `json:"get_quantity"`
	ExpirationDate      time.Time `json:"expiration_date"`
	ProductID           *int      `json:"product_id"`
}

