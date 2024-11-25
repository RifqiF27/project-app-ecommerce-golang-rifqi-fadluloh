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
	ShippingAddress *string     `json:"shipping_address"`
	AddressIndex    int         `json:"-"`
	Shipping        string      `json:"shipping"`
	TotalAmount     float64     `json:"total_amount"`
}

type OrderItem struct {
	UserID        int     `json:"user_id,omitempty"`
	ProductID     []int   `json:"product_id,omitempty"`
	AddressIndex  int     `json:"address_index,omitempty"`
	ProductName   string  `json:"product_name"`
	Image         string  `json:"image"`
	Quantity      int     `json:"quantity,omitempty"`
	SubtotalPrice float64 `json:"subtotal_price"`
}
