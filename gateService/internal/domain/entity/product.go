package entity

type Product struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Category    string `json:"category"`

	// 额外字段
	Prices []float64 `json:"prices"`
	Stocks []int     `json:"stocks"`
	Sizes  []string  `json:"sizes"`
	Colors []string  `json:"colors"`
}

type ProductVariant struct {
	VariantID int     `json:"variant_id"`
	ProductID int     `json:"product_id"`
	Color     string  `json:"color"`
	Price     float64 `json:"price"`
	Stock     int     `json:"stock"`
	Size      string  `json:"size"`
}
