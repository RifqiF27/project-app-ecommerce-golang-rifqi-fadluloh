package model

type Checkout struct {
	ID         int     `json:"id,omitempty"`
	UserID     int     `json:"user_id,omitempty"`
	ProductID  int     `json:"product_id,omitempty"`
	Name       string  `json:"name,omitempty"`
	Image      string  `json:"image,omitempty"`
	Price      float64 `json:"price,omitempty"`
	Quantity   int     `json:"quantity,omitempty"`
	TotalPrice float64 `json:"total_price,omitempty"`
	TotalCarts int     `json:"total_carts,omitempty"`
}

type OrderResponse struct {
	OrderID         int         `json:"order_id"`
	Items           []OrderItem `json:"items"`
	ShippingAddress string      `json:"shipping_address"`
	Shipping        string      `json:"shipping"`
	TotalAmount     float64     `json:"total_amount"`
}

type OrderItem struct {
	ProductID     []int     `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Image         string  `json:"image"`
	SubtotalPrice float64 `json:"subtotal_price"`
}
