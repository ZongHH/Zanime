package database

import (
	"context"
	"database/sql"
	"gateService/internal/domain/entity"
)

type ProductRepositoryImpl struct {
	db *sql.DB
}

func NewProductRepositoryImpl(db *sql.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db}
}

func (p *ProductRepositoryImpl) CreateProduct(ctx context.Context, product *entity.Product, productVariant *entity.ProductVariant) (int, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var productID int
	err = tx.QueryRowContext(ctx, "INSERT INTO products (product_name, description, image, category) VALUES (?, ?, ?, ?)",
		product.ProductName, product.Description, product.Image, product.Category).Scan(&productID)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO product_variants (product_id, color, price, stock, size) VALUES (?, ?, ?, ?, ?)",
		productID, productVariant.Color, productVariant.Price, productVariant.Stock, productVariant.Size)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return productID, nil
}

func (p *ProductRepositoryImpl) DeleteProduct(ctx context.Context, productID int) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM product_variants WHERE product_id = ?", productID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM products WHERE product_id = ?", productID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *ProductRepositoryImpl) UpdateProduct(ctx context.Context, product *entity.Product, productVariant *entity.ProductVariant) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE products SET product_name = ?, description = ?, image = ?, category = ? WHERE product_id = ?",
		product.ProductName, product.Description, product.Image, product.Category, product.ProductID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE product_variants SET color = ?, price = ?, stock = ?, size = ? WHERE product_id = ?",
		productVariant.Color, productVariant.Price, productVariant.Stock, productVariant.Size, product.ProductID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *ProductRepositoryImpl) GetProducts(ctx context.Context, page int) ([]*entity.Product, error) {
	limit := 15
	offset := (page - 1) * limit

	query := `
		SELECT p.product_id, product_name, description, image, category, color, size, price, stock 
		FROM products p 
		JOIN product_variants pv ON p.product_id = pv.product_id 
		WHERE p.product_id >= ? AND p.product_id < ?`

	rows, err := p.db.QueryContext(ctx, query, offset+1, offset+limit+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productMap := make(map[int]*entity.Product)
	var products []*entity.Product

	for rows.Next() {
		var productID int
		var productName, description, image, category, color, size string
		var price float64
		var stock int

		if err := rows.Scan(&productID, &productName, &description, &image, &category, &color, &size, &price, &stock); err != nil {
			return nil, err
		}

		if product, exists := productMap[productID]; exists {
			product.Colors = append(product.Colors, color)
			product.Sizes = append(product.Sizes, size)
			product.Prices = append(product.Prices, price)
			product.Stocks = append(product.Stocks, stock)
		} else {
			product := &entity.Product{
				ProductID:   productID,
				ProductName: productName,
				Description: description,
				Image:       image,
				Category:    category,
				Colors:      []string{color},
				Sizes:       []string{size},
				Prices:      []float64{price},
				Stocks:      []int{stock},
			}
			productMap[productID] = product
			products = append(products, product)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductRepositoryImpl) GetProductStock(ctx context.Context, productID int, size string, color string) (int, error) {
	query := "SELECT stock FROM product_variants WHERE product_id = ? AND size = ? AND color = ?"
	var stock int
	err := p.db.QueryRowContext(ctx, query, productID, size, color).Scan(&stock)
	if err != nil {
		return 0, err
	}
	return stock, nil
}

func (p *ProductRepositoryImpl) IncreaseProductStock(ctx context.Context, productID int, size string, color string, quantity int) error {
	query := "UPDATE product_variants SET stock = stock + ? WHERE product_id = ? AND size = ? AND color = ?"
	_, err := p.db.ExecContext(ctx, query, quantity, productID, size, color)
	return err
}

func (p *ProductRepositoryImpl) DecreaseProductStock(ctx context.Context, productID int, size string, color string, quantity int) error {
	query := "UPDATE product_variants SET stock = stock - ? WHERE product_id = ? AND size = ? AND color = ?"
	_, err := p.db.ExecContext(ctx, query, quantity, productID, size, color)
	return err
}
