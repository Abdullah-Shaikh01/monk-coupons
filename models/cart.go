package models

type CartItem struct {
    ProductID     int     `json:"product_id"`
    Quantity      int     `json:"quantity"`
    Price         float64 `json:"price"`
	TotalDiscount float64 `json:"total_discount"`
}

type Cart struct {
	Items []CartItem `json:"items"`
}