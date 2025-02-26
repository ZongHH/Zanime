// package service 提供了与商品相关的业务逻辑服务
package service

import (
	"context"
	"gateService/internal/interfaces/dto"
)

// ProductService 定义了商品服务的接口
type ProductService interface {
	// CreateProduct 创建新的商品
	// 参数:
	// - ctx: 上下文信息
	// - product: 创建商品的请求数据,包含商品基本信息和变体信息
	// 返回:
	// - *dto.CreateProductResponse: 创建商品的响应数据,包含创建的商品ID等信息
	// - error: 创建过程中的错误信息
	CreateProduct(ctx context.Context, product *dto.CreateProductRequest) (*dto.CreateProductResponse, error)

	// GetProducts 获取商品列表
	// 参数:
	// - ctx: 上下文信息
	// - request: 获取商品列表的请求参数,包含分页等信息
	// 返回:
	// - *dto.GetProductsResponse: 商品列表响应数据
	// - error: 获取过程中的错误信息
	GetProducts(ctx context.Context, request *dto.GetProductsRequest) (*dto.GetProductsResponse, error)
}
