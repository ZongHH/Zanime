package dto

import "gateService/internal/domain/entity"

type CreateProductRequest struct {
	ProductName string `json:"product_name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Category    string `json:"category"`

	Color string  `json:"color"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
	Size  string  `json:"size"`
}

type CreateProductResponse struct {
	Code      int `json:"code"`
	ProductID int `json:"product_id"`
}

type GetProductsRequest struct {
	Page int `form:"page" binding:"required"`
}

type GetProductsResponse struct {
	Code     int               `json:"code"`
	Products []*entity.Product `json:"products"`
}
