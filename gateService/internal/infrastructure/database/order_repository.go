package database

import (
	"context"
	"database/sql"
	"gateService/internal/domain/entity"
)

type OrderRepositoryImpl struct {
	db *sql.DB
}

func NewOrderRepositoryImpl(db *sql.DB) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{db: db}
}

func (r *OrderRepositoryImpl) CreateOrder(ctx context.Context, order *entity.Order) error {
	query := `INSERT INTO orders (order_id, user_id, product_id, product_name, price, 
		selected_size, selected_color, user_name, phone, address, status, create_time) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, order.OrderID, order.UserID, order.ProductID, order.ProductName,
		order.Price, order.SelectedSize, order.SelectedColor, order.UserName,
		order.Phone, order.Address, order.Status, order.CreateTime)
	return err
}

func (r *OrderRepositoryImpl) UpdateOrder(ctx context.Context, order *entity.Order) error {
	query := `UPDATE orders SET product_id=?, product_name=?, price=?, description=?, 
		selected_size=?, selected_color=?, user_name=?, phone=?, address=?, 
		status=?, qr_code_url=? WHERE order_id=?`
	_, err := r.db.ExecContext(ctx, query, order.ProductID, order.ProductName, order.Price,
		order.Description, order.SelectedSize, order.SelectedColor, order.UserName,
		order.Phone, order.Address, order.Status, order.QRCodeURL, order.OrderID)
	return err
}

func (r *OrderRepositoryImpl) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	query := "UPDATE orders SET status = ? WHERE order_id = ? AND status = '待支付'"
	_, err := r.db.ExecContext(ctx, query, status, orderID)
	return err
}

func (r *OrderRepositoryImpl) DeleteOrder(ctx context.Context, orderID string) error {
	query := "DELETE FROM orders WHERE order_id = ?"
	_, err := r.db.ExecContext(ctx, query, orderID)
	return err
}

func (r *OrderRepositoryImpl) GetOrdersByUserID(ctx context.Context, userID, page, limit int) ([]*entity.Order, error) {
	query := `
		SELECT 
			o.order_id, o.create_time, o.status, 
			o.selected_size, o.selected_color, o.price,
			o.user_name, o.phone, o.address, o.product_name, p.image
		FROM orders o 
			JOIN products p ON o.product_id = p.product_id 
		WHERE 
			user_id = ?
		ORDER BY create_time DESC 
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		order := &entity.Order{}
		err := rows.Scan(&order.OrderID, &order.CreateTime, &order.Status, &order.SelectedSize,
			&order.SelectedColor, &order.Price, &order.UserName, &order.Phone, &order.Address, &order.ProductName, &order.Image)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
