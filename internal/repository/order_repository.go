package repository

import (
	"context"
	"database/sql"
	"fmt"
	"order-management-system/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) (int, error)
	GetOrderByID(ctx context.Context, id int) (model.Order, error)
	GetOrders(ctx context.Context, page, size int) ([]model.Order, int, error)
	UpdateOrderStatus(ctx context.Context, id int, status string) error
	CreateOrderItems(ctx context.Context, tx *sql.Tx, orderID int, items []model.OrderItem) error
}

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) CreateOrder(ctx context.Context, order model.Order) (int, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("can't start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var totalAmount float64
	for _, item := range order.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	/*
		start test rollback
	*/

	// var orderIDTest int
	// err = tx.QueryRowContext(
	// 	ctx,
	// 	"insert into orders (customer_name, total_amount, status, created_at, updated_at) values ($1, $2, $3, now(), now()) returning id",
	// 	"test", 100, "created",
	// ).Scan(&orderIDTest)

	// if err != nil {
	// 	return 0, fmt.Errorf("can't insert orders: %w", err)
	// }
	// log.Println("Create Order Test")

	// _, err = tx.ExecContext(
	// 	ctx,
	// 	"insert into order_items (order_id ,product_name, quantity, price) values ($1, $2, $3, $4)",
	// 	orderIDTest, "product_test", 3, 50,
	// )

	// if err != nil {
	// 	return 0, fmt.Errorf("can't insert order items: %w", err)
	// }
	// log.Println("Create Order Item Test")

	var orderID int
	err = tx.QueryRowContext(
		ctx,
		"insert into orders (customer_name, total_amount, status, created_at, updated_at) values ($1, $2, $3, now(), now()) returning id",
		order.CustomerName, totalAmount, "created",
	).Scan(&orderID)

	if err != nil {
		return 0, fmt.Errorf("can't insert orders: %w", err)
	}

	err = r.CreateOrderItems(ctx, tx, orderID, order.Items)
	if err != nil {
		return 0, err
	}

	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("can't commit transaction: %w", err)
	}

	return orderID, nil
}

func (r *PostgresOrderRepository) CreateOrderItems(ctx context.Context, tx *sql.Tx, orderID int, items []model.OrderItem) error {
	for _, item := range items {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		_, err := tx.ExecContext(
			ctx,
			"insert into order_items (order_id ,product_name, quantity, price) values ($1, $2, $3, $4)",
			orderID, item.ProductName, item.Quantity, item.Price,
		)

		if err != nil {
			return fmt.Errorf("can't order item: %w", err)
		}
	}
	return nil
}

func (r *PostgresOrderRepository) GetOrderByID(ctx context.Context, id int) (model.Order, error) {
	var order model.Order

	if ctx.Err() != nil {
		return order, ctx.Err()
	}

	err := r.db.QueryRowContext(
		ctx,
		"select id, customer_name, total_amount from orders where id = $1",
		id,
	).Scan(&order.ID, &order.CustomerName, &order.TotalAmount)

	if err != nil {
		if err == sql.ErrNoRows {
			return order, fmt.Errorf("not found order")
		}
		return order, fmt.Errorf("execption get order: %w", err)
	}

	rows, err := r.db.QueryContext(
		ctx,
		"select  product_name, quantity, price from order_items where order_id = $1",
		id,
	)

	if err != nil {
		return order, fmt.Errorf("execption get order items: %w", err)
	}
	defer rows.Close()

	order.Items = []model.OrderItem{}
	for rows.Next() {
		if ctx.Err() != nil {
			return order, ctx.Err()
		}
		var item model.OrderItem
		err := rows.Scan(&item.ProductName, &item.Quantity, &item.Price)
		if err != nil {
			return order, fmt.Errorf("execption read order items: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return order, fmt.Errorf("execption after read order items: %w", err)
	}

	return order, nil
}

func (r *PostgresOrderRepository) GetOrders(ctx context.Context, page, size int) ([]model.Order, int, error) {
	if ctx.Err() != nil {
		return nil, 0, ctx.Err()
	}

	offset := (page - 1) * size

	var total int
	err := r.db.QueryRowContext(ctx, "select count(1) from orders").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, fmt.Errorf("execption get count orders: %w", err)
	}

	rows, err := r.db.QueryContext(
		ctx,
		"select id, customer_name, total_amount, status, created_at, updated_at from orders order by id desc limit $1 offset $2",
		size, offset,
	)

	if err != nil {
		return nil, 0, fmt.Errorf("execption get orders: %w", err)
	}

	if ctx.Err() != nil {
		return nil, 0, ctx.Err()
	}

	defer rows.Close()

	orders := []model.Order{}
	for rows.Next() {

		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}

		var order model.Order
		err := rows.Scan(&order.ID, &order.CustomerName, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("execption read orders: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("execption after read orders: %w", err)
	}

	return orders, total, nil
}

func (r *PostgresOrderRepository) UpdateOrderStatus(ctx context.Context, id int, status string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("can't start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, err := tx.ExecContext(
		ctx,
		"update orders set status = $1, updated_at = now() where id = $2",
		status, id,
	)

	if err != nil {
		return fmt.Errorf("can't update order: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("execption update orders: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("not found order for update")
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}

	return nil
}
