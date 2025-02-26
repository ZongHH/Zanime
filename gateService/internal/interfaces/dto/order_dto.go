package dto

import "gateService/internal/domain/entity"

type CreateOrderRequest struct {
	UserID      int
	ProductID   int     `json:"productID"`
	ProductName string  `json:"productName"`
	Size        string  `json:"selectedSize"`
	Color       string  `json:"selectedColor"`
	Price       float64 `json:"price"`
	Phone       string  `json:"phone"`
	UserName    string  `json:"userName"`
	Address     string  `json:"address"`
}

type CreateOrderResponse struct {
	Code        int     `json:"code"`
	OrderID     string  `json:"orderID"`
	UserID      int     `json:"userID"`
	ProductID   int     `json:"productID"`
	ProductName string  `json:"productName"`
	Size        string  `json:"selectedSize"`
	Color       string  `json:"selectedColor"`
	Price       float64 `json:"price"`
	Phone       string  `json:"phone"`
	UserName    string  `json:"userName"`
	Address     string  `json:"address"`
	CreatedTime string  `json:"createdTime"`
}

type GetOrderRequest struct {
	UserID   int `form:"userID"`
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type GetOrderResponse struct {
	Code   int             `json:"code"`
	Orders []*entity.Order `json:"orders"`
}

type PayRequest struct {
	OrderID string `json:"orderID"`
}

type PayResponse struct {
	Code int `json:"code"`
}
