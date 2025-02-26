package service

import (
	"context"
	"gateService/internal/domain/entity"
	"gateService/internal/domain/repository"
	"gateService/internal/interfaces/dto"
)

type ProductServiceImpl struct {
	productRepository repository.ProductRepository
}

func NewProductServiceImpl(productRepository repository.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{productRepository: productRepository}
}

func (p *ProductServiceImpl) CreateProduct(ctx context.Context, request *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
	product := entity.Product{
		ProductName: request.ProductName,
		Description: request.Description,
		Image:       request.Image,
		Category:    request.Category,
	}

	productVariant := entity.ProductVariant{
		Color: request.Color,
		Price: request.Price,
		Stock: request.Stock,
		Size:  request.Size,
	}

	productID, err := p.productRepository.CreateProduct(ctx, &product, &productVariant)
	if err != nil {
		return nil, err
	}

	return &dto.CreateProductResponse{Code: 200, ProductID: productID}, nil
}

func (p *ProductServiceImpl) GetProducts(ctx context.Context, request *dto.GetProductsRequest) (*dto.GetProductsResponse, error) {
	products, err := p.productRepository.GetProducts(ctx, request.Page)
	if err != nil {
		return nil, err
	}

	return &dto.GetProductsResponse{Code: 200, Products: products}, nil
}
