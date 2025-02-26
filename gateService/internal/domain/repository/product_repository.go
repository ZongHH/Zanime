package repository

import (
	"context"
	"gateService/internal/domain/entity"
)

type ProductRepository interface {
	// CreateProduct 创建商品及其变体
	// 参数:
	// - ctx: 上下文
	// - product: 商品信息
	// - productVariant: 商品变体信息
	// 返回:
	// - int: 创建的商品ID
	// - error: 错误信息
	CreateProduct(ctx context.Context, product *entity.Product, productVariant *entity.ProductVariant) (int, error)

	// DeleteProduct 删除商品
	// 参数:
	// - ctx: 上下文
	// - productID: 商品ID
	// 返回:
	// - error: 错误信息
	DeleteProduct(ctx context.Context, productID int) error

	// UpdateProduct 更新商品及其变体信息
	// 参数:
	// - ctx: 上下文
	// - product: 更新的商品信息
	// - productVariant: 更新的商品变体信息
	// 返回:
	// - error: 错误信息
	UpdateProduct(ctx context.Context, product *entity.Product, productVariant *entity.ProductVariant) error

	// GetProducts 分页获取商品列表
	// 参数:
	// - ctx: 上下文
	// - page: 页码
	// 返回:
	// - []*entity.Product: 商品列表
	// - error: 错误信息
	GetProducts(ctx context.Context, page int) ([]*entity.Product, error)

	// GetProductStock 获取商品库存
	// 参数:
	// - ctx: 上下文
	// - productID: 商品ID
	// - size: 商品尺码
	// - color: 商品颜色
	// 返回:
	// - int: 库存数量
	// - error: 错误信息
	GetProductStock(ctx context.Context, productID int, size string, color string) (int, error)

	// UpdateProductStock 更新商品库存
	// 参数:
	// - ctx: 上下文
	// - productID: 商品ID
	// - size: 商品尺码
	// - color: 商品颜色
	// - quantity: 增加的库存数量
	// 返回:
	// - error: 错误信息
	IncreaseProductStock(ctx context.Context, productID int, size string, color string, quantity int) error

	// DecreaseProductStock 减少商品库存
	// 参数:
	// - ctx: 上下文
	// - productID: 商品ID
	// - size: 商品尺码
	// - color: 商品颜色
	// - quantity: 减少的库存数量
	// 返回:
	// - error: 错误信息
	DecreaseProductStock(ctx context.Context, productID int, size string, color string, quantity int) error
}
