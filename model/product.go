package model

type Product struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	ThumbnailImage string      `json:"thumbnail_image"`
	Price          float64     `json:"price"`
	Discount       float64     `json:"discount_percentage"`
	DiscountPrice  float64     `json:"discount_price"`
	AverageRating  float64     `json:"average_rating"`
	Sold           int         `json:"sold"`
	IsNEW          bool        `json:"is_new"`
}
type ProductID struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	Images         []string    `json:"images"`
	Category       string      `json:"category_name"`
	Price          float64     `json:"price"`
	Variant        interface{} `json:"variant"`
	Discount       float64     `json:"discount_percentage"`
	DiscountPrice  float64     `json:"discount_price"`
	AverageRating  float64     `json:"average_rating"`
	Sold           int         `json:"sold"`
	IsNEW          bool        `json:"is_new"`
}
