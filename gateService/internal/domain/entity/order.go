package entity

type Order struct {
	OrderID     string  `json:"order_id"`
	UserID      int     `json:"user_id"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`

	// 产品选择相关
	SelectedSize  string `json:"selected_size"`
	SelectedColor string `json:"selected_color"`

	// 用户信息
	UserName string `json:"user_name"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`

	// 订单状态相关
	Status     string `json:"status"`
	QRCodeURL  string `json:"qr_code_url"`
	CreateTime string `json:"create_time"`

	// 产品信息
	Image string `json:"image"`
}
